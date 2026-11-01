package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

// App holds everything a handler needs. Methods on it are the handlers
// themselves, registered onto a mux in main.go.
type App struct {
	db        *sql.DB
	adminUser string
	adminHash string
}

func render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("render %s: %v", name, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

// ---- public site ----

type indexData struct {
	PageTitle string
	Posts     []Post
	Tags      []string
	IsAdmin   bool
}

// home shows published posts to the world, but drafts too — inline, with
// a "draft" badge — to a logged-in admin, the way an ordinary blog shows
// an author their own unpublished work without a separate screen.
func (a *App) home(w http.ResponseWriter, r *http.Request) {
	admin := isAdmin(a.db, r)
	posts, err := a.visiblePosts(admin)
	if err != nil {
		log.Printf("home: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	render(w, "index.html", indexData{Posts: posts, Tags: sortedTags(posts), IsAdmin: admin})
}

func (a *App) visiblePosts(admin bool) ([]Post, error) {
	if admin {
		return ListAll(a.db)
	}
	return ListPublished(a.db)
}

type postPageData struct {
	PageTitle string
	Post      Post
	Body      template.HTML
	Prev      *Post
	Next      *Post
	IsAdmin   bool
}

func (a *App) post(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	admin := isAdmin(a.db, r)

	post, err := GetPost(a.db, slug, !admin)
	if errors.Is(err, sql.ErrNoRows) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		log.Printf("post: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	body, err := post.Render()
	if err != nil {
		log.Printf("render markdown for %s: %v", slug, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	data := postPageData{PageTitle: post.Title, Post: post, Body: body, IsAdmin: admin}
	if all, err := a.visiblePosts(admin); err == nil {
		for i, p := range all {
			if p.Slug != slug {
				continue
			}
			if i > 0 {
				data.Next = &all[i-1]
			}
			if i < len(all)-1 {
				data.Prev = &all[i+1]
			}
			break
		}
	}
	render(w, "post.html", data)
}

type searchEntry struct {
	Slug  string   `json:"slug"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
	Text  string   `json:"text"`
}

func (a *App) searchIndex(w http.ResponseWriter, r *http.Request) {
	posts, err := ListPublished(a.db)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	entries := make([]searchEntry, 0, len(posts))
	for _, p := range posts {
		text := plainText(p.Body)
		const maxLen = 4000
		if len(text) > maxLen {
			text = text[:maxLen]
		}
		entries = append(entries, searchEntry{Slug: p.Slug, Title: p.Title, Tags: p.Tags, Text: text})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

func (a *App) feed(w http.ResponseWriter, r *http.Request) {
	posts, err := ListPublished(a.db)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	scheme := "http"
	if isTLS(r) {
		scheme = "https"
	}
	base := scheme + "://" + r.Host

	channel := rssChannel{
		Title:       "Workbench",
		Link:        base + "/",
		Description: "notes on an FPGA SmartNIC security/performance project.",
	}
	for _, p := range posts {
		link := base + "/posts/" + p.Slug
		channel.Items = append(channel.Items, rssItem{
			Title:       p.Title,
			Link:        link,
			GUID:        link,
			PubDate:     p.PostDate.Format(time.RFC1123Z),
			Description: p.Excerpt(),
		})
	}

	w.Header().Set("Content-Type", "application/rss+xml")
	w.Write([]byte(xml.Header))
	xml.NewEncoder(w).Encode(rssFeed{Version: "2.0", Channel: channel})
}

// ---- admin ----

type loginData struct {
	PageTitle string
	Error     string
	IsAdmin   bool
}

func (a *App) loginPage(w http.ResponseWriter, r *http.Request) {
	render(w, "admin_login.html", loginData{PageTitle: "Log in"})
}

func (a *App) loginSubmit(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username != a.adminUser || !checkPassword(a.adminHash, password) {
		render(w, "admin_login.html", loginData{PageTitle: "Log in", Error: "Incorrect username or password."})
		return
	}

	token, err := createSession(a.db)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	setSessionCookie(w, r, token, int(sessionTTL.Seconds()))
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (a *App) logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookie); err == nil {
		deleteSession(a.db, cookie.Value)
	}
	setSessionCookie(w, r, "", -1)
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

type dashboardData struct {
	PageTitle string
	Posts     []Post
	Message   string
	IsAdmin   bool
}

func (a *App) dashboard(w http.ResponseWriter, r *http.Request) {
	posts, err := ListAll(a.db)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	msg := ""
	switch {
	case r.URL.Query().Get("created") != "":
		msg = "Created “" + r.URL.Query().Get("created") + "”."
	case r.URL.Query().Get("updated") != "":
		msg = "Updated “" + r.URL.Query().Get("updated") + "”."
	}
	render(w, "admin_dashboard.html", dashboardData{PageTitle: "Admin", Posts: posts, Message: msg, IsAdmin: true})
}

type formData struct {
	PageTitle string
	Post      Post
	IsEdit    bool
	Error     string
	IsAdmin   bool
}

func (a *App) newPostPage(w http.ResponseWriter, r *http.Request) {
	render(w, "admin_form.html", formData{PageTitle: "New post", Post: Post{PostDate: time.Now()}, IsAdmin: true})
}

func (a *App) createPost(w http.ResponseWriter, r *http.Request) {
	form, ok := a.parsePostForm(w, r)
	if !ok {
		return
	}
	if form.Post.Title == "" {
		form.Error = "Title can't be empty."
		render(w, "admin_form.html", form)
		return
	}
	if slugify(form.Post.Title) == "" {
		form.Error = "Title needs at least one letter or number (it's used to build the URL)."
		render(w, "admin_form.html", form)
		return
	}
	slug, err := CreatePost(a.db, form.Post.Title, form.Post.Body, form.Post.TagsJoined(), form.Post.Published, form.Post.PostDate)
	if err != nil {
		log.Printf("create post: %v", err)
		form.Error = "Could not save the post."
		render(w, "admin_form.html", form)
		return
	}
	http.Redirect(w, r, "/admin?created="+slug, http.StatusSeeOther)
}

func (a *App) editPostPage(w http.ResponseWriter, r *http.Request) {
	post, err := GetPost(a.db, r.PathValue("slug"), false)
	if errors.Is(err, sql.ErrNoRows) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	render(w, "admin_form.html", formData{PageTitle: "Edit post", Post: post, IsEdit: true, IsAdmin: true})
}

func (a *App) updatePost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	form, ok := a.parsePostForm(w, r)
	if !ok {
		return
	}
	form.PageTitle = "Edit post"
	form.IsEdit = true
	form.Post.Slug = slug
	if form.Post.Title == "" {
		form.Error = "Title can't be empty."
		render(w, "admin_form.html", form)
		return
	}
	if err := UpdatePost(a.db, slug, form.Post.Title, form.Post.Body, form.Post.TagsJoined(), form.Post.Published, form.Post.PostDate); err != nil {
		log.Printf("update post %s: %v", slug, err)
		form.Error = "Could not save the post."
		render(w, "admin_form.html", form)
		return
	}
	http.Redirect(w, r, "/admin?updated="+slug, http.StatusSeeOther)
}

func (a *App) deletePost(w http.ResponseWriter, r *http.Request) {
	if err := DeletePost(a.db, r.PathValue("slug")); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// parsePostForm reads the shared new/edit form fields. The bool return is
// false if it already wrote an error response (bad date), so the caller
// should return immediately without rendering anything else.
func (a *App) parsePostForm(w http.ResponseWriter, r *http.Request) (formData, bool) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return formData{}, false
	}
	postDate, err := time.Parse("2006-01-02", r.FormValue("post_date"))
	if err != nil {
		postDate = time.Now()
	}
	post := Post{
		Title:     strings.TrimSpace(r.FormValue("title")),
		Body:      r.FormValue("body"),
		Tags:      parseTags(r.FormValue("tags")),
		Published: r.FormValue("published") == "on",
		PostDate:  postDate,
	}
	return formData{PageTitle: "New post", Post: post, IsAdmin: true}, true
}

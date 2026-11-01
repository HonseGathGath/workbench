package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

//go:embed templates/*.html
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

var tmpl = template.Must(template.ParseFS(templateFS, "templates/*.html"))

type indexData struct {
	PageTitle string
	Posts     []Post
}

type postData struct {
	PageTitle string
	Post      Post
	Prev      *Post
	Next      *Post
}

func main() {
	posts, err := LoadPosts("posts")
	if err != nil {
		log.Fatalf("loading posts: %v", err)
	}

	bySlug := make(map[string]int, len(posts))
	for i, p := range posts {
		bySlug[p.Slug] = i
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(http.FS(staticFS)))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		render(w, "index.html", indexData{Posts: posts})
	})

	mux.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.Path[len("/posts/"):]
		i, ok := bySlug[slug]
		if !ok {
			http.NotFound(w, r)
			return
		}
		data := postData{PageTitle: posts[i].Title, Post: posts[i]}
		// posts are sorted newest first, so "next" chronologically is the
		// previous slice entry and vice versa.
		if i > 0 {
			data.Next = &posts[i-1]
		}
		if i < len(posts)-1 {
			data.Prev = &posts[i+1]
		}
		render(w, "post.html", data)
	})

	addr := ":8080"
	log.Printf("Workbench running at http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("render %s: %v", name, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

const postsDir = "posts"

type writerListData struct {
	Posts []Post
}

type writerEditData struct {
	Post  Post
	IsNew bool
	Today string
}

// RunWriter starts the local-only writing UI. It binds to 127.0.0.1
// specifically (not ":<port>", which would listen on every interface),
// so nothing outside this machine can ever reach it — that's the whole
// reason it needs no login.
func RunWriter(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", writerList)
	mux.HandleFunc("GET /new", writerNew)
	mux.HandleFunc("GET /edit/{slug}", writerEdit)
	mux.HandleFunc("POST /api/save", writerSave)
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(writerStaticFS))))

	addr := "127.0.0.1:" + port
	url := "http://" + addr
	log.Printf("Workbench writer running at %s (only reachable from this machine)", url)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func renderWriter(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("render %s: %v", name, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func writerList(w http.ResponseWriter, r *http.Request) {
	posts, err := LoadPosts(postsDir)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	renderWriter(w, "writer_list.html", writerListData{Posts: posts})
}

func writerNew(w http.ResponseWriter, r *http.Request) {
	renderWriter(w, "writer_edit.html", writerEditData{IsNew: true, Today: time.Now().Format("2006-01-02")})
}

func writerEdit(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	posts, err := LoadPosts(postsDir)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	for _, p := range posts {
		if p.Slug == slug {
			renderWriter(w, "writer_edit.html", writerEditData{Post: p})
			return
		}
	}
	http.NotFound(w, r)
}

type saveRequest struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
	Tags  string `json:"tags"`
	Date  string `json:"date"`
	Body  string `json:"body"`
}

func writerSave(w http.ResponseWriter, r *http.Request) {
	// Requiring the exact application/json content type means a plain HTML
	// form or fetch() from another page can't hit this endpoint as a CORS
	// "simple request" (which skips preflight) — that would otherwise let
	// any page open in the same browser silently write posts while this
	// localhost-only, auth-free server is running.
	if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		http.Error(w, "expected application/json", http.StatusUnsupportedMediaType)
		return
	}

	var req saveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		http.Error(w, "title required", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		date = time.Now()
	}

	// slugify() runs on a client-supplied slug too, not just a generated
	// one — otherwise a crafted slug like "../../whatever" would pass
	// straight through to WritePostFile's filepath.Join and write outside
	// posts/, anywhere the OS user has access.
	slug := slugify(req.Slug)
	if slug == "" {
		slug, err = uniqueSlug(postsDir, req.Title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if err := WritePostFile(postsDir, slug, req.Title, date, parseTags(req.Tags), req.Body); err != nil {
		http.Error(w, "could not save: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"slug": slug})
}

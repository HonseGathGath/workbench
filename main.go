package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

//go:embed templates/*.html
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

var tmpl = template.Must(template.New("").Funcs(template.FuncMap{
	"join": strings.Join,
}).ParseFS(templateFS, "templates/*.html"))

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {
	case "serve":
		serve()
	case "hash-password":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, `usage: go run . hash-password "your password"`)
			os.Exit(1)
		}
		hash, err := hashPassword(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(hash)
	default:
		usage()
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage:
  go run . serve                     run the blog (reads env: PORT, DB_PATH, ADMIN_USERNAME, ADMIN_PASSWORD_HASH)
  go run . hash-password "password"  print a bcrypt hash to use as ADMIN_PASSWORD_HASH`)
	os.Exit(1)
}

// cacheControl adds a cache header for the embedded static assets: embed.FS
// entries have a zero ModTime, so http.FileServer can't emit Last-Modified
// or do conditional GETs on its own.
func cacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600")
		next.ServeHTTP(w, r)
	})
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func serve() {
	dbPath := env("DB_PATH", "workbench.db")
	db, err := openDB(dbPath)
	if err != nil {
		log.Fatalf("opening database %s: %v", dbPath, err)
	}
	defer db.Close()

	adminHash := os.Getenv("ADMIN_PASSWORD_HASH")
	if adminHash == "" {
		log.Fatal("ADMIN_PASSWORD_HASH is not set — generate one with: go run . hash-password \"your password\"")
	}

	app := &App{
		db:        db,
		adminUser: env("ADMIN_USERNAME", "admin"),
		adminHash: adminHash,
	}

	mux := http.NewServeMux()

	// public
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /posts/{slug}", app.post)
	mux.HandleFunc("GET /search-index.json", app.searchIndex)
	mux.HandleFunc("GET /feed.xml", app.feed)
	mux.Handle("GET /static/", cacheControl(http.FileServer(http.FS(staticFS))))

	// admin
	mux.HandleFunc("GET /admin/login", app.loginPage)
	mux.HandleFunc("POST /admin/login", app.loginSubmit)
	mux.HandleFunc("POST /admin/logout", requireAuth(db, app.logout))
	mux.HandleFunc("GET /admin", requireAuth(db, app.dashboard))
	mux.HandleFunc("GET /admin/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin", http.StatusMovedPermanently)
	})
	mux.HandleFunc("GET /admin/posts/new", requireAuth(db, app.newPostPage))
	mux.HandleFunc("POST /admin/posts/new", requireAuth(db, app.createPost))
	mux.HandleFunc("GET /admin/posts/{slug}/edit", requireAuth(db, app.editPostPage))
	mux.HandleFunc("POST /admin/posts/{slug}/edit", requireAuth(db, app.updatePost))
	mux.HandleFunc("POST /admin/posts/{slug}/delete", requireAuth(db, app.deletePost))

	addr := ":" + env("PORT", "8080")
	log.Printf("Workbench running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

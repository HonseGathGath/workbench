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

const outDir = "public"

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {
	case "new":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, `usage: go run . new "Post Title"`)
			os.Exit(1)
		}
		if err := NewPost(strings.Join(os.Args[2:], " ")); err != nil {
			log.Fatal(err)
		}
	case "build":
		if err := Build(outDir); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("built site into ./%s\n", outDir)
	case "serve":
		if err := Build(outDir); err != nil {
			log.Fatal(err)
		}
		serve(outDir)
	default:
		usage()
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage:
  go run . new "Post Title"   scaffold a new post in /posts
  go run . build              build the static site into ./public
  go run . serve               build, then serve ./public on :8080`)
	os.Exit(1)
}

func serve(dir string) {
	addr := ":8080"
	log.Printf("Workbench running at http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, http.FileServer(http.Dir(dir))))
}

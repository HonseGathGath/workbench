package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"strings"
)

func mustSub(fsys embed.FS, dir string) fs.FS {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		panic(err)
	}
	return sub
}

//go:embed templates/*.html
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

//go:embed writer_static/*
var writerStaticFSRaw embed.FS

// writerStaticFS strips the "writer_static/" prefix so it serves at
// /assets/ rather than /assets/writer_static/.
var writerStaticFS = mustSub(writerStaticFSRaw, "writer_static")

var tmpl = template.Must(template.New("").Funcs(template.FuncMap{
	"join": strings.Join,
}).ParseFS(templateFS, "templates/*.html"))

// chdirToRepo lets the installed "workbench" binary be run from anywhere:
// if WORKBENCH_REPO is set (one-time export after `go install .`), it
// switches into that directory first so posts/, site.url, and giscus.json
// are found there instead of relative to wherever the shell happens to be.
func chdirToRepo() {
	repo := os.Getenv("WORKBENCH_REPO")
	if repo == "" {
		return
	}
	if err := os.Chdir(repo); err != nil {
		log.Fatalf("WORKBENCH_REPO=%s: %v", repo, err)
	}
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	chdirToRepo()

	switch os.Args[1] {
	case "write":
		port := "8090"
		if len(os.Args) > 2 {
			port = os.Args[2]
		}
		RunWriter(port)
	case "build":
		if err := Build("public"); err != nil {
			log.Fatal(err)
		}
		fmt.Println("built site into ./public")
	default:
		usage()
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage:
  workbench write [port]   start the local writer (default port 8090, localhost only)
  workbench build          build the public static site into ./public

(also works as: go run . write / go run . build, from inside the repo)

Set WORKBENCH_REPO to the repo's absolute path to run either command
from any directory once "workbench" is installed on your PATH.`)
	os.Exit(1)
}

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

func main() {
	if len(os.Args) < 2 {
		usage()
	}

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
  go run . write [port]   start the local writer (default port 8090, localhost only)
  go run . build          build the public static site into ./public`)
	os.Exit(1)
}

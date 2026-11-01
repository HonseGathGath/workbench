package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var slugNonWord = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(title string) string {
	s := slugNonWord.ReplaceAllString(strings.ToLower(title), "-")
	return strings.Trim(s, "-")
}

// NewPost scaffolds a new Markdown file in /posts with today's date and
// the given title already filled into the front matter, so the only
// thing left to do is open it and write.
func NewPost(title string) error {
	slug := slugify(title)
	if slug == "" {
		return fmt.Errorf("title %q doesn't produce a usable slug", title)
	}
	if err := os.MkdirAll("posts", 0o755); err != nil {
		return err
	}
	path := filepath.Join("posts", slug+".md")
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("%s already exists", path)
	}

	content := fmt.Sprintf("---\ntitle: %s\ndate: %s\ntags: \n---\n\n",
		title, time.Now().Format("2006-01-02"))

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return err
	}
	fmt.Printf("wrote %s\n", path)
	return nil
}

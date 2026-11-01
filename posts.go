package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

type Post struct {
	Slug    string
	Title   string
	Date    time.Time
	Tags    []string
	Body    template.HTML
	Excerpt string
}

func (p Post) DateDisplay() string { return p.Date.Format("Jan 2, 2006") }

// LoadPosts reads every *.md file in dir, parses its front matter and
// Markdown body, and returns them sorted newest first.
func LoadPosts(dir string) ([]Post, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return nil, err
	}

	posts := make([]Post, 0, len(files))
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		post, err := parsePost(raw)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", f, err)
		}
		post.Slug = strings.TrimSuffix(filepath.Base(f), ".md")
		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool { return posts[i].Date.After(posts[j].Date) })
	return posts, nil
}

var frontMatterRe = regexp.MustCompile(`(?s)^---\n(.*?)\n---\n(.*)$`)

func parsePost(raw []byte) (Post, error) {
	m := frontMatterRe.FindSubmatch(raw)
	if m == nil {
		return Post{}, fmt.Errorf("missing --- front matter block")
	}
	var p Post
	for _, line := range strings.Split(string(m[1]), "\n") {
		key, val, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		val = strings.TrimSpace(val)
		switch strings.TrimSpace(key) {
		case "title":
			p.Title = strings.Trim(val, `"`)
		case "date":
			d, err := time.Parse("2006-01-02", val)
			if err != nil {
				return p, fmt.Errorf("bad date %q: %w", val, err)
			}
			p.Date = d
		case "tags":
			val = strings.Trim(val, "[]")
			for _, t := range strings.Split(val, ",") {
				t = strings.TrimSpace(strings.Trim(t, `"`))
				if t != "" {
					p.Tags = append(p.Tags, t)
				}
			}
		}
	}

	body := strings.TrimSpace(string(m[2]))
	p.Excerpt = excerpt(body)

	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(body), &buf); err != nil {
		return p, err
	}
	p.Body = template.HTML(buf.String())
	return p, nil
}

// excerpt grabs the first paragraph of raw Markdown and strips the most
// common inline syntax so it reads as plain text on the index page.
func excerpt(body string) string {
	first := strings.TrimSpace(strings.SplitN(body, "\n\n", 2)[0])
	first = regexp.MustCompile(`^#+\s*`).ReplaceAllString(first, "")
	first = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`).ReplaceAllString(first, "$1")
	first = regexp.MustCompile("[*_`]").ReplaceAllString(first, "")
	first = strings.Join(strings.Fields(first), " ")
	const max = 220
	if len(first) > max {
		first = strings.TrimSpace(first[:max]) + "…"
	}
	return first
}

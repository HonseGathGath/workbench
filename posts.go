package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Post struct {
	Slug    string
	Title   string
	Date    time.Time
	Updated time.Time // zero when the post has never been revised post-publish
	Tags    []string
	Body    string
	Seq     int // chronological position (oldest = 1), assigned by LoadPosts
}

func (p Post) DateDisplay() string    { return p.Date.Format("Jan 2, 2006") }
func (p Post) PostDateValue() string  { return p.Date.Format("2006-01-02") }
func (p Post) HasUpdate() bool        { return !p.Updated.IsZero() }
func (p Post) UpdatedDisplay() string { return p.Updated.Format("Jan 2, 2006") }
func (p Post) UpdatedDateValue() string {
	if p.Updated.IsZero() {
		return ""
	}
	return p.Updated.Format("2006-01-02")
}
func (p Post) SheetNo() string                { return fmt.Sprintf("%03d", p.Seq) }
func (p Post) TagsJoined() string             { return strings.Join(p.Tags, ", ") }
func (p Post) Excerpt() string                { return excerpt(p.Body) }
func (p Post) ReadingMinutes() int            { return readingMinutes(plainText(p.Body)) }
func (p Post) Render() (template.HTML, error) { return renderMarkdown(p.Body) }

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

	// Seq is each post's position in current chronological order, oldest
	// first — a sheet number, not a stable database ID. It's recomputed on
	// every load/build, so deleting an old post reflows the numbers after
	// it; nothing in the site (slugs, links, comment threads) depends on a
	// Seq value staying fixed.
	for i := range posts {
		posts[i].Seq = len(posts) - i
	}
	return posts, nil
}

var frontMatterRe = regexp.MustCompile(`(?s)^---\n(.*?)\n---\n(.*)$`)

func parsePost(raw []byte) (Post, error) {
	m := frontMatterRe.FindSubmatch(raw)
	if m == nil {
		return Post{}, fmt.Errorf("missing --- front matter block")
	}
	p := Post{Tags: []string{}}
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
		case "updated":
			if val != "" {
				d, err := time.Parse("2006-01-02", val)
				if err != nil {
					return p, fmt.Errorf("bad updated date %q: %w", val, err)
				}
				p.Updated = d
			}
		case "tags":
			p.Tags = parseTags(val)
		}
	}
	p.Body = strings.TrimSpace(string(m[2]))
	return p, nil
}

func parseTags(raw string) []string {
	tags := []string{}
	raw = strings.Trim(raw, "[]")
	for _, t := range strings.Split(raw, ",") {
		t = strings.TrimSpace(strings.Trim(t, `"`))
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}

var slugNonWord = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(title string) string {
	s := slugNonWord.ReplaceAllString(strings.ToLower(title), "-")
	return strings.Trim(s, "-")
}

// uniqueSlug slugifies the title and appends -2, -3, ... if a post with
// that slug already exists on disk.
func uniqueSlug(dir, title string) (string, error) {
	base := slugify(title)
	if base == "" {
		return "", fmt.Errorf("title %q doesn't produce a usable slug (needs at least one letter or number)", title)
	}
	slug := base
	for n := 2; ; n++ {
		if _, err := os.Stat(filepath.Join(dir, slug+".md")); os.IsNotExist(err) {
			return slug, nil
		}
		slug = fmt.Sprintf("%s-%d", base, n)
	}
}

// WritePostFile serializes title/date/updated/tags/body as a front-matter
// Markdown file at postsDir/slug.md, overwriting whatever was there. This is
// the only thing that ever creates or changes a post: the local writer calls
// it on every save, so a post file always round-trips through the same
// format whether it was first written by hand or through the UI. updated is
// the zero time.Time when the post has no revision marker.
func WritePostFile(postsDir, slug, title string, date, updated time.Time, tags []string, body string) error {
	if err := os.MkdirAll(postsDir, 0o755); err != nil {
		return err
	}
	var sb strings.Builder
	sb.WriteString("---\n")
	fmt.Fprintf(&sb, "title: %s\n", title)
	fmt.Fprintf(&sb, "date: %s\n", date.Format("2006-01-02"))
	if !updated.IsZero() {
		fmt.Fprintf(&sb, "updated: %s\n", updated.Format("2006-01-02"))
	}
	fmt.Fprintf(&sb, "tags: %s\n", strings.Join(tags, ", "))
	sb.WriteString("---\n\n")
	sb.WriteString(strings.TrimSpace(body))
	sb.WriteString("\n")
	return os.WriteFile(filepath.Join(postsDir, slug+".md"), []byte(sb.String()), 0o644)
}

// DeletePostFile removes postsDir/slug.md. slug must already be the
// sanitized output of slugify (callers are responsible for that — see
// writerDelete), so this never takes a raw path from a request.
func DeletePostFile(postsDir, slug string) error {
	path := filepath.Join(postsDir, slug+".md")
	if _, err := os.Stat(path); err != nil {
		return err
	}
	return os.Remove(path)
}

func sortedTags(posts []Post) []string {
	seen := map[string]bool{}
	var tags []string
	for _, p := range posts {
		for _, t := range p.Tags {
			if !seen[t] {
				seen[t] = true
				tags = append(tags, t)
			}
		}
	}
	sort.Strings(tags)
	return tags
}

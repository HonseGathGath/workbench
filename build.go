package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type indexData struct {
	Root       string
	PageTitle  string
	Posts      []Post
	Tags       []string
	TopbarMeta template.HTML // empty on the index: it has no single sheet/date to show
}

type postData struct {
	Root       string
	PageTitle  string
	Post       Post
	Body       template.HTML
	Prev       *Post
	Next       *Post
	Giscus     giscusConfig
	TopbarMeta template.HTML
}

// topbarMeta renders the title-block strip shown in the shared topbar on a
// post page: sheet number, publish date, and (only when set) the revision
// mark — the same three fields a drawing's title block carries in its
// sheet/rev/date corner. All inputs are Go-formatted numbers/dates, never
// user text, so building this as raw HTML here is safe.
func topbarMeta(p Post) template.HTML {
	s := fmt.Sprintf(`<span class="topbar-sheet">#%s</span> &middot; %s`, p.SheetNo(), p.DateDisplay())
	if p.HasUpdate() {
		s += fmt.Sprintf(` &middot; <span class="topbar-revised">revised %s</span>`, p.UpdatedDisplay())
	}
	return template.HTML(s)
}

// giscusConfig is read from ./giscus.json, which isn't committed until you
// run the one-time giscus.app setup. Until all four fields are filled in,
// Configured stays false and the template skips the comments embed
// entirely, rather than showing every visitor a broken "giscus is not
// installed on this repository" error box.
type giscusConfig struct {
	Repo       string `json:"repo"`
	RepoID     string `json:"repoId"`
	Category   string `json:"category"`
	CategoryID string `json:"categoryId"`
}

func (g giscusConfig) Configured() bool {
	return g.Repo != "" && g.RepoID != "" && g.CategoryID != ""
}

func readGiscusConfig() giscusConfig {
	var cfg giscusConfig
	b, err := os.ReadFile("giscus.json")
	if err != nil {
		return cfg
	}
	_ = json.Unmarshal(b, &cfg)
	if cfg.Category == "" {
		cfg.Category = "General"
	}
	return cfg
}

// Build renders the whole site (from /posts) into outDir: index.html, one
// posts/<slug>/index.html per post, the static assets, a search index,
// and an RSS feed. This is the only thing that ever runs against the
// public deploy path — it never touches the local writer's own templates
// or assets.
func Build(outDir string) error {
	posts, err := LoadPosts("posts")
	if err != nil {
		return fmt.Errorf("loading posts: %w", err)
	}

	if err := os.RemoveAll(outDir); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(outDir, "posts"), 0o755); err != nil {
		return err
	}

	if err := writeTemplate(filepath.Join(outDir, "index.html"), "index.html", indexData{
		Posts: posts,
		Tags:  sortedTags(posts),
	}); err != nil {
		return err
	}

	giscus := readGiscusConfig()
	for i, p := range posts {
		body, err := p.Render()
		if err != nil {
			return fmt.Errorf("rendering %s: %w", p.Slug, err)
		}
		data := postData{Root: "../../", PageTitle: p.Title, Post: p, Body: body, Giscus: giscus, TopbarMeta: topbarMeta(p)}
		if i > 0 {
			data.Next = &posts[i-1]
		}
		if i < len(posts)-1 {
			data.Prev = &posts[i+1]
		}
		dir := filepath.Join(outDir, "posts", p.Slug)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		if err := writeTemplate(filepath.Join(dir, "index.html"), "post.html", data); err != nil {
			return err
		}
	}

	if err := copyStaticAssets(outDir); err != nil {
		return err
	}
	if err := writeSearchIndex(outDir, posts); err != nil {
		return err
	}
	if err := writeFeed(outDir, posts); err != nil {
		return err
	}

	return nil
}

func writeTemplate(path, name string, data any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.ExecuteTemplate(f, name, data)
}

// copyStaticAssets copies only the public static/ embed — the writer's
// own assets live in a separate embedded FS (see write.go) that this
// never touches, so the deployed site can't end up shipping the writer.
func copyStaticAssets(outDir string) error {
	return fs.WalkDir(staticFS, "static", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := fs.ReadFile(staticFS, path)
		if err != nil {
			return err
		}
		rel := strings.TrimPrefix(path, "static/")
		return os.WriteFile(filepath.Join(outDir, rel), data, 0o644)
	})
}

type searchEntry struct {
	Slug  string   `json:"slug"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
	Text  string   `json:"text"`
}

func writeSearchIndex(outDir string, posts []Post) error {
	entries := make([]searchEntry, 0, len(posts))
	for _, p := range posts {
		text := plainText(p.Body)
		const maxLen = 4000
		if len(text) > maxLen {
			text = text[:maxLen]
		}
		entries = append(entries, searchEntry{Slug: p.Slug, Title: p.Title, Tags: p.Tags, Text: text})
	}
	data, err := json.Marshal(entries)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outDir, "search-index.json"), data, 0o644)
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

// writeFeed emits RSS 2.0. Item links are absolute if a one-line base URL
// is present in ./site.url (e.g. "https://user.github.io/repo"), otherwise
// they fall back to relative paths, which most feed readers still resolve
// fine against the feed's own URL.
func writeFeed(outDir string, posts []Post) error {
	base := strings.TrimSuffix(readBaseURL(), "/")

	link := func(slug string) string {
		if base == "" {
			return "posts/" + slug + "/"
		}
		return base + "/posts/" + slug + "/"
	}

	siteLink := "."
	if base != "" {
		siteLink = base + "/"
	}
	channel := rssChannel{
		Title:       "Workbench",
		Link:        siteLink,
		Description: "Breaking things until they make sense.",
	}
	for _, p := range posts {
		channel.Items = append(channel.Items, rssItem{
			Title:       p.Title,
			Link:        link(p.Slug),
			GUID:        link(p.Slug),
			PubDate:     p.Date.Format(time.RFC1123Z),
			Description: p.Excerpt(),
		})
	}

	out, err := xml.MarshalIndent(rssFeed{Version: "2.0", Channel: channel}, "", "  ")
	if err != nil {
		return err
	}
	out = append([]byte(xml.Header), out...)
	return os.WriteFile(filepath.Join(outDir, "feed.xml"), out, 0o644)
}

func readBaseURL() string {
	b, err := os.ReadFile("site.url")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

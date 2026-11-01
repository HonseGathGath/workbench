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
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type Post struct {
	Slug           string
	Title          string
	Date           time.Time
	Tags           []string
	Body           template.HTML
	Excerpt        string
	PlainText      string
	ReadingMinutes int
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
	p.PlainText = plainText(body)
	p.ReadingMinutes = readingMinutes(p.PlainText)

	html, err := renderMarkdown(body)
	if err != nil {
		return p, err
	}
	p.Body = html
	return p, nil
}

// renderMarkdown converts Markdown to HTML, rewriting footnote references
// (`[^1]`) into Tufte-style sidenotes: the note is inlined right where it's
// referenced, with CSS floating it into the margin on wide screens and a
// checkbox/label toggle revealing it inline on narrow ones.
func renderMarkdown(body string) (template.HTML, error) {
	md := goldmark.New(goldmark.WithExtensions(extension.Footnote))

	doc := md.Parser().Parse(text.NewReader([]byte(body)))

	// Render footnote fragments with a throwaway renderer: md.Renderer()
	// must not render anything before the sidenote renderer is registered
	// below, since AddOptions is a no-op after a renderer's first Render.
	fragmentRenderer := goldmark.New().Renderer()

	footnotes := map[int]string{}
	if last := doc.LastChild(); last != nil && last.Kind() == extast.KindFootnoteList {
		for fn := last.FirstChild(); fn != nil; fn = fn.NextSibling() {
			f, ok := fn.(*extast.Footnote)
			if !ok {
				continue
			}
			stripFootnoteBacklinks(f)
			var buf bytes.Buffer
			for c := f.FirstChild(); c != nil; c = c.NextSibling() {
				if err := fragmentRenderer.Render(&buf, []byte(body), c); err != nil {
					return "", err
				}
			}
			footnotes[f.Index] = buf.String()
		}
		doc.RemoveChild(doc, last)
	}

	md.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&sidenoteRenderer{footnotes: footnotes}, 0),
	))

	var buf bytes.Buffer
	if err := md.Renderer().Render(&buf, []byte(body), doc); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

// stripFootnoteBacklinks removes the "↩" return-link nodes goldmark appends
// to a footnote's content; sidenotes don't need them since the note is
// already sitting next to its reference.
func stripFootnoteBacklinks(f *extast.Footnote) {
	var container gast.Node = f
	if fc := container.LastChild(); fc != nil && gast.IsParagraph(fc) {
		container = fc
	}
	var next gast.Node
	for c := container.FirstChild(); c != nil; c = next {
		next = c.NextSibling()
		if c.Kind() == extast.KindFootnoteBacklink {
			container.RemoveChild(container, c)
		}
	}
}

type sidenoteRenderer struct {
	footnotes map[int]string
}

func (s *sidenoteRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(extast.KindFootnoteLink, s.renderFootnoteLink)
}

func (s *sidenoteRenderer) renderFootnoteLink(w util.BufWriter, source []byte, n gast.Node, entering bool) (gast.WalkStatus, error) {
	if !entering {
		return gast.WalkContinue, nil
	}
	link := n.(*extast.FootnoteLink)
	id := fmt.Sprintf("sn-%d-%d", link.Index, link.RefIndex)
	fmt.Fprintf(w, `<label for="%s" class="margin-toggle sidenote-number">%d</label>`, id, link.Index)
	fmt.Fprintf(w, `<input type="checkbox" id="%s" class="margin-toggle">`, id)
	fmt.Fprintf(w, `<span class="sidenote">%s</span>`, s.footnotes[link.Index])
	return gast.WalkContinue, nil
}

// excerpt grabs the first paragraph of raw Markdown and strips the most
// common inline syntax so it reads as plain text on the index page.
func excerpt(body string) string {
	first := strings.TrimSpace(strings.SplitN(body, "\n\n", 2)[0])
	first = stripMarkdown(first)
	const max = 220
	if len(first) > max {
		first = strings.TrimSpace(first[:max]) + "…"
	}
	return first
}

// plainText strips Markdown syntax from the whole post body, for reading
// time estimation and the search index.
func plainText(body string) string {
	return stripMarkdown(body)
}

var (
	headingRe  = regexp.MustCompile(`(?m)^#+\s*`)
	mdLinkRe   = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	footnoteRe = regexp.MustCompile(`\[\^[^\]]*\](?::[^\n]*)?`)
	emphasisRe = regexp.MustCompile("[*_`#>]")
)

func stripMarkdown(s string) string {
	s = headingRe.ReplaceAllString(s, "")
	s = mdLinkRe.ReplaceAllString(s, "$1")
	s = footnoteRe.ReplaceAllString(s, "")
	s = emphasisRe.ReplaceAllString(s, "")
	return strings.Join(strings.Fields(s), " ")
}

func readingMinutes(plain string) int {
	words := len(strings.Fields(plain))
	minutes := words / 200
	if minutes < 1 {
		minutes = 1
	}
	return minutes
}

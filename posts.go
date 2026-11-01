package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Post struct {
	ID        int64
	Slug      string
	Title     string
	Body      string
	Tags      []string
	Published bool
	PostDate  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p Post) DateDisplay() string            { return p.PostDate.Format("Jan 2, 2006") }
func (p Post) PostDateValue() string          { return p.PostDate.Format("2006-01-02") }
func (p Post) TagsJoined() string             { return strings.Join(p.Tags, ", ") }
func (p Post) Excerpt() string                { return excerpt(p.Body) }
func (p Post) ReadingMinutes() int            { return readingMinutes(plainText(p.Body)) }
func (p Post) Render() (template.HTML, error) { return renderMarkdown(p.Body) }

const timeLayout = "2006-01-02 15:04:05"

func scanPost(row scanner) (Post, error) {
	var p Post
	var tagsRaw, postDate, createdAt, updatedAt string
	var published int
	err := row.Scan(&p.ID, &p.Slug, &p.Title, &p.Body, &tagsRaw, &published, &postDate, &createdAt, &updatedAt)
	if err != nil {
		return Post{}, err
	}
	p.Tags = parseTags(tagsRaw)
	p.Published = published != 0
	p.PostDate, _ = time.Parse("2006-01-02", postDate)
	p.CreatedAt, _ = time.Parse(timeLayout, createdAt)
	p.UpdatedAt, _ = time.Parse(timeLayout, updatedAt)
	return p, nil
}

// scanner is satisfied by both *sql.Row and *sql.Rows.
type scanner interface {
	Scan(dest ...any) error
}

const postColumns = "id, slug, title, body, tags, published, post_date, created_at, updated_at"

func ListPublished(db *sql.DB) ([]Post, error) {
	return queryPosts(db, "SELECT "+postColumns+" FROM posts WHERE published = 1 ORDER BY post_date DESC, id DESC")
}

func ListAll(db *sql.DB) ([]Post, error) {
	return queryPosts(db, "SELECT "+postColumns+" FROM posts ORDER BY post_date DESC, id DESC")
}

func queryPosts(db *sql.DB, query string, args ...any) ([]Post, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		p, err := scanPost(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// GetPost returns sql.ErrNoRows if the slug doesn't exist. onlyPublished
// restricts it to published posts, so a public post URL can never surface
// a draft just because the slug was guessed.
func GetPost(db *sql.DB, slug string, onlyPublished bool) (Post, error) {
	query := "SELECT " + postColumns + " FROM posts WHERE slug = ?"
	if onlyPublished {
		query += " AND published = 1"
	}
	return scanPost(db.QueryRow(query, slug))
}

func CreatePost(db *sql.DB, title, body, tagsRaw string, published bool, postDate time.Time) (string, error) {
	slug, err := uniqueSlug(db, title)
	if err != nil {
		return "", err
	}
	now := time.Now().UTC().Format(timeLayout)
	_, err = db.Exec(
		`INSERT INTO posts (slug, title, body, tags, published, post_date, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		slug, title, body, normalizeTags(tagsRaw), boolToInt(published), postDate.Format("2006-01-02"), now, now,
	)
	return slug, err
}

// UpdatePost deliberately never changes the slug: it's the post's public
// URL, and giscus maps comment threads to it by pathname.
func UpdatePost(db *sql.DB, slug, title, body, tagsRaw string, published bool, postDate time.Time) error {
	now := time.Now().UTC().Format(timeLayout)
	_, err := db.Exec(
		`UPDATE posts SET title = ?, body = ?, tags = ?, published = ?, post_date = ?, updated_at = ?
		 WHERE slug = ?`,
		title, body, normalizeTags(tagsRaw), boolToInt(published), postDate.Format("2006-01-02"), now, slug,
	)
	return err
}

func DeletePost(db *sql.DB, slug string) error {
	_, err := db.Exec("DELETE FROM posts WHERE slug = ?", slug)
	return err
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func parseTags(raw string) []string {
	var tags []string
	for _, t := range strings.Split(raw, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}

func normalizeTags(raw string) string {
	return strings.Join(parseTags(raw), ", ")
}

var slugNonWord = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(title string) string {
	s := slugNonWord.ReplaceAllString(strings.ToLower(title), "-")
	return strings.Trim(s, "-")
}

// uniqueSlug slugifies the title and appends -2, -3, ... on collision.
func uniqueSlug(db *sql.DB, title string) (string, error) {
	base := slugify(title)
	if base == "" {
		return "", fmt.Errorf("title %q doesn't produce a usable slug", title)
	}
	slug := base
	for n := 2; ; n++ {
		var exists int
		if err := db.QueryRow("SELECT COUNT(*) FROM posts WHERE slug = ?", slug).Scan(&exists); err != nil {
			return "", err
		}
		if exists == 0 {
			return slug, nil
		}
		slug = fmt.Sprintf("%s-%d", base, n)
	}
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

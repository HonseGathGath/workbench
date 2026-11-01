# Workbench

A minimal personal engineering blog. Posts are Markdown files with front
matter in `/posts` — no database, no CMS, no auth.

Run it:

    make dev

(or just `go run .`) then open http://localhost:8080.

To add a post, drop a new `.md` file in `/posts` with a front-matter block
like the existing ones (`title`, `date: YYYY-MM-DD`, `tags`), then restart.

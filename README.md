# Workbench

A minimal personal engineering blog. Posts are Markdown files with front
matter in `/posts` — no database, no CMS, no auth. `go run .` is a small
CLI with three subcommands.

## Write

    go run . new "Post Title"

Scaffolds `/posts/post-title.md` with today's date already filled in.
Footnotes (`[^1]` ... `[^1]: text`) render as Tufte-style sidenotes.

## Preview locally

    make dev

Builds the site into `./public` and serves it on http://localhost:8080.
Rerun after editing a post.

## Publish

    git push

Once GitHub Pages is set up (one-time): in the repo's Settings → Pages,
set Source to "GitHub Actions". `.github/workflows/deploy.yml` builds and
deploys on every push to `main`. For a fully correct RSS feed, put your
Pages URL (e.g. `https://you.github.io/repo`) as a single line in
`site.url` at the repo root.

## Comments (giscus)

One-time setup on your GitHub repo:

1. Enable Discussions (Settings → General → Features).
2. Install the giscus app: https://github.com/apps/giscus
3. Go to https://giscus.app, fill in your repo, and copy the
   `data-repo`, `data-repo-id`, `data-category`, `data-category-id`
   values it generates into the `<script>` block in `templates/post.html`
   (currently placeholder `REPLACE_ME` values).

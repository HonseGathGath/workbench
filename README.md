# Workbench

A minimal personal blog with a real admin: public pages are served from
Markdown-ish posts stored in SQLite, and everything (write/edit/delete,
draft vs. published) happens through a logged-in `/admin` UI in the
browser — no editor or terminal needed once it's running.

## Run it locally

    go run . hash-password "pick a password"

Copy the printed hash, then:

    ADMIN_PASSWORD_HASH='<paste the hash>' make dev

Visit http://localhost:8080 for the public site, http://localhost:8080/admin
to write. Username defaults to `admin` (override with `ADMIN_USERNAME`).
Posts live in `workbench.db` (SQLite, gitignored) next to the binary.

Env vars `serve` reads: `PORT` (default 8080), `DB_PATH` (default
`workbench.db`), `ADMIN_USERNAME` (default `admin`), `ADMIN_PASSWORD_HASH`
(required).

## Deploy (Railway)

A plain Go binary with a SQLite file needs a host that runs your code
continuously and keeps a writable disk around — Railway does both without
asking for a Dockerfile.

1. Push this repo to GitHub.
2. On railway.app: New Project → Deploy from GitHub repo → pick it.
   Railway detects `go.mod` and builds/runs the binary automatically.
3. In the service's Variables tab, set `ADMIN_PASSWORD_HASH` (from
   `go run . hash-password "..."`) and optionally `ADMIN_USERNAME`.
4. In the service's Settings → Volumes, attach a volume mounted at
   `/data`, and set the `DB_PATH` variable to `/data/workbench.db`.
   Without this, the database resets on every redeploy.
5. Settings → Networking → Generate Domain. That URL is your live blog.

## Comments (giscus)

One-time setup on your GitHub repo:

1. Enable Discussions (Settings → General → Features).
2. Install the giscus app: https://github.com/apps/giscus
3. Go to https://giscus.app, fill in your repo, and copy the
   `data-repo`, `data-repo-id`, `data-category`, `data-category-id`
   values it generates into the `<script>` block in `templates/post.html`
   (currently placeholder `REPLACE_ME` values).

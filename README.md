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

## How to use it

**Logging in.** Go to `/admin/login` (there's a quiet "Admin" link in the
footer of every page) and sign in with your username/password.

**Writing a post.** Once logged in, every page — including the home
page — gets a small nav next to the site title: *Dashboard*, *+ New
post*, *Log out*. Click **+ New post** from wherever you are. Fill in:

- **Title** — also becomes the post's URL (`/posts/your-title`). The
  URL is fixed at creation and won't change even if you edit the title
  later (so old links and comment threads keep working).
- **Tags** — comma-separated, shown on the post and usable as filters
  on the home page.
- **Date** — defaults to today; backdate it if you're writing about
  something that happened earlier.
- **Body** — Markdown. Tables, `~~strikethrough~~`, and bare URLs all
  work. A footnote (`[^1]` in the text, `[^1]: the note` anywhere else
  in the body) renders as a margin sidenote instead of a bottom-of-page
  footnote.
- **Published** — leave it unchecked to save as a draft. Drafts never
  appear on the public site, in search, or in the RSS feed — only you
  (logged in) can see them, both in the Dashboard and inline on the home
  page with a "draft" badge.

**Editing or deleting.** While logged in, every post — published or
draft — has an **Edit** link right on its home-page card and on its own
post page, plus a full list with Edit/Delete at `/admin`. Deleting asks
for confirmation first and can't be undone.

**Publishing a draft.** Edit it, check **Published**, save. It shows up
everywhere immediately.

**Logging out.** Click **Log out** in the header nav from any page.

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

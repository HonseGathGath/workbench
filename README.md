# Workbench

A minimal personal blog. The public site is a plain static build with no
backend, login, or database. Writing happens only through a local tool
that binds to `127.0.0.1` and is never reachable from outside your
machine — see `Workbench-Setup-Guide.pdf` for the full walkthrough
(first-time deploy, day-to-day writing, publishing, moving to a new
computer). Quick reference:

    go run . write     # local writer at http://localhost:8090
    go run . build     # builds the public site into ./public

Or install it once as a global command (`go install .`, then set
`WORKBENCH_REPO` to this directory's absolute path — see the PDF) and
run `workbench write` / `workbench build` from anywhere.

Posts are Markdown files with front matter in `/posts`. Footnotes
(`[^1]` ... `[^1]: text`) become margin sidenotes. Deploy is GitHub
Pages via `.github/workflows/deploy.yml`, which runs on every push to
`main`. Comments (giscus) and the RSS feed's absolute URLs activate once
you create `giscus.json` / `site.url` respectively — both documented in
the PDF.

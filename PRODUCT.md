# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Stack

Static site generator written in Go (no framework), `html/template` + goldmark for Markdown, no client-side build step. Existing codebase; not a greenfield choice.

## Users

Two roles, same person on one side:

- **Author (Ghaith Lassoued):** writes and publishes posts through a local-only writer tool that runs on `127.0.0.1`; never exposed publicly. This is a solo blog — no multi-author or collaboration workflow.
- **Readers:** anyone — no fixed niche. The author is self-taught, writing to teach themselves as they go, so posts can span hardware/security engineering (FPGA, RTL, systems internals) and general software engineering (project logs, debugging, tooling) depending on what they're currently working on.

## Product Purpose

A personal technical blog/notebook. The author writes posts as part of learning in public — working through a problem and writing it up, not writing polished tutorials after the fact. Success is a reading experience serious enough to match that intent: a place worth linking to, not a throwaway notes dump.

## Positioning

What distinguishes this blog from a generic engineering blog: **process over results.** Posts are expected to show the actual path to a working solution — debugging, dead ends, wrong turns — not just a clean retrospective writeup of the final answer. ("Breaking things until they make sense.") This is a content/voice commitment for future posts, not a layout requirement (no special "process" UI chrome implied).

## Operating Context

- Posts are Markdown files with YAML-ish front matter (title/date/tags) in `/posts`, edited through a local writer at `localhost:8090` (title, tags, date, Markdown body, autosave, live preview).
- Publishing = committing the post file and pushing to `main`; GitHub Actions rebuilds and deploys to GitHub Pages on every push. No separate "publish" step or draft/published toggle.
- Footnotes (`[^1]`) render as Tufte-style margin sidenotes on the built post page.
- Optional integrations activate only once configured: RSS absolute links (`site.url`), comments via giscus/GitHub Discussions (`giscus.json`). Both degrade silently (relative links / no comments section) when unconfigured.
- Currently zero real posts exist — the blog has not been written in yet. No existing post content constrains the redesign.

## Capabilities and Constraints

Confirmed existing functionality the redesign must preserve: home page (post list, search, tag filtering), individual post pages, RSS feed, light/dark theme toggle, giscus comments (optional), sidenote footnotes, reading-time estimate, syntax highlighting (highlight.js). The local writer (create/edit posts) is a separate non-public surface with its own minimal UI, not the public reading surface, but is expected to visually belong to the same product.

No CMS, no database, no server-rendered dynamic routes on the public site — it is a plain static build. Any new feature (e.g. deleting a post) must fit this model: a local-only, file-based operation through the same writer tool, not a new public/server surface.

## Brand Commitments

- Name: **Workbench**. Established, not open for the redesign.
- Current tagline: "Breaking things until they make sense." Set deliberately in the prior session; reflects the positioning above.
- Existing favicon mark: a letterform "W" on a solid accent-green square (`#3f6a5c`), Georgia/serif. Minor/replaceable if the new palette changes the accent color, but the "W" monogram concept can carry over.
- No existing logo lockup, imagery, or illustration system — typography and layout currently carry the entire identity.

## Evidence on Hand

No real post content exists yet (0 posts in `/posts`). The redesign has no existing articles' content/length/formatting patterns to anchor to — base typographic decisions (line length, scale) on general technical-writing conventions (code blocks, footnotes, headings, links, occasional images) rather than a specific existing post.

## Product Principles

1. Solo, static, file-based — no feature should require a database, server-side auth, or public write surface.
2. Reading experience is the product; the writer tool is a means to an end and should stay minimal, though visually consistent with the public site.
3. Voice over polish — the site should support process-y, in-progress-feeling writing without looking unfinished itself.
4. No draft/publish ceremony — a post existing in `/posts` and pushed to `main` is the entire publishing model; don't design around a state that doesn't exist.
5. Everything optional (comments, absolute RSS links) must degrade invisibly, never as a broken/placeholder UI state.

## Accessibility & Inclusion

No standard explicitly required by the user; default to solid baseline practice (visible focus states, sufficient contrast in both themes, responsive down to mobile, reduced-motion respected) since no stricter requirement was raised.

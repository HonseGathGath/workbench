---
name: Workbench
description: A solo technical notebook that shows its editing history in the margin, like a numbered engineering-drawing sheet.
colors:
  cool-paper: "#eef0ea"
  raised-paper: "#e3e6dc"
  graphite-ink: "#1c1e1a"
  muted-ink: "#6b7066"
  hairline-rule: "#cfd3c6"
  redline: "#b23a2c"
  redline-ink: "#8a2c21"
typography:
  display:
    fontFamily: "Zilla Slab, Georgia, serif"
    fontSize: "clamp(2rem, 4.5vw, 3rem)"
    fontWeight: 700
    lineHeight: 1.1
    letterSpacing: "-0.01em"
  title:
    fontFamily: "Zilla Slab, Georgia, serif"
    fontSize: "1.3rem"
    fontWeight: 600
    lineHeight: 1.2
  body:
    fontFamily: "Spectral, Georgia, serif"
    fontSize: "19px"
    fontWeight: 400
    lineHeight: 1.65
  dek:
    fontFamily: "Spectral, Georgia, serif"
    fontSize: "1.1rem"
    fontWeight: 400
    lineHeight: 1.5
  label:
    fontFamily: "JetBrains Mono, SF Mono, monospace"
    fontSize: "0.78rem"
    fontWeight: 500
    letterSpacing: "0.01em"
rounded:
  none: "0px"
  circle: "50%"
spacing:
  xs: "0.5rem"
  sm: "1rem"
  md: "1.6rem"
  lg: "2rem"
  xl: "3rem"
components:
  entry-row:
    textColor: "{colors.graphite-ink}"
    typography: "{typography.title}"
    padding: "1.6rem 0"
  entry-row-sheet:
    textColor: "{colors.muted-ink}"
    typography: "{typography.label}"
  topbar:
    backgroundColor: "{colors.cool-paper}"
    textColor: "{colors.graphite-ink}"
    padding: "1.7rem 1.5rem 1.2rem"
  tag-chip:
    textColor: "{colors.muted-ink}"
    typography: "{typography.label}"
  tag-chip-active:
    textColor: "{colors.redline-ink}"
    typography: "{typography.label}"
  w-delete:
    backgroundColor: "transparent"
    textColor: "{colors.redline-ink}"
    padding: "0.3rem 0.7rem"
    typography: "{typography.label}"
  w-delete-hover:
    backgroundColor: "{colors.redline-ink}"
    textColor: "{colors.cool-paper}"
---

# Design System: Workbench

## Overview

**Creative North Star: "Rev Mark"**

Workbench reads as a single numbered sheet in an ongoing set of engineering drawings, not a magazine feature. Every post carries a sheet number assigned by chronological position (oldest = 1), and the shared masthead is built as a title block — brand on the left, sheet/date/revision stamp on the right, divided from the body by a heavy rule — rather than a conventional nav bar. On wide viewports a vertical margin rail restates the sheet number and, only when the post actually has an `updated` front-matter date, a second accent-colored revision tick beneath it. This is the "process over results" positioning made structural rather than asserted in copy: the reader sees, without being told, that this piece was corrected after first publishing.

The palette was deliberately revised mid-build away from a warm-cream-and-rust first pass — too close to the single most common AI-generated-design cliché — to a cooler, grayer paper tone paired with a genuine "redline red" accent. Red was chosen because it is the real-world convention for marking revisions and corrections (red-pen, redline drawings), not as arbitrary decoration; that causal link between subject matter and color choice is a load-bearing decision, not a style preference. Fonts follow the same logic: Zilla Slab, Spectral, and JetBrains Mono were chosen specifically to avoid the Fraunces/Newsreader/Space Grotesk/IBM Plex Mono combination flagged as an overused AI-default pairing.

The system is flat and content-first: no card grids, no shadows, no rounded corners on content surfaces. Depth and hierarchy come from rules (hairline dividers, heavy 2–3px title-block rules) and from the accent color's restraint, not from elevation.

**Key Characteristics:**
- Title-block masthead (brand + sheet/date stamp, heavy rule beneath) in place of a conventional nav bar.
- A flat "sheet list" post index with a leader-dot connector between title and date — never a card grid.
- A margin rail (≥1320px only) restating the sheet number and, at most, one revision tick.
- Redline-red accent reserved for revision marks and active/interactive state — never static decoration.
- One shared token system across the public site and the local writer tool.

## Colors

A cool, graphite-and-paper neutral field with a single warm "redline" accent held in reserve for revisions and interaction.

### Primary
- **Redline** (`#b23a2c` light / `#e07a68` dark): the revision-mark accent — sidenote markers, blockquote rule, tag-chip active/hover underline, focus outline, selection highlight.
- **Redline Ink** (`#8a2c21` light / `#ef9686` dark): the same accent darkened/lightened for text-weight use — hover states on links and the brand, `revised {date}` labels, the writer tool's Delete button.

### Neutral
- **Cool Paper** (`#eef0ea` light / `#13150f` dark): page background.
- **Raised Paper** (`#e3e6dc` light / `#1b1e17` dark): slightly-lifted panels — the dark-mode-aware inline sidenote fallback box, writer-tool input and editor panes.
- **Graphite Ink** (`#1c1e1a` light / `#e4e6dd` dark): primary text color, heavy title-block rules.
- **Muted Ink** (`#6b7066` light / `#a3a89b` dark): secondary text — dates, metadata, deks, tag chips at rest.
- **Hairline Rule** (`#cfd3c6` light / `#333730` dark): 1px dividers between entries, table rows, sidenote borders, the leader-dot connector.

### Named Rules
**The Redline Rule.** `redline`/`redline-ink` are reserved for actual revision marks and active/interactive state (sidenote markers, `.tag-chip.active`, hover, `:focus-visible`, the writer's Delete button) — never for decorating a plain sheet number or other static label. This was enforced as a corrective fix during finish review: an earlier draft colored every sheet number with the accent, diluting it into decoration, and it was corrected to neutral ink specifically because of this rule.

## Typography

**Display Font:** Zilla Slab (with Georgia, serif fallback)
**Body Font:** Spectral (with Georgia, serif fallback) — italic for deks and excerpts
**Label/Mono Font:** JetBrains Mono (with SF Mono, monospace fallback)

**Character:** A sturdy, engineering-adjacent slab serif for headlines paired with a calm literary serif for reading text, with every piece of metadata — dates, sheet numbers, tags, byline, nav labels, table cells — set in monospace so it reads as stamped data, not prose.

### Hierarchy
- **Display** (700, `clamp(2rem, 4.5vw, 3rem)`, 1.1): post/index headline (`.headline`).
- **Title** (600, 1.3rem): sheet-list entry titles (`.entry-title`), writer table row titles.
- **Dek** (400 italic, 1.1rem): the post/index subhead directly under the headline (`.dek`), max-width 34em.
- **Body** (400, 19px / line-height 1.65, 660px measure): post reading copy.
- **Label** (500–600, 0.75–0.78rem, mono, tabular-nums on dates/sheet numbers): every piece of metadata — byline, entry date/sheet/meta, topbar meta, tag chips, table headers and cells.

### Named Rules
**The Mono-Metadata Rule.** Anything that is data about a post rather than the post's own prose — dates, sheet numbers, tags, byline, reading time, nav labels — is set in JetBrains Mono. Zilla Slab is reserved for headlines and titles; Spectral is reserved for reading prose and deks.

## Layout

Single reading column, `max-width: 660px`, centered (`.topbar`, `.masthead`, `.entry-list`, `.post`, `.post-nav`, `.comments`, `.controls` all share this measure). The post index is a flat "sheet list" — each entry is a row with a sheet number, title, a dotted leader connecting title to date, then a mono metadata line and dek beneath — explicitly not a card grid; the brief's content-first requirement ruled cards out.

Entries are divided by a 1px hairline rule (`border-top: 1px solid var(--rule)`) at `1.6rem` vertical rhythm; the first entry drops its top rule. The masthead and post masthead are divided from the body by a heavy 2–3px rule (`.topbar` border-bottom, `.title-rule`) — the title-block convention carried through the whole system.

At `≥1320px`, a vertical margin rail appears to the left of the reading column (`position: absolute`, `writing-mode: vertical-rl`, offset `-100% - 2.5rem` from the column edge) showing the sheet number and, only if the post has an `updated` date, one revision tick beneath it. Tufte-style sidenotes float in the right-hand margin at the same `≥1320px` breakpoint (`width: 230px`, `margin-right: -300px`); below that width they collapse to an inline checkbox-toggle box with a raised-paper background. Below `480px`, the leader-dot connector is dropped and the entry row wraps.

## Elevation & Depth

Flat by design — no `box-shadow` anywhere in the system. Depth and hierarchy come entirely from rule weight (1px hairline dividers vs. 2–3px heavy title-block rules) and from background contrast between `--bg` and `--bg-raised`, not from elevation.

### Named Rules
**The Flat Ledger Rule.** Surfaces never lift off the page. A raised panel (sidenote fallback box, writer input/editor panes) is distinguished only by `--bg-raised` and a hairline border — never a shadow.

## Shapes

No border-radius on any content surface — entry rows, inputs, buttons, panels, and the post-body code block are all square-cornered, consistent with a flat, hairline-rule, drafting-adjacent aesthetic. The one exception is the circular `#theme-toggle` icon button (`border-radius: 50%`), which is a deliberate icon-button convention, not a departure from the no-rounding rule for content surfaces.

## Components

### Buttons
- **Shape:** square (0px radius), except the circular theme toggle noted above.
- **Primary (writer "+ New post"):** `.w-button` — `redline-ink` background, `cool-paper` text, mono label, `0.4rem 0.9rem` padding.
- **Destructive (writer Delete):** `.w-delete` — border + text in `redline-ink` *at rest*, not just on hover, so a destructive action reads as such from a static glance; on hover it inverts to a filled `redline-ink` background with paper text.
- **Icon (theme toggle):** transparent, 1px `--rule` border, circular, `currentColor` icon; hover shifts border/icon to the accent.

### Chips
- **Style (tag chips):** no background or border by default — mono label text in `--ink-muted`, bottom border transparent.
- **State:** hover and `.active` both switch text to `redline-ink` and reveal a solid `redline` bottom border; this is the only place a tag reads as "selected."

### Inputs / Fields
- **Style:** no box; a baseline-only underline (`border-bottom: 1px solid var(--rule)`), italic placeholder text, mono/body font depending on context. Writer tool fields use a full 1px `--rule` border on a `--bg-raised` panel instead, since they're a denser authoring UI rather than inline search.
- **Focus:** underline/border shifts to `--accent`; global `:focus-visible` uses a 2px `--accent` outline with 2px offset everywhere.

### Navigation
- **Topbar:** a title block, not a nav bar — brand (Zilla Slab, bold) on the left, optional mono meta string on the right, divided from the page body by a heavy 2px `--ink` rule. Brand hovers to `redline-ink`.
- **Post nav (prev/next):** mono labels in `--ink-muted` beneath a hairline top rule; hover shifts to `redline-ink`.

### Sheet Rail (signature component)
A vertical, mono-set margin column (`≥1320px` only) running beside the reading column: the sheet number in muted ink, and — only when the post has a real `updated` front-matter date — one accent-colored "revised {date}" tick beneath it. The rail deliberately holds at most two ticks; the data model stores only one optional `updated` date per post with no edit-by-edit history, so this is not a truncated history log but the full intended shape of the device. The same sheet number + revision pairing is echoed in the entry-list row and the post's title-block meta string, so all three stay in sync from the same build-time computation.

## Do's and Don'ts

### Do:
- **Do** set every piece of post metadata (dates, sheet numbers, tags, byline, reading time, nav labels, table cells) in JetBrains Mono.
- **Do** reserve `redline`/`redline-ink` for revision marks and active/interactive state (The Redline Rule).
- **Do** keep the post index as a flat sheet-list row layout with a dotted leader connector — never a card grid.
- **Do** style a destructive action (writer Delete) with its warning color at rest, not only on hover.
- **Do** convey depth with rule weight and background contrast only; no `box-shadow`.

### Don't:
- **Don't** color a static label (a plain sheet number, a date with no revision) with the accent — it's reserved for actual revision marks and interaction, per the Redline Rule.
- **Don't** add rounded corners to content surfaces (buttons, inputs, panels, code blocks); the circular theme toggle is the one established exception, not a precedent for others.
- **Don't** render an icon as a Unicode/emoji glyph. The theme toggle is an authored inline SVG (`{{define "theme-icon"}}` in `templates/partials.html`, reused identically on the public site and the writer tool) specifically because an earlier draft's raw "◐" character was flagged as a craft-floor violation and replaced.
- **Don't** expand the margin rail into a multi-revision history log — the data model stores only one optional `updated` date per post, so the rail holds at most a sheet number and one revision tick, by design.

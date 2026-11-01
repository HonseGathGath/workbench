---
title: Mapping the field — who's actually doing FPGA research in Europe right now
date: 2026-09-22
tags: research, landscape, europe
---

I wanted a map of who's working on FPGA security and reconfigurable-computing research in Europe, partly out of curiosity and partly because this project is supposed to feed into research applications later. My first instinct was to start from university rankings and work down to labs. That turned out to be almost useless, and it took me longer than I'd like to admit to notice why.

Rankings measure institutions, not groups, and research happens in groups. A well-ranked university can have zero people currently working on anything adjacent to FPGA security, while a less famous one can have a three-person lab publishing steadily at the venues that actually matter for this area. Going in by reputation, I kept landing on departments with a single FPGA-adjacent course and no active research output. Not useless information, but not what I needed.

So I flipped the approach: start from the papers, not the institutions. Concretely, that meant pulling the last two or three years of proceedings from the venues where this kind of work actually gets published — FPL, FCCM, DATE, a few sessions at CHES — and just reading author affiliations off the papers themselves. That surfaces something a ranking never will: who is *currently* active, publishing with students (meaning there's an open lab, not a solo retired-professor situation), and working on something close enough to my threat model to be worth contacting.

A few things fell out of doing it this way:

- **Activity clusters, it doesn't scatter evenly.** A handful of groups across Germany, the Netherlands, Belgium, and the UK account for a disproportionate share of recent FPGA-security output. That's useful to know on its own — it tells me where the community's center of gravity actually is, versus where I assumed it was.
- **Co-authorship networks reveal the actual social graph.** Once I had a list of active authors, cross-referencing who publishes together repeatedly mapped out collaborations that aren't visible from any single department's website. Some "separate" labs are really one extended network.
- **Recency matters more than volume.** A group with fewer total papers but consistent output in the last 18 months is a better target than a group with a big back-catalog and nothing recent — people move, labs wind down, funding runs out.

What I ended up with is a working list of roughly a dozen groups, each tagged with what specifically they're working on and how it overlaps with my three threats from the last post. I'm treating this as a living document rather than a one-time survey — I'll revisit it before I get to the point of actually reaching out to anyone, since "active six months ago" and "active now" aren't the same claim.

Next practical step: read deeply into two or three of the closest-match groups' recent papers, not just the abstracts, so I actually understand their methodology before I start comparing it to mine.

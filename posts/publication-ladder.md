---
title: Publication ladder — picking targets before I have results
date: 2026-10-18
tags: research, process, publishing
---

This felt backwards when I first considered doing it: pick publication venues before I have anything resembling a result. The usual story is you do the work, get something worth telling people about, and then go find somewhere to tell them. I'm doing it the other way round, on purpose, and I think it's the right call for a project of this shape.

The reasoning is simple once stated. A capstone project has a fixed amount of time and a strong temptation to spend all of it building, with "evaluation" as a rushed afterthought in the last two weeks. If I don't know in advance what a credible result looks like to an outside audience, I will not build the measurement infrastructure to produce one — I'll just build the thing and hope the numbers look good whenever I finally check. Picking a venue first forces me to ask, now, what a reviewer at that venue would actually want to see, which tells me what to instrument from day one instead of retrofitting benchmarks after the fact.

So I built a small ladder rather than picking one target, because I don't yet know how far this will get:

**Rung 1 — workshop or poster track.** Realistic floor outcome. Workshops attached to the bigger FPGA venues (FPL, FCCM) take shorter, more exploratory submissions. Even a partial result — "here's the attack, here's a mitigation, here's the overhead" for just one of my three threats — is workshop-shaped. This is the outcome I'm planning the project's timeline around, not the stretch goal.

**Rung 2 — a full workshop or short paper.** If the bitstream-tampering thread (the one I'm starting with) produces a clean attack/mitigation/overhead triplet with real numbers, that's a short paper, not just a poster. This is the "things went reasonably well" outcome.

**Rung 3 — a full conference paper.** Only realistic if I get two of the three threats to a comparable level of completeness, with a comparison between them that says something more general about the fabric-sharing tenant-isolation problem. I'm treating this as a bonus, not a plan.

The actual use of this ladder isn't the ambition-signaling, it's the reverse-engineering: for each rung, I wrote down what the evaluation section of that paper would need to contain, then checked what instrumentation and baseline measurements that implies, then checked whether my current build plan actually produces that data. It already caught one gap — I had no plan to measure power draw with anything more rigorous than eyeballing it, and side-channel work without a real power trace setup isn't a result, it's an anecdote. That's now on the build list before I write another line of the attack logic.

I'll revisit this ladder every few weeks rather than treat it as fixed. The point isn't to lock in a target now and white-knuckle toward it — it's to make sure that whatever I end up with in three months is evaluable by someone other than me, instead of a demo that only makes sense if you already trust me.

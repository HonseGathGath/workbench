---
title: Why I'm keeping this log
date: 2026-08-10
tags: meta, process
---

I'm starting a PCD (I'll explain that term properly in the next post) on FPGA and SmartNIC security and performance, and I want a place to write down what I'm actually doing while I'm doing it — not a polished account of it afterward.

A few reasons, roughly in order of how much they actually motivate me day to day:

**I understand things better once I've written them down badly.** There's a specific kind of clarity that shows up only when I try to explain a decision in full sentences instead of leaving it as a vague shape in my head. Half the value of this log will come from the act of writing it, not from anyone reading it.

**Research projects drift, and I want a paper trail of why.** Six months from now I'll have changed my threat model, dropped an approach that seemed promising, or picked up a constraint I didn't know about at the start. Without notes, that history collapses into "I guess I always thought this." With notes, I can actually check.

**It's a public record of the work, which matters for where this is going.** I'm doing this PCD with an eye toward research applications afterward — PhD programs, research groups, maybe a paper or two along the way. A committee or a potential advisor reading a changelog of real decisions, including the wrong turns, tells them more than a tidy CV line ever could. I'd rather be judged on how I think through a problem than on a one-paragraph summary of the outcome.

So here's what this will and won't be. It won't be tutorial content — there are better places to learn what a PCIe TLP is or how AXI handshaking works, and I'm not trying to compete with them. It will be: design decisions as I make them, dead ends and why they were dead ends, notes from papers that changed how I think about something, and periodic stock-takes of where the project actually stands versus where I thought it would be.

The name is deliberately unglamorous. This is a workbench, not a showroom — things get left half-assembled on it, and that's fine. I'd rather have an honest record of the mess than a retroactively cleaned-up story. If you're reading this later and something seems to contradict a post from a few months earlier, that's not an error — that's the project actually moving.

Next up: what "security architecture" needs to mean for this project before I can call anything about it real.

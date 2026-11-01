---
title: Scoping the PCD — what "security architecture" actually has to mean
date: 2026-08-28
tags: security, threat-model, scoping
---

Quick definition since I used the term without explaining it last time: a PCD, in the program I'm doing this under, is a Projet de Conception et Développement — a design-and-build capstone project. Mine is an FPGA-based SmartNIC with a stated goal of "security and performance." I've now spent two weeks finding out how little that sentence actually commits me to, and how much work it takes to turn it into something I can be held to.

"Security architecture" is one of those phrases that sounds load-bearing and is actually empty until you attach it to a threat model. Secure against whom, doing what, with what resources? Without that, every design choice I make is unfalsifiable — I can always claim it was "for security" after the fact. So before writing any RTL, I made myself name the actual threats this project addresses. Three came out of it:

**Side-channel leakage.** An FPGA sharing a fabric with other tenants (or even just other logic from the same tenant) leaks information through power draw and timing — this is well-trodden ground in FPGA security research, but it's real and it's measurable on real hardware, which means it's testable on my timeline. I'm not going to solve side-channel resistance in general; I'm scoping to a specific leakage channel I can instrument and measure before and after a mitigation.

**Bitstream tampering and provenance.** If the bitstream loaded onto the device can be swapped, modified, or replayed without detection, every other security property is moot. This one has a cleaner story — bitstream authentication and integrity checking is a solved-ish problem in the literature, which means my job is to implement it correctly and measure the cost, not invent something new.

**Tenant isolation on shared fabric.** This is the one closest to the "SmartNIC" half of the project. If the device partitions logic or data paths for different tenants (different VMs, different traffic classes), what actually enforces that boundary, and what happens when it's stressed? This is the threat that's hardest to pin down precisely, and I suspect it'll need its own follow-up post once I've read more.

What I'm explicitly *not* claiming: formal verification of isolation properties, resistance to a nation-state adversary with physical access, or anything approaching a general secure-enclave story. Scoping down is uncomfortable because it makes the project sound smaller than "security architecture" implies. But an unscoped security claim is worth zero — reviewers (and advisors, and my own future self six months from now) will ask "secure against what?" and I'd rather have an answer than a vibe.

Practically, this changes what I build next. Instead of starting with the datapath and bolting security on, I'm going to pick one of these three threats — probably bitstream tampering, since it's the most self-contained — and build a minimal, measurable version of an attack and a mitigation for it first. Get one real result before trying to generalize.

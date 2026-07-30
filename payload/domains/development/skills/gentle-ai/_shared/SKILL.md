---
name: _shared
description: "Shared SDD references for installed skills. Not invokable."
disable-model-invocation: true
user-invocable: false
license: MIT
metadata:
  author: gentleman-programming
  version: "1.0"
---

## Purpose

This directory stores shared reference documents consumed by real SDD skills:
`sdd-phase-common.md` (the phase protocol and the canonical return contract) and
`engram-convention.md` (artifact naming and retrieval).
<!-- matecito-ai: `persistence-contract.md` was deleted and dropped from this list. No phase was ever
     instructed to read it, and its content — the engram/none mode table and the "never write project
     files" rules — was a parallel copy of Section C of `sdd-phase-common.md`, which all ten executors do
     read. A duplicated contract leaves N-1 copies stale on every edit. The list is now exhaustive, not
     illustrative: "for example" invited exactly this kind of leftover to sit here unnoticed. -->

## Not Invokable

`_shared` is a support package only. Do not invoke it as a skill.

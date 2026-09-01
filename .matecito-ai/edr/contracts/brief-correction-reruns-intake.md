# EDR — A correction to the brief is resolved by re-running intake, never by the orchestrator editing the brief

- **Status:** Accepted
- **Date:** 2026-09-01

## Contexto
The new Brief Confirmation Gate offers the intake brief for one of two answers: accept, or correct. The behavior on correction is already fixed by the change's spec — re-dispatch `sdd-intake` with the original request plus the user's words — but the reasoning behind picking that path over the orchestrator editing the persisted brief in place has no durable home yet.

## Decisión
A correction to the brief is resolved by re-dispatching `sdd-intake` with the original request plus the user's words, verbatim — never by the orchestrator editing the persisted brief. `sdd-intake` runs a fresh single pass and overwrites its own brief under the same `topic_key`. This preserves the single-writer invariant (only the phase that writes an artifact may rewrite it) and needs no new plumbing: downstream phases retrieve `sdd/{change-name}/intake` at their own dispatch time, so nothing has to forward the accepted version by a second channel.

## Reglas verificables
- **[manual]** A correction at the Brief Confirmation Gate re-dispatches `sdd-intake` with the original request plus the user's words verbatim; the orchestrator never writes into `sdd/{change-name}/intake` directly.
- **[manual]** `sdd-intake` overwrites its own brief under the same `topic_key` on a correction re-dispatch, exactly as on its first pass — no new artifact key, no second channel forwarding the accepted version.

## Alternativas consideradas
The orchestrator edits the persisted brief directly. One fewer dispatch, but it makes the orchestrator write an artifact that only `sdd-intake` writes — new plumbing that breaks the single-writer invariant `contracts/single-writer-per-batch.md` already establishes for this flow. Rejected.

## Consecuencias
A correction costs a full `sdd-intake` re-dispatch instead of an in-place edit, and the correction loop has no bound — it closes only when the user accepts. In exchange, the brief keeps exactly one writer, and no new mechanism is needed to carry the accepted version to downstream phases: they already retrieve the key at their own dispatch time.

## Relacionados
- `relacionado-con` → [single-writer-per-batch.md](single-writer-per-batch.md) — el mismo invariante de un-solo-escritor que esta decisión preserva para el brief.
- `relacionado-con` → [intake-brief-passthrough-sections.md](intake-brief-passthrough-sections.md) — qué secciones lleva el brief que esta decisión re-emite en cada corrección.

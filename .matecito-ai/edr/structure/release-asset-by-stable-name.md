# EDR — The installer selects the release asset by its stable name, qmd.tgz, never one derived from the tag

- **Status:** Accepted
- **Date:** 2026-09-01

## Contexto
The live latest release of the qmd fork (v2.8.3-mate.6) publishes exactly two assets of identical size (274954 bytes): `qmd.tgz`, a version-independent name, and `tobilu-qmd-2.8.3-mate.6.tgz`, the `npm pack` output of the scoped package `@tobilu/qmd`. There is no `checksums.txt` in the release.

## Decisión
qmdLatestTarballURL matches the release asset by the exact, fixed name `qmd.tgz` — never a name built from `tag_name` — and returns the matching asset's `browser_download_url`. When `qmd.tgz` is absent, it errors naming the tag and listing every asset the release did publish.

## Reglas verificables
- **[manual]** qmdLatestTarballURL matches the asset by exact name "qmd.tgz", never a name built from tag_name.
- **[manual]** when qmd.tgz is absent from the release's assets, the helper returns an error naming the tag and listing every available asset name.

## Alternativas consideradas
Deriving `tobilu-qmd-<version>.tgz` from `tag_name`, the way `releasedl` derives its goreleaser asset name. Discarded: the derived name depends on two independently drifting things — the tag→version mapping (`v2.8.3-mate.6` → `2.8.3-mate.6`) and the package scope (`@tobilu`) — and a change to either produces a silent miss that reads as "no release asset". `qmd.tgz` depends on neither.

## Consecuencias
This decision binds the fork's own release process going forward: every future release has to keep publishing an asset literally named `qmd.tgz`, or new installs stop resolving. That assumption cannot be verified from this repo — it is raised as an unverified assumption on the person who owns that release process, not silently absorbed.

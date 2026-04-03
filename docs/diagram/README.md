# Architecture Diagrams

**Path:** `docs/diagram`

## Purpose

This folder stores source diagrams and generated images used to explain architecture and runtime flows.

## Diagram Catalog

Keep this folder for diagrams that clarify:

- entrypoint and bootstrap flow
- event backbone and projection flow
- contract or integration boundaries that are hard to read from code alone

## Editing Workflow

- update the diagram source first
- regenerate or refresh the exported image if the published docs depend on it
- keep README links and image paths aligned with the source files

## Boundary Rules

- diagrams should explain a boundary that already exists in code or docs; they should not invent architecture
- generated images are support artifacts, not the canonical source

## Risks And Compatibility Notes

- stale images are worse than missing images when they contradict the current code path

---

<!-- doc-nav:start -->
## Navigation
- [Back to docs index](../index.md)
<!-- doc-nav:end -->

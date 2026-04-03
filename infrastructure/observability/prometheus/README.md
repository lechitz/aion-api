# Prometheus Scrape Config

**Path:** `infrastructure/observability/prometheus`

## Purpose

This folder stores Prometheus scrape configuration for the local observability stack.

## Current File

| File | Purpose |
| --- | --- |
| scrape config in this folder | target discovery and scrape rules |

## Validate

- start the stack
- confirm the expected scrape targets are healthy
- verify key application metrics appear in Prometheus

## Boundary Rules

- scrape target ownership belongs here; metric instrumentation belongs in application code

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->

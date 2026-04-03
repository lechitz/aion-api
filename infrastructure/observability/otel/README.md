# OpenTelemetry Collector Config

**Path:** `infrastructure/observability/otel`

## Purpose

This folder stores collector config for trace and metric ingestion in the local observability stack.

## Current Pipeline

- receive OTLP signals from services
- normalize and route traces and metrics to downstream backends

## Source File

| File | Purpose |
| --- | --- |
| collector config in this folder | OTLP ingestion and exporter wiring |

## Validate

- start the stack
- hit instrumented endpoints
- confirm traces and metrics reach the downstream backends

## Boundary Rules

- collector routing belongs here; application instrumentation belongs in `internal/platform/observability`

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->

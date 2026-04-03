# Fluent Bit Configuration

**Path:** `infrastructure/observability/fluentbit`

## Purpose

This folder stores log collection and forwarding config for the local observability stack.

## Current Flow

- tail container or runtime logs
- normalize and label records
- forward them to the configured log backend

## Validate

- start the local stack
- produce application logs
- confirm they appear downstream with the expected labels

## Boundary Rules

- keep log parsing and forwarding rules here, not in app code
- if log labels change, update downstream queries and dashboards in the same PR

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->

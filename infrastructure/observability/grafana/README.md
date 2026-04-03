# Grafana Provisioning Assets

**Path:** `infrastructure/observability/grafana`

## Purpose

This folder stores Grafana datasources, dashboards, and provisioning assets used by the local observability stack.

## Current Layout

| Area | Role |
| --- | --- |
| `dashboards/` | dashboard JSON definitions |
| `datasources/` | datasource provisioning |
| `provisioning/` | Grafana provisioning structure |

## Validate

- start the local stack
- open Grafana
- verify that datasources and dashboards provision without manual repair

## Boundary Rules

- keep dashboard truth in versioned assets, not in ad-hoc UI edits
- if a dashboard depends on a new label or metric, update the matching collector or exporter config in the same PR

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->

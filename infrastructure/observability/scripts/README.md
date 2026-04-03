# Observability Helper Scripts

**Path:** `infrastructure/observability/scripts`

## Purpose

This folder stores helper automation for validating and refreshing observability wiring in the local stack.

## Package Composition

- `setup-improvements.sh`
  - validate current observability assets, restart the stack, and check key endpoints plus Grafana provisioning

## Boundary Rules

- helper scripts are operational support, not the source of truth for telemetry topology
- production deployment logic and secrets do not belong here

## Validate

```bash
./infrastructure/observability/scripts/setup-improvements.sh
```

## Risks And Compatibility Notes

- script success is useful but not sufficient; config files remain the canonical wiring definition

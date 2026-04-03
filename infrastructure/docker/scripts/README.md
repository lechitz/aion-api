# Docker Runtime Scripts

**Path:** `infrastructure/docker/scripts`

## Purpose

This folder stores shell scripts used by Docker workflows, Make targets, and container startup.

## Package Composition

- `entrypoint.sh`
  - container entrypoint for the API image

## Boundary Rules

- keep Docker concerns isolated from application code
- make destructive actions explicit and opt-in
- do not embed secrets in scripts

## Validate

- build the image
- start the stack
- confirm the container boots through `entrypoint.sh` without manual patching

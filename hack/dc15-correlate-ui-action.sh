#!/usr/bin/env bash
set -euo pipefail

# Correlates draft UI actions across dashboard, API, and chat logs.
# Usage:
#   ./hack/dc15-correlate-ui-action.sh --since 45m
#   ./hack/dc15-correlate-ui-action.sh --since 45m --draft-id tag-abc123

SINCE="45m"
DRAFT_ID=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --since)
      SINCE="${2:-}"
      shift 2
      ;;
    --draft-id)
      DRAFT_ID="${2:-}"
      shift 2
      ;;
    -h|--help)
      cat <<'EOF'
Usage:
  ./hack/dc15-correlate-ui-action.sh [--since <duration>] [--draft-id <id>]

Options:
  --since      Docker logs window (default: 45m)
  --draft-id   Optional filter for one specific draft id
EOF
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      exit 1
      ;;
  esac
done

if ! command -v docker >/dev/null 2>&1; then
  echo "docker not found in PATH" >&2
  exit 1
fi

declare -A PATTERNS
PATTERNS["aion-dev-web"]="draft_id|draft_accept|draft_cancel|action_result|/chat/text|VITE"
PATTERNS["aion-dev-api"]="HTTP chat request includes UI action|ui_action_type|draft_id|consent_required|consent_confirmed|consent_policy_version|Chat request cancelled"
PATTERNS["aion-dev-chat"]="UI action metadata detected|Handling UI action|UI action handled|draft_id|status=|High-risk action blocked"

CONTAINERS=("aion-dev-web" "aion-dev-api" "aion-dev-chat")

print_header() {
  local title="$1"
  printf '\n===== %s =====\n' "$title"
}

container_exists() {
  local name="$1"
  docker ps -a --format '{{.Names}}' | rg -x --quiet "$name"
}

for container in "${CONTAINERS[@]}"; do
  print_header "$container (since=${SINCE})"

  if ! container_exists "$container"; then
    echo "SKIP: container not found."
    continue
  fi

  base_pattern="${PATTERNS[$container]}"
  lines="$(docker logs --since "$SINCE" "$container" 2>&1 | rg -n "$base_pattern" || true)"

  if [[ -n "$DRAFT_ID" ]]; then
    lines="$(printf '%s\n' "$lines" | rg -n "$DRAFT_ID|draft_id|status=|ui_action_type" || true)"
  fi

  if [[ -z "${lines// }" ]]; then
    echo "No matching log lines."
    continue
  fi

  printf '%s\n' "$lines"
done

#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
CONTRACTS_DIR="${ROOT_DIR}/contracts/graphql"
OUTPUT_FILE="${CONTRACTS_DIR}/manifest.json"

if [[ ! -d "${CONTRACTS_DIR}" ]]; then
  echo "contracts/graphql not found: ${CONTRACTS_DIR}" >&2
  exit 1
fi

TMP_ENTRIES="$(mktemp)"
trap 'rm -f "${TMP_ENTRIES}"' EXIT

while IFS= read -r abs_file; do
  rel_file="${abs_file#${ROOT_DIR}/}"

  kind_raw="$(sed -E 's/#.*$//' "${abs_file}" | tr '\n' ' ' | sed -E 's/^[[:space:]]*([A-Za-z]+).*/\1/' | tr '[:upper:]' '[:lower:]')"
  case "${kind_raw}" in
    query) kind="query" ;;
    mutation) kind="mutation" ;;
    *)
      echo "Unknown GraphQL operation type in ${rel_file}" >&2
      exit 1
      ;;
  esac

  operation_name="$(sed -E 's/#.*$//' "${abs_file}" | tr '\n' ' ' | sed -E 's/^[[:space:]]*(query|mutation)[[:space:]]+([A-Za-z_][A-Za-z0-9_]*).*/\2/')"
  root_field="$(tr '\n' ' ' < "${abs_file}" | sed -E 's/[[:space:]]+/ /g' | sed -E 's/^[^{]*\{[[:space:]]*//' | sed -E 's/^([A-Za-z_][A-Za-z0-9_]*).*/\1/')"

  if [[ -z "${operation_name}" || -z "${root_field}" ]]; then
    echo "Failed to parse operation metadata in ${rel_file}" >&2
    exit 1
  fi

  checksum="$(sha256sum "${abs_file}" | awk '{print $1}')"
  printf '%s\t%s\t%s\t%s\t%s\n' "${kind}" "${operation_name}" "${root_field}" "${rel_file}" "${checksum}" >> "${TMP_ENTRIES}"
done < <(find "${CONTRACTS_DIR}/queries" "${CONTRACTS_DIR}/mutations" -type f -name '*.graphql' | LC_ALL=C sort)

{
  printf '{\n'
  printf '  "version": 1,\n'
  printf '  "generatedAt": "__GENERATED_AT__",\n'
  printf '  "operations": [\n'

  total="$(wc -l < "${TMP_ENTRIES}" | tr -d ' ')"
  idx=0
  while IFS=$'\t' read -r kind operation_name root_field rel_file checksum; do
    idx=$((idx + 1))
    comma=","
    if [[ ${idx} -eq ${total} ]]; then
      comma=""
    fi

    printf '    {"type":"%s","name":"%s","rootField":"%s","path":"%s","sha256":"%s"}%s\n' \
      "${kind}" "${operation_name}" "${root_field}" "${rel_file}" "${checksum}" "${comma}"
  done < "${TMP_ENTRIES}"

  printf '  ]\n'
  printf '}\n'
} > "${OUTPUT_FILE}"

# Keep output deterministic for drift checks.
sed -i 's/"generatedAt": "__GENERATED_AT__"/"generatedAt": "deterministic"/' "${OUTPUT_FILE}"

echo "✅ GraphQL manifest generated: ${OUTPUT_FILE}"

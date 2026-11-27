#!/usr/bin/env bash

# Tags a commit for a container image release.
#
# The tag name is the version the pipeline bakes into every image tag, so it
# must be <major>.<minor>.<patch> optionally followed by a pre-release suffix:
#   1.1.0        -> images tagged 1.1.0 and latest
#   1.2.0-rc.1   -> images tagged 1.2.0-rc.1 only
#
# Usage:
#   ./scripts/release-version.sh tag 1.1.0

set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  ./scripts/release-version.sh tag <major.minor.patch>[-<pre-release>]

Examples:
  ./scripts/release-version.sh tag 1.2.0
  ./scripts/release-version.sh tag 1.2.0-rc.1

The script only creates the tag locally; push it to trigger the pipeline:
  git push origin <tag>
EOF
}

tag_release() {
  local version="$1"

  if [[ ! "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.]+)?$ ]]; then
    echo "error: '$version' is not <major>.<minor>.<patch>[-<pre-release>]" >&2
    exit 1
  fi

  local tag="v${version}"

  if git rev-parse --verify --quiet "refs/tags/${tag}" >/dev/null; then
    echo "error: tag ${tag} already exists" >&2
    exit 1
  fi

  git tag --annotate "${tag}" --message "Release ${version}"
  echo "Created ${tag} on $(git rev-parse --short HEAD)."
  echo "Push it with: git push origin ${tag}"
}

if [[ $# -lt 1 ]]; then
  usage
  exit 1
fi

case "$1" in
  tag)
    [[ $# -eq 2 ]] || { usage; exit 1; }
    tag_release "$2"
    ;;
  *)
    usage
    exit 1
    ;;
esac

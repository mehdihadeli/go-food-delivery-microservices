#!/usr/bin/env bash

# Calculates the container image version for the current commit.
#
# The repository has no root go.mod, so this is not a Go module release
# version: it only names the images built from this commit.
#   - a pushed tag vX.Y.Z (or vX.Y.Z-<suffix>) is a release, used verbatim
#   - anything else is a preview build derived from the latest release tag
#
# Outputs: version, docker_tag, environment, commit.
#   environment=production  a released vX.Y.Z tag, the only run that may publish
#   environment=staging     a pre-release tag or a branch preview, images only
#   environment=none        nothing to publish

set -euo pipefail

emit() {
  if [ -n "${GITHUB_OUTPUT:-}" ]; then
    printf '%s=%s\n' "$1" "$2" >>"$GITHUB_OUTPUT"
  else
    printf '%s=%s\n' "$1" "$2"
  fi
}

die() {
  printf 'error: %s\n' "$1" >&2
  exit 1
}

# Strips the version out of a tag, tolerating the <module directory>/vX.Y.Z
# form that subdirectory Go modules use.
tag_version() {
  local tag="$1"
  tag="${tag##*/}"
  [[ "$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.]+)?$ ]] ||
    die "tag '$1' is not vX.Y.Z or vX.Y.Z-<suffix>"
  printf '%s' "${tag#v}"
}

readonly commit="$(git rev-parse --short HEAD)"

if [[ "${GITHUB_REF:-}" == refs/tags/* ]]; then
  version="$(tag_version "${GITHUB_REF_NAME:?GITHUB_REF_NAME must be set}")"
else
  # Branch builds are previews of the next release: they sort above the last
  # release and never get the mutable :latest tag.
  if latest="$(git describe --tags --abbrev=0 --match 'v*' 2>/dev/null)"; then
    latest="${latest##*/}"
    latest="${latest#v}"

    # A released version is done, so its previews belong to the next patch. A
    # pre-release is still in flight, so keep its version and only add the
    # preview marker, which sorts below the pre-release itself.
    base="${latest%%-*}"
    if [[ "$latest" != *-* ]]; then
      major="${base%%.*}"
      rest="${base#*.}"
      minor="${rest%%.*}"
      patch="${rest#*.}"
      base="${major}.${minor}.$((patch + 1))"
    fi
  else
    base="0.0.0"
  fi
  version="${base}-preview.$(git rev-list --count HEAD).${commit}"
fi


if [[ "$version" == *-* ]]; then
  environment=staging
else
  environment=production
fi

printf 'version:     %s\n' "$version"
printf 'environment: %s\n' "$environment"
printf 'commit:      %s\n' "$commit"

emit "version" "$version"
# The image tag is the version verbatim: one commit, one tag, no build metadata
# suffix that would make two different images look identical.
emit "docker_tag" "$version"
emit "environment" "$environment"
emit "commit" "$commit"

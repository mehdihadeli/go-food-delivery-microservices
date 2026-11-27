#!/bin/bash

# Single source of truth for the module list the mise tasks iterate over.
#
#   ./scripts/modules.sh            every module, shared package first
#   ./scripts/modules.sh services   the microservices only
#   ./scripts/modules.sh pkg        the shared module only
#
# Order matters for `mise run build`: every service module replaces
# .../internal/pkg with the local ../../pkg directory, so the shared module is
# built first.
#
# Adding a microservice means adding one word here and one matrix entry in
# .github/workflows/ci.yml. Nothing else in the repository lists the modules.
set -e

readonly selection="${1:-all}"

case "$selection" in
all) printf '%s\n' pkg catalogwriteservice catalogreadservice orderservice ;;
services) printf '%s\n' catalogwriteservice catalogreadservice orderservice ;;
pkg) printf '%s\n' pkg ;;
*)
    echo "unknown selection: $selection (expected all, services or pkg)" >&2
    exit 1
    ;;
esac

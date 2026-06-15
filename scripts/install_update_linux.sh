#!/bin/bash
# Backward-compatible wrapper. Prefer: plugin/update.sh update
exec "$(cd "$(dirname "$0")/.." && pwd)/plugin/update.sh" update "$@"

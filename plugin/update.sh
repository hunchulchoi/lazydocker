#!/usr/bin/env bash
# Lazydocker update plugin — separate from the lazydocker binary.
# Updates the upstream release binary; fork config/features stay independent.
#
# Usage:
#   plugin/update.sh check
#   plugin/update.sh update
#   DIR=~/.local/bin plugin/update.sh update
#
# Env:
#   LAZYDOCKER_REPO  GitHub repo (default: jesseduffield/lazydocker)
#   DIR              Install dir (default: ~/.local/bin)

set -euo pipefail

LAZYDOCKER_REPO="${LAZYDOCKER_REPO:-jesseduffield/lazydocker}"
DIR="${DIR:-$HOME/.local/bin}"
GITHUB_API="https://api.github.com/repos/${LAZYDOCKER_REPO}/releases/latest"

usage() {
  cat <<EOF
Usage: $(basename "$0") check|update

  check   Compare installed lazydocker with latest GitHub release
  update  Download and install latest release binary

Linux and macOS only. Windows not supported.

Environment:
  LAZYDOCKER_REPO  default: jesseduffield/lazydocker
  DIR              install directory, default: ~/.local/bin
EOF
}

require_unix() {
  case "$(uname -s)" in
    Linux|Darwin) ;;
    *)
      echo "error: self-update plugin supports Linux and macOS only" >&2
      exit 1
      ;;
  esac
}

map_arch() {
  local arch
  arch="$(uname -m)"
  case "$arch" in
    x86_64|amd64) echo x86_64 ;;
    i386|i686) echo x86 ;;
    arm64|aarch64) echo arm64 ;;
    armv6*) echo armv6 ;;
    armv7*) echo armv7 ;;
    arm) echo armv7 ;;
    *)
      echo "error: unsupported architecture: $arch" >&2
      exit 1
      ;;
  esac
}

latest_tag() {
  curl -fsSL -H 'Accept: application/json' "$GITHUB_API" \
    | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' \
    | head -1
}

current_version() {
  if ! command -v lazydocker >/dev/null 2>&1; then
    echo "not installed"
    return
  fi
  lazydocker --version 2>/dev/null | head -1 | tr -d '\r'
}

normalize_tag() {
  local v="$1"
  v="${v#v}"
  v="${v#Version: }"
  v="${v%% *}"
  echo "$v"
}

asset_name() {
  local tag="$1"
  local version os arch
  version="$(normalize_tag "$tag")"
  os="$(uname -s)"
  arch="$(map_arch)"
  echo "lazydocker_${version}_${os}_${arch}.tar.gz"
}

download_url() {
  local tag="$1"
  echo "https://github.com/${LAZYDOCKER_REPO}/releases/download/${tag}/$(asset_name "$tag")"
}

asset_exists() {
  local url="$1"
  local code
  code="$(curl -fsSL -o /dev/null -w '%{http_code}' -I "$url" || true)"
  [ "$code" = "200" ]
}

cmd_check() {
  require_unix
  local latest current
  latest="$(latest_tag)"
  if [ -z "$latest" ]; then
    echo "error: could not read latest release from ${LAZYDOCKER_REPO}" >&2
    exit 1
  fi

  current="$(current_version)"
  echo "repo:    ${LAZYDOCKER_REPO}"
  echo "current: ${current:-unknown}"
  echo "latest:  ${latest}"

  local cur_norm lat_norm
  cur_norm="$(normalize_tag "$current")"
  lat_norm="$(normalize_tag "$latest")"

  if [ "$cur_norm" = "$lat_norm" ]; then
    echo "status:  up to date"
    exit 0
  fi

  echo "status:  update available"
  echo "run:     $(dirname "$0")/update.sh update"
  exit 2
}

cmd_update() {
  require_unix
  local tag url tmpdir archive dest
  tag="$(latest_tag)"
  if [ -z "$tag" ]; then
    echo "error: could not read latest release" >&2
    exit 1
  fi

  url="$(download_url "$tag")"
  if ! asset_exists "$url"; then
    echo "error: release asset not found: $url" >&2
    exit 1
  fi

  mkdir -p "$DIR"
  tmpdir="$(mktemp -d)"
  trap 'rm -rf "$tmpdir"' EXIT
  archive="${tmpdir}/lazydocker.tar.gz"

  echo "downloading ${tag} ..."
  curl -fsSL -o "$archive" "$url"
  tar xzf "$archive" -C "$tmpdir" lazydocker
  dest="${DIR}/lazydocker"
  install -m 755 "${tmpdir}/lazydocker" "$dest"

  echo "installed: ${dest} (${tag})"
  "${dest}" --version 2>/dev/null | head -1 || true
}

main() {
  case "${1:-}" in
    check) cmd_check ;;
    update|install) cmd_update ;;
    -h|--help|help|"") usage ;;
    *)
      echo "error: unknown command: $1" >&2
      usage
      exit 1
      ;;
  esac
}

main "$@"

#!/usr/bin/env bash
# Installer for matecito-ai.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/franwerner/matecito-ai/master/scripts/install.sh | bash
#
# Environment variables:
#   VERSION       Release tag to install (default: latest). Example: VERSION=v0.1.0
#   INSTALL_DIR   Where to place the binary (default: $HOME/.local/bin).
#
set -euo pipefail

REPO="franwerner/matecito-ai"
BINARY="matecito-ai"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
VERSION="${VERSION:-latest}"

err() { printf "error: %s\n" "$*" >&2; exit 1; }
info() { printf "==> %s\n" "$*"; }

detect_os() {
  local os
  os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  case "$os" in
    linux)  echo "linux" ;;
    darwin) echo "darwin" ;;
    *) err "OS no soportado por este script: $os. Bajá el binario manualmente desde https://github.com/$REPO/releases" ;;
  esac
}

detect_arch() {
  local arch
  arch="$(uname -m)"
  case "$arch" in
    x86_64|amd64) echo "amd64" ;;
    arm64|aarch64) echo "arm64" ;;
    *) err "arquitectura no soportada: $arch" ;;
  esac
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || err "se necesita '$1' instalado"
}

resolve_version() {
  if [ "$VERSION" = "latest" ]; then
    require_cmd curl
    VERSION="$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
      | grep -oE '"tag_name"[[:space:]]*:[[:space:]]*"[^"]+"' \
      | head -n1 \
      | sed -E 's/.*"([^"]+)"$/\1/')"
    [ -n "$VERSION" ] || err "no pude resolver la última release"
  fi
}

verify_sha256() {
  local file="$1" sums="$2"
  local expected name actual
  name="$(basename "$file")"
  expected="$(grep " $name\$" "$sums" | awk '{print $1}' || true)"
  [ -n "$expected" ] || err "checksum para $name no encontrado en checksums.txt"
  if command -v sha256sum >/dev/null 2>&1; then
    actual="$(sha256sum "$file" | awk '{print $1}')"
  elif command -v shasum >/dev/null 2>&1; then
    actual="$(shasum -a 256 "$file" | awk '{print $1}')"
  else
    err "se necesita 'sha256sum' o 'shasum' instalado"
  fi
  [ "$expected" = "$actual" ] || err "checksum no coincide: esperado $expected, obtenido $actual"
}

main() {
  require_cmd curl
  require_cmd tar

  local os arch version_no_v archive tmp
  os="$(detect_os)"
  arch="$(detect_arch)"
  resolve_version
  version_no_v="${VERSION#v}"
  archive="${BINARY}_${version_no_v}_${os}_${arch}.tar.gz"

  info "Descargando $BINARY $VERSION para $os/$arch"

  tmp="$(mktemp -d)"
  trap "rm -rf '$tmp'" EXIT

  curl -fsSL "https://github.com/$REPO/releases/download/$VERSION/$archive" -o "$tmp/$archive"
  curl -fsSL "https://github.com/$REPO/releases/download/$VERSION/checksums.txt" -o "$tmp/checksums.txt"

  info "Verificando checksum SHA256"
  verify_sha256 "$tmp/$archive" "$tmp/checksums.txt"

  info "Extrayendo binario"
  tar -xzf "$tmp/$archive" -C "$tmp"
  [ -f "$tmp/$BINARY" ] || err "no encontré $BINARY dentro del archivo"

  mkdir -p "$INSTALL_DIR"
  install -m 0755 "$tmp/$BINARY" "$INSTALL_DIR/$BINARY"

  info "Instalado en $INSTALL_DIR/$BINARY"

  case ":$PATH:" in
    *":$INSTALL_DIR:"*)
      info "Listo. Probá: $BINARY verify"
      ;;
    *)
      cat <<EOF

⚠ $INSTALL_DIR no está en tu PATH.

Agregalo al rc de tu shell (~/.bashrc, ~/.zshrc, etc.):
  export PATH="$INSTALL_DIR:\$PATH"

Después reabrí la shell o corré: source ~/.bashrc (o el rc que corresponda).
EOF
      ;;
  esac
}

main "$@"

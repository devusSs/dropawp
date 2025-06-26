#!/usr/bin/env bash
set -euo pipefail

PKG_VERSION=${PKG_VERSION:-$(git describe --tags --abbrev=0 2>/dev/null || echo "dev")}
PKG_COMMIT=${PKG_COMMIT:-$(git rev-parse --short HEAD || echo "none")}
PKG_DATE=${PKG_DATE:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}

echo "→ Building dropawp"
echo "   Version: ${PKG_VERSION}"
echo "   Commit:  ${PKG_COMMIT}"
echo "   Date:    ${PKG_DATE}"

LDFLAGS=(
  "-X github.com/devusSs/dropawp/internal/version.Version=${PKG_VERSION}"
  "-X github.com/devusSs/dropawp/internal/version.Commit=${PKG_COMMIT}"
  "-X github.com/devusSs/dropawp/internal/version.Date=${PKG_DATE}"
)

go build -v -ldflags "${LDFLAGS[*]}" .
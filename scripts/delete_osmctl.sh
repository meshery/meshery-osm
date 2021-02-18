#!/bin/sh

set -e

: "${OSM_VERSION:=}"
: "${OSM_ARCH:=amd64}"
: "${NAMESPACE:=osm-system}"
: "${OS:=$(uname | awk '{print tolower($0)}')}"
URL="https://github.com/openservicemesh/osm/releases/download/$OSM_VERSION/osm-$OSM_VERSION-$OS-$OSM_ARCH.tar.gz"

if ! curl -L "$URL" | tar xz; then
  exit 1
fi

if ! ./$OS-$OSM_ARCH/osm mesh -n ${NAMESPACE} delete -f; then
	exit 2
fi


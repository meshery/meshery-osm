#!/bin/sh

set -e

: "${OSM_VERSION:=}"
: "${OSM_ARCH:=amd64}"
: "${OS:=$(uname | awk '{print tolower($0)}')}"
URL="https://github.com/openservicemesh/osm/releases/download/$OSM_VERSION/osm-$OSM_VERSION-$OS-$OSM_ARCH.tar.gz"

if ! type "grep" > /dev/null 2>&1; then
  exit 1;
fi
if ! type "curl" > /dev/null 2>&1; then
  exit 2;
fi
if ! type "tar" > /dev/null 2>&1; then
  exit 3;
fi
if ! type "gzip" > /dev/null 2>&1; then
  exit 4;
fi

if ! curl -s --head $URL | head -n 1 | grep "HTTP/1.[01] [23].." > /dev/null; then
  exit 5;
fi


if ! curl -L "$URL" | tar xz; then
  exit 6;
fi

if ! ./$OS-$OSM_ARCH/osm install; then
	exit 7;
fi


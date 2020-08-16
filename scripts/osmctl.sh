#!/bin/sh

set -e

: "${OSM_VERSION:=}"
: "${OSM_ARCH:=amd64}"
: "${OS:=$(uname | awk '{print tolower($0)}')}"
URL="https://github.com/openservicemesh/osm/releases/download/$OSM_VERSION/osm-$OSM_VERSION-$OS-$OSM_ARCH.tar.gz"

printf "\n"
printf "INFO\tWelcome to the Open service mesh automated download!\n"

if ! type "grep" > /dev/null 2>&1; then
  printf "ERROR\tgrep cannot be found\n"
  exit 1;
fi
if ! type "curl" > /dev/null 2>&1; then
  printf "ERROR\tcurl cannot be found\n"
  exit 1;
fi
if ! type "tar" > /dev/null 2>&1; then
  printf "ERROR\ttar cannot be found\n"
  exit 1;
fi
if ! type "gzip" > /dev/null 2>&1; then
  printf "ERROR\tgzip cannot be found\n"
  exit 1;
fi

printf "INFO\tOSM version: %s\n" "$OSM_VERSION"
printf "INFO\tOperating system: %s\n" "$OS"
printf "INFO\tOperating system architecture: %s\n" "$OSM_ARCH"


if ! curl -s --head $URL | head -n 1 | grep "HTTP/1.[01] [23].." > /dev/null; then
  printf "ERROR\tUnable to download OSM at the following URL: %s\n" "$URL"
  exit 1
fi

printf "INFO\tDownloading osmctl from: %s" "$URL"
printf "\n\n"

if curl -L "$URL" | tar xz; then
  printf "\n"
  printf "INFO\tosmctl %s has been downloaded!\n" "$OSM_VERSION"
  printf "\n"
else
  printf "\n"
  printf "ERROR\tUnable to download osmctl\n"
  exit 1
fi

printf "INFO\tStarting deployment......\n"
if ! ./$OS-$OSM_ARCH/osm install; then
	printf "ERROR\tUnable to deploy\n"
	exit 1
fi
printf "INFO\tDeployment successfull!!\n"

if ! rm -rf $OS-$OSM_ARCH; then
	printf "ERROR\tUnable to clear temperory files!"
fi

printf "INFO\tOpen service mesh has been installed successfully!!\n"

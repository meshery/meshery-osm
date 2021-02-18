#!/bin/sh

set -e

: "${OSM_VERSION:=}"
: "${OSM_ARCH:=amd64}"
: "${OS:=$(uname | awk '{print tolower($0)}')}"

URL="https://github.com/openservicemesh/osm/archive/${OSM_VERSION}.tar.gz"

# if ! curl -L "$URL" | tar xz; then
# 	exit 1
# fi

# ./$OS-$OSM_ARCH/osm
#cd $PWD/osm-"${OSM_VERSION[@]:1}"
#cat .env.example > .env
# ./demo/run-osm-demo.sh

CTR_REGISTRY="osmci.azurecr.io/osm"

function install() {
	for i in bookstore bookbuyer bookthief bookwarehouse; do kubectl create ns $i; done

	./$OS-$OSM_ARCH/osm namespace add bookstore bookbuyer bookthief bookwarehouse

	kubectl apply -f "$PWD/osm-${OSM_VERSION[@]:1}/docs/example/manifests/apps/"
}

function delete() {
	for i in bookstore bookbuyer bookthief bookwarehouse; do kubectl delete ns $i; done
}

if [ "$1" = "install" ]; then
	echo "installing sample app"
	install
elif [ "$1" = "delete" ]; then
	echo "deleting sample app"
	delete
else
	exit 1
fi

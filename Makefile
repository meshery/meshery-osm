

include build/Makefile.show-help.mk

GOPATH = $(shell go env GOPATH)

## Run Meshery Adapter.
run:
	go mod tidy; \
	DEBUG=true GOPROXY=direct GOSUMDB=off go run main.go

protoc-setup:
	cd meshes
	wget https://raw.githubusercontent.com/layer5io/meshery/master/meshes/meshops.proto

check: error
	golangci-lint run

check-clean-cache:
	golangci-lint cache clean

proto:	
	protoc -I meshes/ meshes/meshops.proto --go_out=plugins=grpc:./meshes/

## Build Docker Image for Meshery Adapter.
docker:
	docker build -t layer5/meshery-osm .

## Build and Run Docker Image for Meshery Adapter.
docker-run:
	(docker rm -f meshery-osm) || true
	docker run --name meshery-osm -d \
	-p 10009:10009 \
	-e DEBUG=true \
	layer5/meshery-osm

# setup-adapter sets up a new adapter with the given name & port
setup-adapter:
	mv "osm" ${ADAPTER}
	find . -type f -exec sed -i '' -e 's/osm/${ADAPTER}/g' {} +
	find . -type f -exec sed -i '' -e 's/<port>/${PORT}/g' {} +
	find . -type f -exec sed -i '' -e 's/<go_version>/${GO_VERSION}/g' {} +

## Run all Go checks.
go-all: go-tidy go-fmt go-vet golangci-lint

go-fmt:
	go fmt ./...

go-vet:
	go vet ./...

go-tidy:
	@echo "Executing go mod tidy"
	go mod tidy

golangci-lint: $(GOLANGCILINT)
	@echo
	$(GOPATH)/bin/golangci-lint run

error:
	go run github.com/layer5io/meshkit/cmd/errorutil -d . update -i ./helpers -o ./helpers

$(GOLANGCILINT):
	(cd /; GO111MODULE=on GOPROXY="direct" GOSUMDB=off go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.30.0)

.PHONY: error golangci-lint tidy go-vet go-fmt go-all
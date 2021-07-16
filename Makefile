
GOPATH = $(shell go env GOPATH)

protoc-setup:
	cd meshes
	wget https://raw.githubusercontent.com/layer5io/meshery/master/meshes/meshops.proto

check: error
	golangci-lint run

check-clean-cache:
	golangci-lint cache clean

proto:	
	protoc -I meshes/ meshes/meshops.proto --go_out=plugins=grpc:./meshes/

docker:
	docker build -t layer5/meshery-osm .

docker-run:
	(docker rm -f meshery-osm) || true
	docker run --name meshery-osm -d \
	-p 10009:10009 \
	-e DEBUG=true \
	layer5/meshery-osm

run:
	DEBUG=true GOPROXY=direct GOSUMDB=off go run main.go

# setup-adapter sets up a new adapter with the given name & port
setup-adapter:
	mv "osm" ${ADAPTER}
	find . -type f -exec sed -i '' -e 's/osm/${ADAPTER}/g' {} +
	find . -type f -exec sed -i '' -e 's/<port>/${PORT}/g' {} +
	find . -type f -exec sed -i '' -e 's/<go_version>/${GO_VERSION}/g' {} +

.PHONY: local-check
local-check: tidy
local-check: go-fmt
local-check: go-vet
local-check: golangci-lint

.PHONY: go-fmt
go-fmt:
	go fmt ./...

.PHONY: go-vet
go-vet:
	go vet ./...

.PHONY: tidy
tidy:
	@echo "Executing go mod tidy"
	go mod tidy

.PHONY: golangci-lint
golangci-lint: $(GOLANGCILINT)
	@echo
	$(GOPATH)/bin/golangci-lint run

.PHONY: error
error:
	go run github.com/layer5io/meshkit/cmd/errorutil -d . update -i ./helpers -o ./helpers

$(GOLANGCILINT):
	(cd /; GO111MODULE=on GOPROXY="direct" GOSUMDB=off go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.30.0)

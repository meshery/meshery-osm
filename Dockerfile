# syntax=docker/dockerfile:1
FROM golang:1.19 as build-env
ARG VERSION
ARG GIT_COMMITSHA
WORKDIR /github.com/layer5io/meshery-osm
COPY go.mod go.sum ./
RUN go mod download
COPY main.go main.go
COPY internal/ internal/
COPY osm/ osm/
COPY build/ build/
RUN GOPROXY=https://proxy.golang.org,direct CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -ldflags="-w -s -X main.version=$VERSION -X main.gitsha=$GIT_COMMITSHA" -a -o meshery-osm main.go
FROM gcr.io/distroless/nodejs:16
ENV DISTRO="debian"
ENV SERVICE_ADDR="meshery-osm"
ENV MESHERY_SERVER="http://meshery:9081"
WORKDIR /
COPY templates/ ./templates
COPY --from=build-env /github.com/layer5io/meshery-osm/meshery-osm .
ENTRYPOINT ["/meshery-osm"]

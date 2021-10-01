FROM golang:1.15 as build-env
ARG VERSION
ARG GIT_COMMITSHA
WORKDIR /github.com/layer5io/meshery-osm
COPY go.mod go.sum ./
RUN go mod download
COPY main.go main.go
COPY internal/ internal/
COPY osm/ osm/
RUN GOPROXY=https://proxy.golang.org,direct CGO_ENABLED=1 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -ldflags="-w -s -X main.version=$VERSION -X main.gitsha=$GIT_COMMITSHA" -a -o meshery-osm main.go

FROM alpine:3.14 as jsonschema-util
RUN apk add --no-cache curl
WORKDIR /
RUN curl -LO https://github.com/layer5io/kubeopenapi-jsonschema/releases/download/v0.1.0/kubeopenapi-jsonschema

FROM gcr.io/distroless/nodejs:14
ENV DISTRO="debian"
ENV GOARCH="amd64"
ENV SERVICE_ADDR="meshery-osm"
ENV MESHERY_SERVER="http://meshery:9081"
WORKDIR /
COPY templates/ ./templates
COPY --from=build-env /github.com/layer5io/meshery-osm/meshery-osm .
COPY --from=jsonschema-util /kubeopenapi-jsonschema /root/.meshery/bin/kubeopenapi-jsonschema
ENTRYPOINT ["/meshery-osm"]

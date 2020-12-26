FROM golang:1.15 as build-env
WORKDIR /github.com/layer5io/meshery-osm
COPY go.mod go.sum ./
RUN go mod download
COPY main.go main.go
COPY internal/ internal/
COPY osm/ osm/
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -ldflags="-w -s" -a -o meshery-osm main.go

FROM gcr.io/distroless/base
ENV DISTRO="debian"
ENV GOARCH="amd64"
WORKDIR /
COPY --from=build-env /github.com/layer5io/meshery-osm/meshery-osm .
ENTRYPOINT ["/meshery-osm"]
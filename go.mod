module github.com/layer5io/meshery-osm

go 1.15

replace (
	github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200806180306-b7e46afd657f
	github.com/layer5io/meshery-adapter-library v0.1.11 => ../meshery-adapter-library
	github.com/layer5io/meshkit v0.1.31 => ../meshkit
)

require (
	github.com/layer5io/meshery-adapter-library v0.1.11
	github.com/layer5io/meshkit v0.1.31
	github.com/layer5io/service-mesh-performance v0.3.2
)

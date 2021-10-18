module github.com/layer5io/meshery-osm

go 1.15

replace github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200806180306-b7e46afd657f

require (
	github.com/layer5io/meshery-adapter-library v0.1.23
	github.com/layer5io/meshkit v0.2.31
	github.com/layer5io/service-mesh-performance v0.3.3
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.21.0
)

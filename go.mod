module github.com/layer5io/meshery-osm

go 1.15

replace github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200806180306-b7e46afd657f

require (
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/layer5io/meshery-adapter-library v0.1.20
	github.com/layer5io/meshkit v0.2.24
	github.com/layer5io/service-mesh-performance v0.3.3
	github.com/pkg/errors v0.9.1
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.18.12
)

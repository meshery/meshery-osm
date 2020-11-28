module github.com/layer5io/meshery-osm

go 1.15

replace github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200806180306-b7e46afd657f

require (
	github.com/layer5io/meshery-adapter-library v0.1.7
	github.com/layer5io/meshkit v0.1.27
	k8s.io/api v0.18.12
	k8s.io/apimachinery v0.18.12
)

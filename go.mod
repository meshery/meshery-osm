module github.com/layer5io/meshery-osm

go 1.13

replace (
	github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200806180306-b7e46afd657f
	vbom.ml/util => github.com/fvbommel/util v0.0.0-20180919145318-efcd4e0f9787
)

require (
	github.com/Azure/go-autorest/autorest/adal v0.9.0 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/gophercloud/gophercloud v0.4.0 // indirect
	github.com/layer5io/gokit v0.1.12
	github.com/layer5io/learn-layer5/smi-conformance v0.0.0-20200916172547-79cf11334bd7
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	google.golang.org/grpc v1.31.0
	gopkg.in/yaml.v2 v2.3.0
	helm.sh/helm/v3 v3.3.4
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8
	rsc.io/letsencrypt v0.0.3 // indirect
)

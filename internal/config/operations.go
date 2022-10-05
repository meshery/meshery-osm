package config

import (
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/meshes"
	"github.com/layer5io/meshkit/utils"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var (
	OSMOperation          = strings.ToLower(smp.ServiceMesh_OPEN_SERVICE_MESH.Enum().String())
	OSMBookStoreOperation = "osm_bookstore_app"
	ServiceName           = "service_name"
)

func getOperations(op adapter.Operations) adapter.Operations {
	var adapterVersions []adapter.Version
	versions, _ := utils.GetLatestReleaseTagsSorted("openservicemesh", "osm")
	for _, v := range versions {
		adapterVersions = append(adapterVersions, adapter.Version(v))
	}
	op[OSMOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_INSTALL),
		Description: "Open Service Mesh",
		Versions:    adapterVersions,
		Templates: []adapter.Template{
			"templates/open_service_mesh.yaml",
		},
	}

	op[OSMBookStoreOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_SAMPLE_APPLICATION),
		Description: "Bookstore Application",
		Versions:    adapterVersions,
		Templates: []adapter.Template{
			"https://raw.githubusercontent.com/openservicemesh/osm-docs/release-v1.2/manifests/apps/bookbuyer.yaml",
			"https://raw.githubusercontent.com/openservicemesh/osm-docs/release-v1.2/manifests/apps/bookstore-v2.yaml",
			"https://raw.githubusercontent.com/openservicemesh/osm-docs/release-v1.2/manifests/apps/bookthief.yaml",
			"https://raw.githubusercontent.com/openservicemesh/osm-docs/release-v1.2/manifests/apps/bookwarehouse.yaml",
			"https://raw.githubusercontent.com/openservicemesh/osm-docs/release-v1.2/manifests/split/traffic-split-v2.yaml",
			"https://raw.githubusercontent.com/openservicemesh/osm-docs/main/manifests/apps/mysql.yaml",
			"file://templates/osm-bookstore-traffic-access-v1.yaml",
		},
		AdditionalProperties: map[string]string{
			ServiceName: "bookstore",
		},
	}

	return op
}

package build

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-osm/internal/config"

	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/manifests"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var DefaultGenerationMethod string
var LatestVersion string
var WorkloadPath string
var CRDnames []string
var OverrideURL string
var AllVersions []string

//NewConfig creates the configuration for creating components
func NewConfig(version string) manifests.Config {
	return manifests.Config{
		Name:        smp.ServiceMesh_Type_name[int32(smp.ServiceMesh_OPEN_SERVICE_MESH)],
		MeshVersion: version,
		Filter: manifests.CrdFilter{
			RootFilter:    []string{"$[?(@.kind==\"CustomResourceDefinition\")]"},
			NameFilter:    []string{"$..[\"spec\"][\"names\"][\"kind\"]"},
			VersionFilter: []string{"$[0]..spec.versions[0]"},
			GroupFilter:   []string{"$[0]..spec"},
			SpecFilter:    []string{"$[0]..openAPIV3Schema.properties.spec"},
			ItrFilter:     []string{"$[?(@.spec.names.kind"},
			ItrSpecFilter: []string{"$[?(@.spec.names.kind"},
			VField:        "name",
			GField:        "group",
		},
	}
}
func GetDefaultURL(crd string) string {
	if OverrideURL != "" {
		return OverrideURL
	}
	return strings.Join([]string{"https://raw.githubusercontent.com/openservicemesh/osm/main/cmd/osm-bootstrap/crds", crd}, "/")
}
func init() {
	wd, _ := os.Getwd()
	WorkloadPath = filepath.Join(wd, "templates", "oam", "workloads")
	AllVersions, _ = utils.GetLatestReleaseTagsSorted("openservicemesh", "osm")
	if len(AllVersions) == 0 {
		return
	}
	CRDnames, _ = config.GetFileNames("openservicemesh", "osm", "/cmd/osm-bootstrap/crds/**")
	LatestVersion = AllVersions[len(AllVersions)-1]
	DefaultGenerationMethod = adapter.Manifests
}

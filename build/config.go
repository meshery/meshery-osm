package build

import (
	"encoding/json"
	"fmt"
	"io"
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
var MeshModelPath string
var CRDnames []string
var OverrideURL string
var AllVersions []string

const Component = "OSM"

var meshmodelmetadata = make(map[string]interface{})

var MeshModelConfig = adapter.MeshModelConfig{ //Move to build/config.go
	Category: "Orchestration & Management",
	Metadata: meshmodelmetadata,
}

// NewConfig creates the configuration for creating components
func NewConfig(version string) manifests.Config {
	return manifests.Config{
		Name:        smp.ServiceMesh_Type_name[int32(smp.ServiceMesh_OPEN_SERVICE_MESH)],
		Type:        Component,
		MeshVersion: version,
		CrdFilter: manifests.NewCueCrdFilter(manifests.ExtractorPaths{
			NamePath:    "spec.names.kind",
			IdPath:      "spec.names.kind",
			VersionPath: "spec.versions[0].name",
			GroupPath:   "spec.group",
			SpecPath:    "spec.versions[0].schema.openAPIV3Schema.properties.spec"}, false),
		ExtractCrds: func(manifest string) []string {
			crds := strings.Split(manifest, "---")
			return crds
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
	f, _ := os.Open("./build/meshmodel_metadata.json")
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	byt, _ := io.ReadAll(f)

	_ = json.Unmarshal(byt, &meshmodelmetadata)
	WorkloadPath = filepath.Join(wd, "templates", "oam", "workloads")
	MeshModelPath = filepath.Join(wd, "templates", "meshmodel", "components")
	AllVersions, _ = utils.GetLatestReleaseTagsSorted("openservicemesh", "osm")
	if len(AllVersions) == 0 {
		return
	}
	CRDnames, _ = config.GetFileNames("openservicemesh", "osm", "/cmd/osm-bootstrap/crds/**")
	LatestVersion = AllVersions[len(AllVersions)-1]
	DefaultGenerationMethod = adapter.Manifests
}

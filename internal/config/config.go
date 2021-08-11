package config

import (
	"path"

	"github.com/layer5io/meshery-adapter-library/common"
	"github.com/layer5io/meshery-adapter-library/config"
	configprovider "github.com/layer5io/meshery-adapter-library/config/provider"
	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshkit/utils"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var (
	configRootPath = path.Join(utils.GetHome(), ".meshery")

	ServerDefaults = map[string]string{
		"name":     smp.ServiceMesh_OPEN_SERVICE_MESH.Enum().String(),
		"type":     "adapter",
		"port":     "10009",
		"traceurl": "none",
	}

	MeshSpecDefaults = map[string]string{
		"name":    smp.ServiceMesh_OPEN_SERVICE_MESH.Enum().String(),
		"status":  status.NotInstalled,
		"version": "none",
	}

	ProviderConfigDefaults = map[string]string{
		configprovider.FilePath: configRootPath,
		configprovider.FileType: "yaml",
		configprovider.FileName: "osm",
	}

	KubeConfigDefaults = map[string]string{
		configprovider.FilePath: configRootPath,
		configprovider.FileType: "yaml",
		configprovider.FileName: "kubeconfig",
	}

	OperationsDefaults = getOperations(common.Operations)
)

func New(provider string) (config.Handler, error) {
	opts := configprovider.Options{
		ServerConfig:   ServerDefaults,
		MeshSpec:       MeshSpecDefaults,
		ProviderConfig: ProviderConfigDefaults,
		Operations:     OperationsDefaults,
	}
	switch provider {
	case configprovider.ViperKey:
		return configprovider.NewViper(opts)
	case configprovider.InMemKey:
		return configprovider.NewInMem(opts)
	}
	return nil, config.ErrEmptyConfig
}

func NewKubeconfigBuilder(provider string) (config.Handler, error) {

	opts := configprovider.Options{}

	// Config environment
	opts.ProviderConfig = KubeConfigDefaults

	// Config provider
	switch provider {
	case configprovider.ViperKey:
		return configprovider.NewViper(opts)
	case configprovider.InMemKey:
		return configprovider.NewInMem(opts)
	}
	return nil, config.ErrEmptyConfig
}

// RootPath returns the root config path
func RootPath() string {
	return configRootPath
}

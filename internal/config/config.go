package config

import (
	"path"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/common"
	"github.com/layer5io/meshery-adapter-library/config"
	"github.com/layer5io/meshery-adapter-library/status"
	configprovider "github.com/layer5io/meshkit/config/provider"
	"github.com/layer5io/meshkit/utils"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

const (
	// OAM Metadata constants
	OAMAdapterNameMetadataKey       = "adapter.meshery.io/name"
	OAMComponentCategoryMetadataKey = "ui.meshery.io/category"
)

var (
	configRootPath = path.Join(utils.GetHome(), ".meshery")

	ServerDefaults = map[string]string{
		"name":     smp.ServiceMesh_OPEN_SERVICE_MESH.Enum().String(),
		"type":     "adapter",
		"port":     "10009",
		"traceurl": status.None,
	}

	MeshSpecDefaults = map[string]string{
		"name":    smp.ServiceMesh_OPEN_SERVICE_MESH.Enum().String(),
		"status":  status.NotInstalled,
		"version": status.None,
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

func New(provider string) (h config.Handler, err error) {
	opts := configprovider.Options{
		FilePath: configRootPath,
		FileName: "osm",
		FileType: "yaml",
	}
	switch provider {
	case configprovider.ViperKey:
		h, err = configprovider.NewViper(opts)
		if err != nil {
			return nil, err
		}
	case configprovider.InMemKey:
		h, err = configprovider.NewInMem(opts)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrEmptyConfig
	}
	// Setup server config
	if err := h.SetObject(adapter.ServerKey, ServerDefaults); err != nil {
		return nil, err
	}

	// Setup mesh config
	if err := h.SetObject(adapter.MeshSpecKey, MeshSpecDefaults); err != nil {
		return nil, err
	}

	// Setup Operations Config
	if err := h.SetObject(adapter.OperationsKey, OperationsDefaults); err != nil {
		return nil, err
	}

	return h, nil
}

func NewKubeconfigBuilder(provider string) (config.Handler, error) {
	opts := configprovider.Options{
		FilePath: configRootPath,
		FileType: "yaml",
		FileName: "kubeconfig",
	}

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

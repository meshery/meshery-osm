// Copyright 2020 Layer5, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/api/grpc"
	internalconfig "github.com/layer5io/meshery-osm/internal/config"
	"github.com/layer5io/meshery-osm/osm"
	"github.com/layer5io/meshery-osm/osm/oam"
	configprovider "github.com/layer5io/meshkit/config/provider"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshkit/utils/manifests"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var (
	serviceName = "osm-adapter"
	version     = "none"
	gitsha      = "none"
)

func main() {
	log, err := logger.New(serviceName, logger.Options{Format: logger.SyslogLogFormat, DebugLevel: displayDebugLogs()})
	if err != nil {
		fmt.Println("Logger Init Failed", err.Error())
		os.Exit(1)
	}

	if err = os.Setenv("KUBECONFIG", path.Join(
		internalconfig.KubeConfigDefaults[configprovider.FilePath],
		fmt.Sprintf("%s.%s", internalconfig.KubeConfigDefaults[configprovider.FileName],
			internalconfig.KubeConfigDefaults[configprovider.FileType],
		)),
	); err != nil {
		// Fail silently
		log.Warn(err)
	}

	cfg, err := internalconfig.New(configprovider.ViperKey)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	kubeconfigHandler, err := internalconfig.NewKubeconfigBuilder(configprovider.ViperKey)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	service := &grpc.Service{}
	_ = cfg.GetObject(adapter.ServerKey, &service)

	service.Handler = osm.New(cfg, log, kubeconfigHandler)
	service.Channel = make(chan interface{}, 100)
	service.StartedAt = time.Now()
	service.Version = version
	service.GitSHA = gitsha

	go registerCapabilities(service.Port, log)        //Registering static capabilities
	go registerDynamicCapabilities(service.Port, log) //Registering latest capabilities periodically
	// Server Initialization
	log.Info("Adaptor Listening at port: ", service.Port)
	err = grpc.Start(service, nil)
	if err != nil {
		log.Error(grpc.ErrGrpcServer(err))
		os.Exit(1)
	}
}

// This init function can help adapters create the configuration logic work well, so do not remove it although that's
// not a good idea.
func init() {
	err := os.MkdirAll(path.Join(utils.GetHome(), ".meshery", "bin"), 0750)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

// displayDebugLogs will return true if the "DEBUG" env var
// is set to "true"
func displayDebugLogs() bool {
	return os.Getenv("DEBUG") == "true"
}

func mesheryServerAddress() string {
	meshReg := os.Getenv("MESHERY_SERVER")

	if meshReg != "" {
		if strings.HasPrefix(meshReg, "http") {
			return meshReg
		}

		return "http://" + meshReg
	}

	return "http://localhost:9081"
}

func serviceAddress() string {
	svcAddr := os.Getenv("SERVICE_ADDR")

	if svcAddr != "" {
		return svcAddr
	}

	return "mesherylocal.layer5.io"
}

func registerCapabilities(port string, log logger.Handler) {
	// Register workloads
	if err := oam.RegisterWorkloads(mesheryServerAddress(), serviceAddress()+":"+port); err != nil {
		log.Info(err.Error())
	}

	// Register traits
	if err := oam.RegisterTraits(mesheryServerAddress(), serviceAddress()+":"+port); err != nil {
		log.Info(err.Error())
	}
}

func registerDynamicCapabilities(port string, log logger.Handler) {
	registerWorkloads(port, log)
	//Start the ticker
	const reRegisterAfter = 24
	ticker := time.NewTicker(reRegisterAfter * time.Hour)
	for {
		<-ticker.C
		registerWorkloads(port, log)
	}

}

func registerWorkloads(port string, log logger.Handler) {
	appVersion, chartVersion, err := getLatestValidAppVersionAndChartVersion()
	if err != nil {
		log.Info("Could not get latest version")
		return
	}
	log.Info("Registering latest workload components for version ", appVersion)
	// Register workloads
	if err := adapter.RegisterWorkLoadsDynamically(mesheryServerAddress(), serviceAddress()+":"+port, &adapter.DynamicComponentsConfig{
		TimeoutInMinutes: 60,
		URL:              "https://openservicemesh.github.io/osm/osm-" + chartVersion + ".tgz",
		GenerationMethod: adapter.HelmCHARTS,
		Config: manifests.Config{
			Name:        smp.ServiceMesh_Type_name[int32(smp.ServiceMesh_OPEN_SERVICE_MESH)],
			MeshVersion: appVersion,
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
		},
		Operation: internalconfig.OSMOperation,
	}); err != nil {
		log.Info(err.Error())
		return
	}
	log.Info("Latest workload components successfully registered.")
}

// returns latest valid appversion and chartversion
func getLatestValidAppVersionAndChartVersion() (string, string, error) {
	release, err := internalconfig.GetLatestReleases(100)
	if err != nil {
		return "", "", osm.ErrGetLatestRelease(err)
	}
	//loops through latest  app versions untill it finds one which is available in helm chart's index.yaml
	for _, rel := range release {
		if chartVersion, err := kubernetes.HelmAppVersionToChartVersion("https://openservicemesh.github.io/osm", "osm", rel.TagName); err == nil {
			return rel.TagName, chartVersion, nil
		}

	}
	return "", "", osm.ErrGetLatestRelease(err)
}

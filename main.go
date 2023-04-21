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

	"github.com/google/uuid"
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/api/grpc"
	"github.com/layer5io/meshery-osm/build"
	internalconfig "github.com/layer5io/meshery-osm/internal/config"
	"github.com/layer5io/meshery-osm/osm"
	"github.com/layer5io/meshery-osm/osm/oam"
	configprovider "github.com/layer5io/meshkit/config/provider"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/events"
)

var (
	serviceName = "osm-adapter"
	version     = "edge"
	gitsha      = "none"
	instanceID  = uuid.NewString()
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
	e := events.NewEventStreamer()
	service.Handler = osm.New(cfg, log, kubeconfigHandler, e)
	service.EventStreamer = e
	service.StartedAt = time.Now()
	service.Version = version
	service.GitSHA = gitsha

	go registerCapabilities(service.Port, log)        //Registering static capabilities
	go registerDynamicCapabilities(service.Port, log) //Registering latest capabilities periodically
	// Server Initialization
	log.Info("Adapter listening on port: ", service.Port)
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

	return "localhost"
}

func registerCapabilities(port string, log logger.Handler) {
	// Register meshmodel components
	if err := oam.RegisterMeshModelComponents(instanceID, mesheryServerAddress(), serviceAddress(), port); err != nil {
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
	//First we create and store any new components if available
	version := build.LatestVersion
	gm := build.DefaultGenerationMethod
	// Prechecking to skip comp gen
	if os.Getenv("FORCE_DYNAMIC_REG") != "true" && oam.AvailableVersions[version] {
		log.Info("Components available statically for version ", version, ". Skipping dynamic component registeration")
		return
	}
	log.Info("Registering latest workload components for version ", version)
	// Register workloads
	for _, crd := range build.CRDnames {
		log.Info("Generating components for ", crd)
		if err := adapter.CreateComponents(adapter.StaticCompConfig{
			URL:             build.GetDefaultURL(crd),
			Method:          gm,
			MeshModelPath:   build.MeshModelPath,
			MeshModelConfig: build.MeshModelConfig,
			DirName:         version,
			Config:          build.NewConfig(version),
		}); err != nil {
			log.Error(err)
			return
		}
		log.Info(crd, " created")
	}

	//*The below log is checked in the workflows. If you change this log, reflect that change in the workflow where components are generated
	log.Info("Component creation completed for version ", version)
	//Now we will register in case
	log.Info("Registering workloads with Meshery Server for version ", version)
	if err := oam.RegisterMeshModelComponents(instanceID, mesheryServerAddress(), serviceAddress(), port); err != nil {
		log.Info(err.Error())
		return
	}

	log.Info("Latest workload components successfully registered.")
}

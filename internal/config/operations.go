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

package config

import (
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/meshes"
)

var (
	OSMOperation          = "osm"
	OSMBookStoreOperation = "osm_bookstore_app"
	ServiceName           = "service_name"
)

func getOperations(op adapter.Operations) adapter.Operations {
	op[OSMOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_INSTALL),
		Description: "OSM Service Mesh",
		Versions:    []adapter.Version{"v0.6.0", "v0.5.0"},
		Templates:   adapter.NoneTemplate,
	}

	op[OSMBookStoreOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_SAMPLE_APPLICATION),
		Description: "Bookstore Application",
		Versions:    []adapter.Version{"v0.6.0", "v0.5.0"},
		Templates: []adapter.Template{
			"https://raw.githubusercontent.com/openservicemesh/osm/main/docs/example/manifests/apps/bookbuyer.yaml",
			"https://raw.githubusercontent.com/openservicemesh/osm/main/docs/example/manifests/apps/bookstore-v1.yaml",
			"https://raw.githubusercontent.com/openservicemesh/osm/main/docs/example/manifests/apps/bookthief.yaml",
			"https://raw.githubusercontent.com/openservicemesh/osm/main/docs/example/manifests/apps/bookwarehouse.yaml",
			"https://raw.githubusercontent.com/openservicemesh/osm/main/docs/example/manifests/apps/traffic-split.yaml",
			"file://templates/osm-bookstore-traffic-access-v1.yaml",
		},
		AdditionalProperties: map[string]string{
			ServiceName: "bookstore",
		},
	}

	return op
}

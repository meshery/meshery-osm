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
	OSMOperation       = "osm"
	OSMSampleBookBuyer = "osm-sample-book-buyer"
)

func getOperations(op adapter.Operations) adapter.Operations {
	op[OSMOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_INSTALL),
		Description: "OSM Service Mesh",
		Versions:    []adapter.Version{"v0.3.0", "v0.2.0"},
		Templates:   adapter.NoneTemplate,
	}

	// add sample book buyer operation
	op[OSMSampleBookBuyer] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_SAMPLE_APPLICATION),
		Description: "OSM Sample Application",
		Versions:    []adapter.Version{"v0.3.0", "v0.2.0"},
		Templates:   adapter.NoneTemplate,
	}

	return op
}

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

package osm

import (
	"fmt"
	"sync"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

func (h *Handler) installOSM(del bool, version, ns string, kubeconfigs []string) (string, error) {
	h.Log.Debug(fmt.Sprintf("Requested install of version: %s", version))
	h.Log.Debug(fmt.Sprintf("Requested action is delete: %v", del))
	h.Log.Debug(fmt.Sprintf("Requested action is in namespace: %s", ns))

	st := status.Installing
	if del {
		st = status.Removing
	}

	err := h.Config.GetObject(adapter.MeshSpecKey, h)
	if err != nil {
		return st, ErrMeshConfig(err)
	}

	h.Log.Info("Installing...")
	err = h.applyHelmChart(del, version, ns, kubeconfigs)
	if err != nil {
		return st, ErrApplyHelmChart(err)
	}

	st = status.Installed
	if del {
		st = status.Removed
	}

	return st, nil
}

func (h *Handler) applyHelmChart(del bool, version, namespace string, kubeconfigs []string) error {

	repo := "https://openservicemesh.github.io/osm/"
	chart := "osm"
	var act mesherykube.HelmChartAction
	if del {
		act = mesherykube.UNINSTALL
	} else {
		act = mesherykube.INSTALL
	}
	var wg sync.WaitGroup
	var errMx sync.Mutex
	var errs []error
	for _, kubeconfig := range kubeconfigs {
		wg.Add(1)
		go func(kubeconfig string) {
			defer wg.Done()
			kClient, err := mesherykube.New([]byte(kubeconfig))
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}
			err = kClient.ApplyHelmChart(mesherykube.ApplyHelmChartConfig{
				ChartLocation: mesherykube.HelmChartLocation{
					Repository: repo,
					Chart:      chart,
					Version:    version,
				},
				Namespace:       namespace,
				Action:          act,
				CreateNamespace: true,
			})
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}
		}(kubeconfig)
	}
	wg.Wait()
	if len(errs) != 0 {
		return mergeErrors(errs)
	}
	return nil
}

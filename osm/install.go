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
	"net/http"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"gopkg.in/yaml.v2"
)

// HelmIndex holds the index.yaml data in the struct format
type HelmIndex struct {
	APIVersion string      `yaml:"apiVersion"`
	Entries    HelmEntries `yaml:"entries"`
}

// HelmEntries holds the data for all of the entries present
// in the helm repository
type HelmEntries map[string][]HelmEntryMetadata

// HelmEntryMetadata is the struct for holding the metadata
// associated with a helm repositories' entry
type HelmEntryMetadata struct {
	APIVersion string `yaml:"apiVersion"`
	AppVersion string `yaml:"appVersion"`
	Name       string `yaml:"name"`
	Version    string `yaml:"version"`
}

func (h *Handler) installOSM(del bool, version, ns string) (string, error) {
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
	err = h.applyHelmChart(del, version, ns)
	if err != nil {
		return st, ErrApplyHelmChart(err)
	}

	st = status.Installed
	if del {
		st = status.Removed
	}

	return st, nil
}

func (h *Handler) applyHelmChart(del bool, version, namespace string) error {
	kClient := h.MesheryKubeclient

	repo := "https://openservicemesh.github.io/osm/"
	chart := "osm"

	chartVersion, err := ConvertAppVersionToChartVersion(repo, chart, version)
	if err != nil {
		return ErrConvertingAppVersionToChartVersion(err)
	}

	err = kClient.ApplyHelmChart(mesherykube.ApplyHelmChartConfig{
		ChartLocation: mesherykube.HelmChartLocation{
			Repository: repo,
			Chart:      chart,
			Version:    chartVersion,
		},
		Namespace:       namespace,
		Delete:          del,
		CreateNamespace: true,
	})

	return err
}

// ConvertAppVersionToChartVersion takes in the repo, chart and app version and
// returns the corresponding chart version for the same
func ConvertAppVersionToChartVersion(repo, chart, appVersion string) (string, error) {
	appVersion = normalizeVersion(appVersion)

	helmIndex, err := CreateHelmIndex(repo)
	if err != nil {
		return "", ErrCreatingHelmIndex(err)
	}

	entryMetadata, exists := helmIndex.Entries.GetEntryWithAppVersion(chart, appVersion)
	if !exists {
		return "", ErrEntryWithAppVersionNotExists(chart, appVersion)
	}

	return entryMetadata.Version, nil
}

// CreateHelmIndex takes in the repo name and creates a
// helm index for it. Helm index is basically marshalled version of
// index.yaml file present in the remote helm repository
func CreateHelmIndex(repo string) (*HelmIndex, error) {
	url := fmt.Sprintf("%s/index.yaml", repo)

	// helm repository path will alaways be varaible hence,
	// #nosec
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, ErrHelmRepositoryNotFound(repo, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var hi HelmIndex
	dec := yaml.NewDecoder(resp.Body)
	if err := dec.Decode(&hi); err != nil {
		return nil, ErrDecodeYaml(err)
	}

	return &hi, nil
}

// GetEntryWithAppVersion takes in the entry name and the appversion and returns the corresponding
// metadata for the parameters if it exists
func (helmEntries HelmEntries) GetEntryWithAppVersion(entry, appVersion string) (HelmEntryMetadata, bool) {
	hem, ok := helmEntries[entry]
	if !ok {
		return HelmEntryMetadata{}, false
	}

	for _, v := range hem {
		if v.Name == entry && v.AppVersion == appVersion {
			return v, true
		}
	}

	return HelmEntryMetadata{}, false
}

// normalizeVerion takes in a version and adds "v" prefix
// if it isn't already present
func normalizeVersion(version string) string {
	if strings.HasPrefix(version, "v") {
		return version
	}

	return fmt.Sprintf("v%s", version)
}

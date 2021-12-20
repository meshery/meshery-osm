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
	"sync"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshkit/utils/walker"
)

// Release is used to save the release informations
type Release struct {
	ID      int             `json:"id,omitempty"`
	TagName string          `json:"tag_name,omitempty"`
	Name    adapter.Version `json:"name,omitempty"`
	Draft   bool            `json:"draft,omitempty"`
	Assets  []*Asset        `json:"assets,omitempty"`
}

// Asset describes the github release asset object
type Asset struct {
	Name        string `json:"name,omitempty"`
	State       string `json:"state,omitempty"`
	DownloadURL string `json:"browser_download_url,omitempty"`
}

func appendThreadSafe(arr *[]string, s string, m *sync.RWMutex) {
	m.Lock()
	defer m.Unlock()
	*arr = append((*arr), s)
}

// GetFileNames takes the url of a github repo and the path to a directory. Then returns all the filenames from that directory
func GetFileNames(owner string, repo string, path string) ([]string, error) {
	g := walker.NewGit()
	var fs []string
	var m sync.RWMutex
	err := g.Owner(owner).Repo(repo).Branch("master").Root(path).RegisterFileInterceptor(func(f walker.File) error {
		appendThreadSafe(&fs, f.Name, &m)
		return nil
	}).Walk()
	return fs, err
}

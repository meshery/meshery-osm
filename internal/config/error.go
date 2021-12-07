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
	"github.com/layer5io/meshkit/errors"
)

const (
	// ErrEmptyConfigCode represents the error when the configuration is either empty
	// or is invalid
	ErrEmptyConfigCode = "1021"

	// ErrGetLatestReleasesCode represents the error which occurs during the process of getting
	// latest releases
	ErrGetLatestReleasesCode = "1022"

	// ErrGetLatestReleaseNamesCode represents the error which occurs during the process of extracting
	// release names
	ErrGetLatestReleaseNamesCode = "1023"

	ErrGetManifestNamesCode = "1024"
)

var (
	// ErrEmptyConfig error is the error when config is invalid
	ErrEmptyConfig = errors.New(ErrEmptyConfigCode, errors.Alert, []string{"Config is empty"}, []string{}, []string{}, []string{})
)

// ErrGetLatestReleases is the error for fetching nsm-mesh releases
func ErrGetLatestReleases(err error) error {
	return errors.New(ErrGetLatestReleasesCode, errors.Alert, []string{"Unable to fetch release info"}, []string{err.Error()}, []string{}, []string{})
}

// ErrGetLatestReleaseNames is the error for fetching nsm-mesh releases
func ErrGetLatestReleaseNames(err error) error {
	return errors.New(ErrGetLatestReleaseNamesCode, errors.Alert, []string{"Failed to extract release names"}, []string{err.Error()}, []string{}, []string{})
}

// ErrGetManifestNames is the error for fetching consul manifest names
func ErrGetManifestNames(err error) error {
	return errors.New(ErrGetManifestNamesCode, errors.Alert, []string{"Unable to fetch manifest names from github"}, []string{err.Error()}, []string{}, []string{})
}

// Package osm - Error codes for the adapter
package osm

import (
	"github.com/layer5io/meshkit/errors"
)

var (
	// Error code for failed service mesh installation

	// ErrInstallOSMCode represents the errors which are generated
	// during open service mesh install process
	ErrInstallOSMCode = "osm_test_code"

	// ErrTarXZFCode represents the errors which are generated
	// during decompressing and extracting tar.gz file
	ErrTarXZFCode = "osm_test_code"

	// ErrMeshConfigCode represents the errors which are generated
	// when an invalid mesh config is found
	ErrMeshConfigCode = "osm_test_code"

	// ErrRunOsmCtlCmdCode represents the errors which are generated
	// during fetch manifest process
	ErrRunOsmCtlCmdCode = "osm_test_code"

	// ErrDownloadBinaryCode represents the errors which are generated
	// during binary download process
	ErrDownloadBinaryCode = "osm_test_code"

	// ErrInstallBinaryCode represents the errors which are generated
	// during binary installation process
	ErrInstallBinaryCode = "osm_test_code"

	// ErrSampleAppCode represents the errors which are generated
	// duing sample app installation
	ErrSampleAppCode = "osm_test_code"

	// ErrCustomOperationCode represents the errors which are generated
	// when an invalid addon operation is requested
	ErrCustomOperationCode = "osm_test_code"

	// ErrCreatingNSCode represents the errors which are generated
	// during the process of creating a namespace
	ErrCreatingNSCode = "osm_test_code"

	// ErrRunExecutableCode represents the errors which are generated
	// during the running a executable
	ErrRunExecutableCode = "osm_test_code"

	// ErrSidecarInjectionCode represents the errors which are generated
	// during the process of enabling/disabling sidecar injection
	ErrSidecarInjectionCode = "osm_test_code"

	// ErrOpInvalidCode represents the error which is generated when
	// there is an invalid operation
	ErrOpInvalidCode = "osm_test_code"

	// ErrOpInvalid represents the errors which are generated
	// when an invalid operation is requested
	ErrOpInvalid = errors.New(ErrOpInvalidCode, errors.Alert, []string{"Invalid operation"}, []string{}, []string{}, []string{})
)

// ErrInstallOSM is the error for install mesh
func ErrInstallOSM(err error) error {
	return errors.New(ErrInstallOSMCode, errors.Alert, []string{"Error with osm operation: ", err.Error()}, []string{}, []string{}, []string{})
}

// ErrTarXZF is the error for unzipping the file
func ErrTarXZF(err error) error {
	return errors.New(ErrTarXZFCode, errors.Alert, []string{"Error while extracting file: ", err.Error()}, []string{}, []string{}, []string{})
}

// ErrMeshConfig is the error for mesh config
func ErrMeshConfig(err error) error {
	return errors.New(ErrMeshConfigCode, errors.Alert, []string{"Error configuration mesh: ", err.Error()}, []string{}, []string{}, []string{})
}

// ErrRunOsmCtlCmd is the error for mesh port forward
func ErrRunOsmCtlCmd(err error, des string) error {
	return errors.New(ErrRunOsmCtlCmdCode, errors.Alert, []string{"Error running osmctl command: ", des}, []string{err.Error()}, []string{}, []string{})
}

// ErrDownloadBinary is the error while downloading osm binary
func ErrDownloadBinary(err error) error {
	return errors.New(ErrDownloadBinaryCode, errors.Alert, []string{"Error downloading osmctl binary: ", err.Error()}, []string{}, []string{}, []string{})
}

// ErrInstallBinary is the error while downloading osm binary
func ErrInstallBinary(err error) error {
	return errors.New(ErrInstallBinaryCode, errors.Alert, []string{"Error installing osmctl binary: ", err.Error()}, []string{}, []string{}, []string{})
}

// ErrSampleApp is the error for streaming event
func ErrSampleApp(err error) error {
	return errors.New(ErrSampleAppCode, errors.Alert, []string{"Error with sample app operation: ", err.Error()}, []string{}, []string{}, []string{})
}

// ErrCustomOperation is the error for streaming event
func ErrCustomOperation(err error) error {
	return errors.New(ErrCustomOperationCode, errors.Alert, []string{"Error with custom operation: ", err.Error()}, []string{}, []string{}, []string{})
}

// ErrCreatingNS is the error while creating the namespace
func ErrCreatingNS(err error) error {
	return errors.New(ErrCreatingNSCode, errors.Alert, []string{"Error creating namespace: ", err.Error()}, []string{}, []string{}, []string{})
}

// ErrRunExecutable is the error while running an executable
func ErrRunExecutable(err error) error {
	return errors.New(ErrRunExecutableCode, errors.Alert, []string{"Error running executable: ", err.Error()}, []string{}, []string{}, []string{})
}

// ErrSidecarInjection is the error while enabling/disabling sidecar injection
// on a particular namespace
func ErrSidecarInjection(err error) error {
	return errors.New(ErrSidecarInjectionCode, errors.Alert, []string{"Error sidecar injection: ", err.Error()}, []string{}, []string{}, []string{})
}

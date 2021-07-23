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
	ErrOpInvalid = errors.New(ErrOpInvalidCode, errors.Alert, []string{"Invalid operation"}, []string{"Istio adapter recived an invalid operation from the meshey server"}, []string{"The operation is not supported by the adapter", "Invalid operation name"}, []string{"Check if the operation name is valid and supported by the adapter"})
)

// ErrInstallOSM is the error for install mesh
func ErrInstallOSM(err error) error {
	return errors.New(ErrInstallOSMCode, errors.Alert, []string{"Error with osm operation"}, []string{"Error occured while installing osm mesh through osmctl", err.Error()}, []string{}, []string{})
}

// ErrTarXZF is the error for unzipping the file
func ErrTarXZF(err error) error {
	return errors.New(ErrTarXZFCode, errors.Alert, []string{"Error while extracting file"}, []string{err.Error()}, []string{"The gzip might be corrupt"}, []string{"Retry the operation"})
}

// ErrMeshConfig is the error for mesh config
func ErrMeshConfig(err error) error {
	return errors.New(ErrMeshConfigCode, errors.Alert, []string{"Error configuration mesh"}, []string{err.Error(), "Error getting MeshSpecKey config from in-memory configuration"}, []string{}, []string{"Reconnect the adaptor to the meshkit server"})
}

// ErrRunOsmCtlCmd is the error for mesh port forward
func ErrRunOsmCtlCmd(err error, des string) error {
	return errors.New(ErrRunOsmCtlCmdCode, errors.Alert, []string{"Error running istioctl command"}, []string{err.Error()}, []string{"Corrupted istioctl binary", "Command might be invalid"}, []string{})
}

// ErrDownloadBinary is the error while downloading osm binary
func ErrDownloadBinary(err error) error {
	return errors.New(ErrDownloadBinaryCode, errors.Alert, []string{"Error downloading osm binary"}, []string{err.Error(), "Error occured while download osm binary from its github release"}, []string{"Checkout https://docs.github.com/en/rest/reference/repos#releases for more details"}, []string{})
}

// ErrInstallBinary is the error while downloading osm binary
func ErrInstallBinary(err error) error {
	return errors.New(ErrInstallBinaryCode, errors.Alert, []string{"Error installing osm binary"}, []string{err.Error()}, []string{"Corrupted osm release binary", "Invalid installation location"}, []string{})
}

// ErrSampleApp is the error for streaming event
func ErrSampleApp(err error) error {
	return errors.New(ErrSampleAppCode, errors.Alert, []string{"Error with sample app operation"}, []string{err.Error(), "Error occured while trying to install a sample application using manifests"}, []string{"Invalid kubeclient config", "Invalid manifest"}, []string{"Reconnect your adapter to meshery server to refresh the kubeclient"})
}

// ErrCustomOperation is the error for streaming event
func ErrCustomOperation(err error) error {
	return errors.New(ErrCustomOperationCode, errors.Alert, []string{"Error with custom operation"}, []string{"Error occured while applying custom manifest to the cluster", err.Error()}, []string{"Invalid kubeclient config", "Invalid manifest"}, []string{"Reupload the kubconfig in the Meshery Server and reconnect the adapter"})
}

// ErrCreatingNS is the error while creating the namespace
func ErrCreatingNS(err error) error {
	return errors.New(ErrCreatingNSCode, errors.Alert, []string{"Error creating namespace"}, []string{"Error occured while applying manifest to create a namespace", err.Error()}, []string{"Invalid kubeclient config", "Invalid manifest"}, []string{"Reupload the kubeconfig in the Meshery Server and reconnect the adapter"})
}

// ErrRunExecutable is the error while running an executable
func ErrRunExecutable(err error) error {
	return errors.New(ErrRunExecutableCode, errors.Alert, []string{"Error running executable"}, []string{err.Error()}, []string{"Corrupted binary", "Invalid operation"}, []string{"Check if the adaptor is executing a deprecated command"})
}

// ErrSidecarInjection is the error while enabling/disabling sidecar injection
// on a particular namespace
func ErrSidecarInjection(err error) error {
	return errors.New(ErrSidecarInjectionCode, errors.Alert, []string{"Error occured while injection sidecar"}, []string{"Error occured while injecting sidercar using osm(ctl) `osm namespace add/remove <name>` ", err.Error()}, []string{"Corrupted binary", "Invalidoperation"}, []string{"Check if the adaptor is executing a deprecated command"})
}

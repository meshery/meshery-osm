// Package osm - Error codes for the adapter
package osm

import (
	"github.com/layer5io/meshkit/errors"
)

var (
	// Error code for failed service mesh installation

	// ErrInstallOSMCode represents the errors which are generated
	// during open service mesh install process
	ErrInstallOSMCode = "1000"

	// ErrTarXZFCode represents the errors which are generated
	// during decompressing and extracting tar.gz file
	ErrTarXZFCode = "1001"

	// ErrMeshConfigCode represents the errors which are generated
	// when an invalid mesh config is found
	ErrMeshConfigCode = "1002"

	// ErrRunOsmCtlCmdCode represents the errors which are generated
	// during fetch manifest process
	ErrRunOsmCtlCmdCode = "1003"

	// ErrDownloadBinaryCode represents the errors which are generated
	// during binary download process
	ErrDownloadBinaryCode = "1004"

	// ErrInstallBinaryCode represents the errors which are generated
	// during binary installation process
	ErrInstallBinaryCode = "1005"

	// ErrSampleAppCode represents the errors which are generated
	// duing sample app installation
	ErrSampleAppCode = "1006"

	// ErrCustomOperationCode represents the errors which are generated
	// when an invalid addon operation is requested
	ErrCustomOperationCode = "1007"

	// ErrCreatingNSCode represents the errors which are generated
	// during the process of creating a namespace
	ErrCreatingNSCode = "1008"

	// ErrRunExecutableCode represents the errors which are generated
	// during the running a executable
	ErrRunExecutableCode = "1009"

	// ErrSidecarInjectionCode represents the errors which are generated
	// during the process of enabling/disabling sidecar injection
	ErrSidecarInjectionCode = "1010"

	// ErrOpInvalidCode represents the error which is generated when
	// there is an invalid operation
	ErrOpInvalidCode = "1011"

	// ErrApplyHelmChartCode represents the error which are generated
	// during the process of applying helm chart
	ErrApplyHelmChartCode = "1012"

	// ErrNilClientCode represents the error code which is
	// generated when kubernetes client is nil
	ErrNilClientCode = "1013"

	// ErrInvalidOAMComponentTypeCode represents the error code which is
	// generated when an invalid oam component is requested
	ErrInvalidOAMComponentTypeCode = "1014"

	// ErrOSMCoreComponentFailCode represents the error code which is
	// generated when an osm core operations fails
	ErrOSMCoreComponentFailCode = "1015"
	// ErrProcessOAMCode represents the error code which is
	// generated when an OAM operations fails
	ErrProcessOAMCode = "1016"
	// ErrParseOSMCoreComponentCode represents the error code which is
	// generated when osm core component manifest parsing fails
	ErrParseOSMCoreComponentCode = "1017"
	// ErrParseOAMComponentCode represents the error code which is
	// generated during the OAM component parsing
	ErrParseOAMComponentCode = "1018"
	// ErrParseOAMConfigCode represents the error code which is
	// generated during the OAM configuration parsing
	ErrParseOAMConfigCode = "1019"

	// ErrOpInvalid represents the errors which are generated
	// when an invalid operation is requested
	ErrOpInvalid = errors.New(ErrOpInvalidCode, errors.Alert, []string{"Invalid operation"}, []string{"Istio adapter recived an invalid operation from the meshey server"}, []string{"The operation is not supported by the adapter", "Invalid operation name"}, []string{"Check if the operation name is valid and supported by the adapter"})

	// ErrNilClient represents the error which is
	// generated when kubernetes client is nil
	ErrNilClient = errors.New(ErrNilClientCode, errors.Alert, []string{"kubernetes client not initialized"}, []string{"Kubernetes client is nil"}, []string{"kubernetes client not initialized"}, []string{"Reconnect the adaptor to Meshery server"})

	// ErrParseOAMComponent represents the error which is
	// generated during the OAM component parsing
	ErrParseOAMComponent = errors.New(ErrParseOAMComponentCode, errors.Alert, []string{"error parsing the component"}, []string{"Error occured while prasing application component in the OAM request made"}, []string{"Invalid OAM component passed in OAM request"}, []string{"Check if your request has vaild OAM components"})

	// ErrParseOAMConfig represents the error which is
	// generated during the OAM configuration parsing
	ErrParseOAMConfig = errors.New(ErrParseOAMConfigCode, errors.Alert, []string{"error parsing the configuration"}, []string{"Error occured while prasing component config in the OAM request made"}, []string{"Invalid OAM config passed in OAM request"}, []string{"Check if your request has vaild OAM config"})

	// ErrGetLatestReleaseCode represents the error which is
	// generated when the latest stable version could not
	// be fetched during runtime component registeration
	ErrGetLatestReleaseCode = "1020"
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

// ErrApplyHelmChart is the error for applying helm chart
func ErrApplyHelmChart(err error) error {
	return errors.New(ErrApplyHelmChartCode, errors.Alert, []string{"Error occured while applying Helm Chart"}, []string{err.Error()}, []string{}, []string{})
}

// ErrParseOSMCoreComponent is the error when osm core component manifest parsing fails
func ErrParseOSMCoreComponent(err error) error {
	return errors.New(ErrParseOSMCoreComponentCode, errors.Alert, []string{"osm core component manifest parsing failing"}, []string{err.Error()}, []string{}, []string{})
}

// ErrInvalidOAMComponentType is the error when the OAM component name is not valid
func ErrInvalidOAMComponentType(compName string) error {
	return errors.New(ErrInvalidOAMComponentTypeCode, errors.Alert, []string{"invalid OAM component name: ", compName}, []string{}, []string{}, []string{})
}

// ErrOSMCoreComponentFail is the error when core osm component processing fails
func ErrOSMCoreComponentFail(err error) error {
	return errors.New(ErrOSMCoreComponentFailCode, errors.Alert, []string{"error in osm core component"}, []string{err.Error()}, []string{}, []string{})
}

// ErrProcessOAM is a generic error which is thrown when an OAM operations fails
func ErrProcessOAM(err error) error {
	return errors.New(ErrProcessOAMCode, errors.Alert, []string{"error performing OAM operations"}, []string{err.Error()}, []string{}, []string{})
}

// ErrGetLatestRelease is the error for get latest versions
func ErrGetLatestRelease(err error) error {
	return errors.New(ErrGetLatestReleaseCode, errors.Alert, []string{"Could not get latest version"}, []string{err.Error()}, []string{"Latest version could not be found at the specified url"}, []string{})
}

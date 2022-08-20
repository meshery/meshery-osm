package osm

import (
	"context"
	"fmt"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/common"
	"github.com/layer5io/meshery-adapter-library/status"
	internalconfig "github.com/layer5io/meshery-osm/internal/config"
	"github.com/layer5io/meshkit/errors"
)

// ApplyOperation function contains the operation handlers
func (h *Handler) ApplyOperation(ctx context.Context, request adapter.OperationRequest, hchan *chan interface{}) error {
	err := h.CreateKubeconfigs(request.K8sConfigs)
	if err != nil {
		return err
	}
	h.SetChannel(hchan)
	kubeconfigs := request.K8sConfigs
	operations := make(adapter.Operations)
	err = h.Config.GetObject(adapter.OperationsKey, &operations)
	if err != nil {
		return err
	}

	e := &adapter.Event{
		Operationid: request.OperationID,
		Summary:     status.Deploying,
		Details:     "Operation is not supported",
		Component:   internalconfig.ServerDefaults["type"],
		ComponentName: internalconfig.ServerDefaults["name"],
	}

	//deployment
	switch request.OperationName {
	case internalconfig.OSMOperation:
		go func(hh *Handler, ee *adapter.Event) {
			version := string(operations[request.OperationName].Versions[0])
			stat, err := hh.installOSM(request.IsDeleteOperation, version, request.Namespace, kubeconfigs)
			if err != nil {
				summary := fmt.Sprintf("Error while %s Open service mesh", stat)
				hh.streamErr(summary, e, err)
				return
			}
			ee.Summary = fmt.Sprintf("Open service mesh %s successfully", stat)
			ee.Details = fmt.Sprintf("Open service mesh is now %s.", stat)
			hh.StreamInfo(e)
		}(h, e)
	case
		common.BookInfoOperation,
		common.HTTPBinOperation,
		common.ImageHubOperation,
		common.EmojiVotoOperation:
		go func(hh *Handler, ee *adapter.Event) {
			appName := operations[request.OperationName].AdditionalProperties[common.ServiceName]
			stat, err := hh.installSampleApp(request.IsDeleteOperation, request.Namespace, operations[request.OperationName].Templates, kubeconfigs)
			if err != nil {
				summary := fmt.Sprintf("Error while %s %s application", stat, appName)	
				hh.streamErr(summary, e, err)
				return
			}
			ee.Summary = fmt.Sprintf("%s application %s successfully", appName, stat)
			ee.Details = fmt.Sprintf("The %s application is now %s.", appName, stat)
			hh.StreamInfo(e)
		}(h, e)
	case internalconfig.OSMBookStoreOperation:
		go func(hh *Handler, ee *adapter.Event) {
			version := string(operations[request.OperationName].Versions[0])
			appName := operations[request.OperationName].AdditionalProperties[common.ServiceName]
			stat, err := hh.installOSMBookStoreSampleApp(
				request.IsDeleteOperation,
				version,
				operations[request.OperationName].Templates,
				kubeconfigs,
			)
			if err != nil {
				summary := fmt.Sprintf("Error while %s %s application", stat, appName)
				hh.streamErr(summary, e, err)
				return
			}
			ee.Summary = fmt.Sprintf("%s application %s successfully", appName, stat)
			ee.Details = fmt.Sprintf("The %s application is now %s.", appName, stat)
			hh.StreamInfo(e)
		}(h, e)
	case common.SmiConformanceOperation:
		go func(hh *Handler, ee *adapter.Event) {
			name := operations[request.OperationName].Description
			_, err := hh.RunSMITest(adapter.SMITestOptions{
				Ctx:         context.TODO(),
				OperationID: ee.Operationid,
				Manifest:    string(operations[request.OperationName].Templates[0]),
				Namespace:   "meshery",
				Labels: map[string]string{
					"openservicemesh.io/monitored-by": "osm",
				},
				Annotations: make(map[string]string),
			})
			if err != nil {
				summary := fmt.Sprintf("Error while %s %s test", status.Running, name)	
				hh.streamErr(summary ,e, err)
				return
			}
			ee.Summary = fmt.Sprintf("%s test %s successfully", name, status.Completed)
			ee.Details = ""
			hh.StreamInfo(e)
		}(h, e)
	default:
		h.streamErr("Invalid operation", e, ErrOpInvalid)
	}
	return nil
}

func(h *Handler) streamErr(summary string, e *adapter.Event, err error) {
	e.Summary = summary
	e.Details = err.Error()
	e.ErrorCode = errors.GetCode(err)
	e.ProbableCause = errors.GetCause(err)
	e.SuggestedRemediation = errors.GetRemedy(err)
	h.StreamErr(e, err)
}

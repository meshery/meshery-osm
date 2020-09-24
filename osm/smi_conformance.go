package osm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/layer5io/meshery-osm/meshes"
	"github.com/layer5io/meshery-osm/osm/smi"
)

func (iClient *Client) validateSMIConformance(id string, version string) error {

	labels := map[string]string{
		"openservicemesh.io/monitored-by": "osm",
	}

	test, err := smi.New(context.TODO(), id, version, "Open Service Mesh", iClient.k8sClientset)
	if err != nil {
		iClient.eventChan <- &meshes.EventsResponse{
			OperationId: id,
			EventType:   meshes.EventType_ERROR,
			Summary:     "Error while creating smi-conformance tool",
			Details:     err.Error(),
		}
		return err
	}

	result, err := test.Run(labels, nil)
	if err != nil {
		iClient.eventChan <- &meshes.EventsResponse{
			OperationId: id,
			EventType:   meshes.EventType_ERROR,
			Summary:     "Error while Running smi-conformance test",
			Details:     err.Error(),
		}
		return err
	}
	jsondata, _ := json.Marshal(result)

	iClient.eventChan <- &meshes.EventsResponse{
		OperationId: id,
		EventType:   meshes.EventType_INFO,
		Summary:     fmt.Sprintf("Smi conformance test %s successfully", result.Status),
		Details:     string(jsondata),
	}

	return nil
}

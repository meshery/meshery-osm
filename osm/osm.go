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
	"context"
	"fmt"
	"time"

	"github.com/layer5io/meshery-osm/meshes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// CreateMeshInstance creates an instance of the mesh on the cluster
func (iClient *Client) CreateMeshInstance(_ context.Context, k8sReq *meshes.CreateMeshInstanceRequest) (*meshes.CreateMeshInstanceResponse, error) {
	var k8sConfig []byte
	contextName := ""
	if k8sReq != nil {
		k8sConfig = k8sReq.K8SConfig
		contextName = k8sReq.ContextName
	}
	// logrus.Debugf("received k8sConfig: %s", k8sConfig)
	logrus.Debugf("received contextName: %s", contextName)

	err := createKubeconfig(k8sConfig, contextName)
	if err != nil {
		err = errors.Wrapf(err, "unable to create kubeconfig")
		logrus.Error(err)
		return nil, err
	}

	ic, err := newClient(k8sConfig, contextName)
	if err != nil {
		err = errors.Wrapf(err, "unable to create a new osm client")
		logrus.Error(err)
		return nil, err
	}

	iClient.k8sClientset = ic.k8sClientset
	iClient.k8sDynamicClient = ic.k8sDynamicClient
	iClient.eventChan = make(chan *meshes.EventsResponse, 100)
	iClient.config = ic.config
	return &meshes.CreateMeshInstanceResponse{}, nil
}

// MeshName just returns the name of the mesh the client is representing
func (iClient *Client) MeshName(context.Context, *meshes.MeshNameRequest) (*meshes.MeshNameResponse, error) {
	return &meshes.MeshNameResponse{Name: "Open Service Mesh"}, nil
}

// ApplyOperation is a method invoked to apply a particular operation on the mesh in a namespace
func (iClient *Client) ApplyOperation(_ context.Context, arReq *meshes.ApplyRuleRequest) (*meshes.ApplyRuleResponse, error) {

	if arReq == nil {
		return nil, errors.New("mesh client has not been created")
	}

	_, ok := supportedOps[arReq.OpName]
	if !ok {
		return nil, fmt.Errorf("operation id: %s, error: %s is not a valid operation name", arReq.OperationId, arReq.OpName)
	}

	if arReq.OpName == customOpCommand && arReq.CustomBody == "" {
		return nil, fmt.Errorf("operation id: %s, error: yaml body is empty for %s operation", arReq.OperationId, arReq.OpName)
	}

	switch arReq.OpName {
	case installOSMCommand:
		go func() {
			opName := "deploying"
			if arReq.DeleteOp {
				opName = "removing"
			}
			if err := iClient.installMesh("osmctl", "v0.3.0", arReq.DeleteOp); err != nil {
				iClient.eventChan <- &meshes.EventsResponse{
					OperationId: arReq.OperationId,
					EventType:   meshes.EventType_ERROR,
					Summary:     fmt.Sprintf("Error while %s Open service mesh", opName),
					Details:     err.Error(),
				}
				return
			}
			opName = "deployed"
			if arReq.DeleteOp {
				opName = "removed"
			}
			iClient.eventChan <- &meshes.EventsResponse{
				OperationId: arReq.OperationId,
				EventType:   meshes.EventType_INFO,
				Summary:     fmt.Sprintf("Open service mesh %s successfully", opName),
				Details:     fmt.Sprintf("The Open service mesh is now %s.", opName),
			}

		}()
		return &meshes.ApplyRuleResponse{
			OperationId: arReq.OperationId,
		}, nil
	case smiConformanceCommand:
		go func() {
			err := iClient.validateSMIConformance(arReq.OperationId, "v0.3.0")
			if err != nil {
				return
			}
		}()
	}
	return &meshes.ApplyRuleResponse{}, nil
}

// SupportedOperations - returns a list of supported operations on the mesh
func (iClient *Client) SupportedOperations(context.Context, *meshes.SupportedOperationsRequest) (*meshes.SupportedOperationsResponse, error) {
	supportedOpsCount := len(supportedOps)
	result := make([]*meshes.SupportedOperation, supportedOpsCount)
	i := 0
	for k, sp := range supportedOps {
		result[i] = &meshes.SupportedOperation{
			Key:      k,
			Value:    sp.name,
			Category: sp.opType,
		}
		i++
	}
	return &meshes.SupportedOperationsResponse{
		Ops: result,
	}, nil
}

// StreamEvents - streams generated/collected events to the client
func (iClient *Client) StreamEvents(in *meshes.EventsRequest, stream meshes.MeshService_StreamEventsServer) error {
	for {
		select {
		case event := <-iClient.eventChan:
			logrus.Debugf("sending event: %+#v", event)
			if err := stream.Send(event); err != nil {
				err = errors.Wrapf(err, "unable to send event")

				// to prevent loosing the event, will re-add to the channel
				go func() {
					iClient.eventChan <- event
				}()
				logrus.Error(err)
				return err
			}
		default:
		}
		time.Sleep(500 * time.Millisecond)
	}
}

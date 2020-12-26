package osm

import (
	"fmt"
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshkit/errors"
	"github.com/layer5io/meshkit/utils"
)

func (h *Handler) installSampleApp(del bool, namespace string, templates []adapter.Template) (string, error) {
	st := status.Installing
	if del {
		st = status.Removing
	}
	for _, template := range templates {
		contents, err := utils.ReadRemoteFile(string(template))
		if err != nil {
			return st, ErrSampleApp(err)
		}
		err = h.applyManifest(del, namespace, []byte(contents))
		if err != nil {
			return st, ErrSampleApp(err)
		}
	}
	return status.Installed, nil
}

// ErrSampleApp is the error for streaming event
func ErrSampleApp(err error) error {
	return errors.NewDefault("osm_test_code", fmt.Sprintf("Error with sample app operation: %s", err.Error()))
}

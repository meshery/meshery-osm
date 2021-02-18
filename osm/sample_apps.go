package osm

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshkit/errors"
	"github.com/layer5io/meshkit/utils"
)

const (
	OperationDelete  = "delete"
	OperationInstall = "install"
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

func (h *Handler) OSMSampleBookBuyerExecute(remove bool, version string) (string, error) {
	var operation, st string

	fmt.Println("sample app request")
	executable, err := exec.LookPath("./scripts/install_sample_book_info_app.sh")
	if err != nil {
		return st, ErrSampleApp(err)
	}

	if remove {
		operation = OperationDelete
		st = status.Removing
	} else {
		operation = OperationInstall
		st = status.Installing
	}

	cmd := &exec.Cmd{
		Path:   executable,
		Args:   []string{executable, operation},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("OSM_VERSION=%s", version),
	)

	err = cmd.Start()
	if err != nil {
		return st, ErrSampleApp(err)
	}

	err = cmd.Wait()
	if err != nil {
		return st, err
	}

	if remove {
		st = status.Removed
	} else {
		st = status.Installed
	}

	return st, nil
}

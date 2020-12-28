package osm

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshkit/utils"
)

// noneNamespace indicates unset namespace
const noneNamespace = ""

// installOSMBookStoreSampleApp installs or uninstalls the default OSM bookstore application
func (h *Handler) installOSMBookStoreSampleApp(del bool, version string, templates []adapter.Template) (string, error) {
	st := status.Installing
	if del {
		st = status.Removing
	}

	namespaces := []string{
		"bookstore",
		"bookbuyer",
		"bookthief",
		"bookwarehouse",
	}

	// Add the namespaces for sidecar injection
	for _, ns := range namespaces {
		if err := createNS(h, ns, del); err != nil {
			return st, ErrCreatingNS(err)
		}
		if err := h.sidecarInjection(del, version, ns); err != nil {
			return st, ErrSidecarInjection(err)
		}
	}

	// Install the manifests
	st, err := h.installSampleApp(del, noneNamespace, templates)
	if err != nil {
		return st, err
	}

	return st, nil
}

func (h *Handler) installSampleApp(del bool, namespace string, templates []adapter.Template) (string, error) {
	st := status.Installing
	if del {
		st = status.Removing
	}
	for _, template := range templates {
		contents, err := utils.ReadFileSource(string(template))
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

// sidecarInjection enables/disables sidecar injection on a namespace
func (h *Handler) sidecarInjection(del bool, version, ns string) error {
	exe, err := h.getExecutable(version)
	if err != nil {
		return err
	}

	injectCmd := "add"
	if del {
		injectCmd = "remove"
	}

	cmd := &exec.Cmd{
		Path: exe,
		Args: []string{
			exe,
			"namespace",
			injectCmd,
			ns,
		},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	if err := cmd.Run(); err != nil {
		return ErrRunExecutable(err)
	}

	return nil
}

// createNS handles the creatin as well as deletion of namespaces
func createNS(h *Handler, ns string, del bool) error {
	manifest := fmt.Sprintf(`
apiVersion: v1
kind: Namespace
metadata:
  name: %s
`,
		ns,
	)

	if err := h.applyManifest(del, noneNamespace, []byte(manifest)); err != nil {
		return err
	}

	return nil
}

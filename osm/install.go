package osm

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// Versions of the Service Mesh
var (
	releases = map[string]struct{}{
		"v0.2.0": struct{}{},
		"v0.3.0": struct{}{},
	}
)

// installMesh installs the mesh
func (iClient *Client) installMesh(method string, version string) error {

	if _, ok := releases[version]; !ok {
		return errors.New(fmt.Sprintf("version %s unavailable", version))
	}

	// Install Method: Using Service Mesh-specific CLI
	switch method {
	case "osmctl":
		if err := applyOSM(version); err != nil {
			return err
		}
	}
	return nil
}

// applyOSM applies mesh resources with osmctl
func applyOSM(version string) error {
	Executable, err := exec.LookPath("./scripts/osmctl.sh")
	if err != nil {
		return err
	}

	cmd := &exec.Cmd{
		Path:   Executable,
		Args:   []string{Executable},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("OSM_VERSION=%s", version),
	)

	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

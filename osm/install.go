package osm

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshery-osm/internal/config"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

// Execute executes the install or delete operation using osmctl
func (h *Handler) Execute(del bool, version, ns string) (string, error) {
	if del {
		return h.deleteOSM(version, ns)
	}
	return h.installOSM(version, ns)
}

func (h *Handler) deleteOSM(version, ns string) (string, error) {
	st := status.Removing
	Executable, err := h.getExecutable(version)
	if err != nil {
		return st, (err)
	}

	cmd := &exec.Cmd{
		Path: Executable,
		Args: []string{
			Executable,
			"mesh",
			"uninstall",
			"-f",
			"--osm-namespace",
			ns,
		},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	if err = cmd.Run(); err != nil {
		return st, ErrInstallOSM(err)
	}

	return status.Removed, nil
}

func (h *Handler) installOSM(version, ns string) (string, error) {
	st := status.Installing
	Executable, err := h.getExecutable(version)
	if err != nil {
		return st, err
	}
	cmd := &exec.Cmd{
		Path: Executable,
		Args: []string{
			Executable,
			"install",
			"--osm-namespace",
			ns,
		},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	if err = cmd.Run(); err != nil {
		return st, ErrInstallOSM(err)
	}

	return status.Installed, nil
}

func (h *Handler) applyManifest(del bool, namespace string, contents []byte) error {
	err := h.MesheryKubeclient.ApplyManifest(contents, mesherykube.ApplyOptions{
		Namespace: namespace,
		Delete:    del,
	})
	if err != nil {
		return err
	}
	return nil
}

// getExecutable looks for the executable in
// 1. $PATH
// 2. Root config path
//
// If it doesn't find the executable in the path then it proceeds
// to download the binary from github releases and installs it
// in the root config path
func (h *Handler) getExecutable(release string) (string, error) {
	const platform = runtime.GOOS
	binaryName := generatePlatformSpecificBinaryName("osm", platform)
	alternateBinaryName := generatePlatformSpecificBinaryName("osm-"+release, platform)

	// Look for the executable in the path
	h.Log.Info("Looking for osm in the path...")
	executable, err := exec.LookPath(binaryName)
	if err == nil {
		return executable, nil
	}
	executable, err = exec.LookPath(alternateBinaryName)
	if err == nil {
		return executable, nil
	}

	binPath := path.Join(config.RootPath(), "bin")

	// Look for config in the root path
	h.Log.Info("Looking for osm in", binPath, "...")
	executable = path.Join(binPath, alternateBinaryName)
	if _, err := os.Stat(executable); err == nil {
		return executable, nil
	}

	// Proceed to download the binary in the config root path
	h.Log.Info("osm not found in the path, downloading...")
	res, err := downloadBinary(platform, runtime.GOARCH, release)
	if err != nil {
		return "", err
	}

	// Install the binary
	h.Log.Info("Installing...")
	if err = installBinary(
		res,
		binPath,
		platform,
		binaryName,
	); err != nil {
		return "", err
	}

	// Rename the binary
	if err = os.Rename(
		path.Join(binPath, fmt.Sprintf("%s-%s", platform, runtime.GOARCH), binaryName),
		path.Join(binPath, alternateBinaryName),
	); err != nil {
		return "", ErrInstallBinary(err)
	}

	// Cleanup
	if err = os.RemoveAll(path.Join(binPath, fmt.Sprintf("%s-%s", platform, runtime.GOARCH))); err != nil {
		return "", ErrInstallBinary(err)
	}

	h.Log.Info("Done")
	return path.Join(binPath, alternateBinaryName), nil
}

func downloadBinary(platform, arch, release string) (*http.Response, error) {
	url := fmt.Sprintf(
		"https://github.com/openservicemesh/osm/releases/download/%s/osm-%s-%s-%s.tar.gz",
		release,
		release,
		platform,
		arch,
	)

	// we need variable url here hence,
	// #nosec
	resp, err := http.Get(url)
	if err != nil {
		return nil, ErrDownloadBinary(err)
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, ErrDownloadBinary(fmt.Errorf("bad status: %s", resp.Status))
	}

	return resp, nil
}

func installBinary(res *http.Response, location, platform, name string) error {
	// Close the response body
	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	err := os.MkdirAll(location, 0750)
	if err != nil {
		return err
	}

	untar, err := tarxzf(res.Body, location)
	if err != nil {
		return ErrInstallBinary(err)
	}

	if platform == "darwin" || platform == "linux" {
		// Change permissions, we need the binary to be executable, hence
		// #nosec
		if err = os.Chmod(path.Join(location, untar, name), 0750); err != nil {
			return ErrInstallBinary(err)
		}
	}

	return nil
}

func tarxzf(stream io.Reader, location string) (string, error) {
	uncompressedStream, err := gzip.NewReader(stream)
	if err != nil {
		return "", err
	}

	tarReader := tar.NewReader(uncompressedStream)
	name := ""

	for {
		header, err := tarReader.Next()
		if name == "" {
			name = header.Name
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return "", ErrTarXZF(err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			// File traversal is required to store the binary at the right place
			// #nosec
			if err := os.MkdirAll(path.Join(location, header.Name), 0750); err != nil {
				return "", ErrTarXZF(err)
			}
		case tar.TypeReg:
			// File traversal is required to store the binary at the right place
			// #nosec
			outFile, err := os.Create(path.Join(location, header.Name))
			if err != nil {
				return "", ErrTarXZF(err)
			}
			// Trust istioctl tar
			// #nosec
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return "", ErrTarXZF(err)
			}
			if err = outFile.Close(); err != nil {
				return "", ErrTarXZF(err)
			}

		default:
			return "", ErrTarXZF(err)
		}
	}

	return name, nil
}

func generatePlatformSpecificBinaryName(binName, platform string) string {
	if platform == "windows" && !strings.HasSuffix(binName, ".exe") {
		return binName + ".exe"
	}

	return binName
}

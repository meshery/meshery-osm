package osm

import (
	"context"
	"fmt"
	"sync"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// noneNamespace indicates unset namespace
const noneNamespace = ""

// installOSMBookStoreSampleApp installs or uninstalls the default OSM bookstore application
func (h *Handler) installOSMBookStoreSampleApp(del bool, version string, templates []adapter.Template, kubeconfigs []string) (string, error) {
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
		if err := createNS(h, ns, del, kubeconfigs); err != nil {
			return st, ErrCreatingNS(err)
		}
		if err := h.sidecarInjection(ns, del, kubeconfigs); err != nil {
			return st, ErrSidecarInjection(err)
		}
	}

	// Install the manifests
	st, err := h.installSampleApp(del, noneNamespace, templates, kubeconfigs)
	if err != nil {
		return st, err
	}

	return st, nil
}

func (h *Handler) installSampleApp(del bool, namespace string, templates []adapter.Template, kubeconfigs []string) (string, error) {
	st := status.Installing
	if del {
		st = status.Removing
	}
	for _, template := range templates {
		err := h.applyManifest([]byte(template.String()), del, namespace, kubeconfigs)
		if err != nil {
			return st, ErrSampleApp(err)
		}
	}
	return status.Installed, nil
}

// sidecarInjection enables/disables sidecar injection on a namespace
func (h *Handler) sidecarInjection(namespace string, del bool, kubeconfigs []string) error {
	var wg sync.WaitGroup
	var errs []error
	var errMx sync.Mutex
	for _, k8sconfig := range kubeconfigs {
		wg.Add(1)
		go func(k8sconfig string) {
			defer wg.Done()
			kClient, err := mesherykube.New([]byte(k8sconfig))
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}
			// updating the label on the namespace
			ns, err := kClient.KubeClient.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}

			if ns.ObjectMeta.Labels == nil {
				ns.ObjectMeta.Labels = map[string]string{}
			}
			ns.ObjectMeta.Labels["openservicemesh.io/monitored-by"] = "osm"

			if del {
				delete(ns.ObjectMeta.Labels, "openservicemesh.io/monitored-by")
			}

			// updating the annotations on the namespace
			if ns.ObjectMeta.Annotations == nil {
				ns.ObjectMeta.Annotations = map[string]string{}
			}
			ns.ObjectMeta.Annotations["openservicemesh.io/sidecar-injection"] = "enabled"

			if del {
				delete(ns.ObjectMeta.Annotations, "openservicemesh.io/sidecar-injection")
			}

			fmt.Println(ns.ObjectMeta)

			_, err = kClient.KubeClient.CoreV1().Namespaces().Update(context.TODO(), ns, metav1.UpdateOptions{})
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}
		}(k8sconfig)
	}
	wg.Wait()
	if len(errs) != 0 {
		return ErrLoadNamespace(mergeErrors(errs), namespace)
	}
	return nil
}

// createNS handles the creatin as well as deletion of namespaces
func createNS(h *Handler, ns string, del bool, kubeconfigs []string) error {
	manifest := fmt.Sprintf(`
apiVersion: v1
kind: Namespace
metadata:
  name: %s
`,
		ns,
	)

	if err := h.applyManifest([]byte(manifest), del, noneNamespace, kubeconfigs); err != nil {
		return err
	}

	return nil
}

func (h *Handler) applyManifest(contents []byte, isDel bool, namespace string, kubeconfigs []string) error {
	var wg sync.WaitGroup
	var errs []error
	var errMx sync.Mutex
	for _, k8sconfig := range kubeconfigs {
		wg.Add(1)
		go func(k8sconfig string) {
			defer wg.Done()
			kClient, err := mesherykube.New([]byte(k8sconfig))
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}
			err = kClient.ApplyManifest(contents, mesherykube.ApplyOptions{
				Namespace: namespace,
				Update:    true,
				Delete:    isDel,
			})
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}
		}(k8sconfig)
	}
	wg.Wait()
	if len(errs) != 0 {
		return ErrLoadNamespace(mergeErrors(errs), namespace)
	}
	return nil
}

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
	"fmt"
	"io/ioutil"
	"os/user"
	"time"

	"github.com/layer5io/meshery-osm/meshes"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	// auth is needed for initialization only
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	// "github.com/ghodss/yaml"
	"gopkg.in/yaml.v2"
)

var (
	kubeconfigpath = ".kube/config"
)

// Client represents an osm client in Meshery
type Client struct {
	config           *rest.Config
	k8sClientset     *kubernetes.Clientset
	k8sDynamicClient dynamic.Interface
	eventChan        chan *meshes.EventsResponse

	osmReleaseVersion     string
	osmReleaseDownloadURL string
	osmReleaseUpdatedAt   time.Time

	smiAddress string
}

type Kubeconfig struct {
	APIVersion string `yaml:"apiVersion,omitempty"`
	Clusters   []struct {
		Cluster struct {
			CertificateAuthorityData string `yaml:"certificate-authority-data,omitempty"`
			Server                   string `yaml:"server,omitempty"`
		} `yaml:"cluster,omitempty"`
		Name string `yaml:"name,omitempty"`
	} `yaml:"clusters,omitempty"`
	Contexts []struct {
		Context struct {
			Cluster   string `yaml:"cluster,omitempty"`
			Namespace string `yaml:"namespace,omitempty"`
			User      string `yaml:"user,omitempty"`
		} `yaml:"context,omitempty"`
		Name string `yaml:"name,omitempty"`
	} `yaml:"contexts,omitempty"`
	CurrentContext string `yaml:"current-context,omitempty"`
	Kind           string `yaml:"kind,omitempty"`
	Preferences    struct {
	} `yaml:"preferences,omitempty"`
	Users []struct {
		Name string `yaml:"name,omitempty"`
		User struct {
			Exec struct {
				APIVersion string   `yaml:"apiVersion,omitempty"`
				Args       []string `yaml:"args,omitempty"`
				Command    string   `yaml:"command,omitempty"`
				Env        []struct {
					Name  string `yaml:"name,omitempty"`
					Value string `yaml:"value,omitempty"`
				} `yaml:"env,omitempty"`
			} `yaml:"exec,omitempty"`
			AuthProvider struct {
				Config struct {
					AccessToken string    `yaml:"access-token,omitempty"`
					CmdArgs     string    `yaml:"cmd-args,omitempty"`
					CmdPath     string    `yaml:"cmd-path,omitempty"`
					Expiry      time.Time `yaml:"expiry,omitempty"`
					ExpiryKey   string    `yaml:"expiry-key,omitempty"`
					TokenKey    string    `yaml:"token-key,omitempty"`
				} `yaml:"config,omitempty"`
				Name string `yaml:"name,omitempty"`
			} `yaml:"auth-provider,omitempty"`
			ClientCertificateData string `yaml:"client-certificate-data,omitempty"`
			ClientKeyData         string `yaml:"client-key-data,omitempty"`
			Token                 string `yaml:"token,omitempty"`
		} `yaml:"user,omitempty,omitempty"`
	} `yaml:"users,omitempty"`
}

func configClient(kubeconfig []byte, contextName string) (*rest.Config, error) {
	if len(kubeconfig) > 0 {
		ccfg, err := clientcmd.Load(kubeconfig)
		if err != nil {
			return nil, err
		}
		if contextName != "" {
			ccfg.CurrentContext = contextName
		}

		return clientcmd.NewDefaultClientConfig(*ccfg, &clientcmd.ConfigOverrides{}).ClientConfig()
	}
	return rest.InClusterConfig()
}

func newClient(kubeconfig []byte, contextName string) (*Client, error) {
	kubeconfig = monkeyPatchingToSupportInsecureConn(kubeconfig)
	client := Client{}
	config, err := configClient(kubeconfig, contextName)
	if err != nil {
		return nil, err
	}
	config.QPS = 100
	config.Burst = 200

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	client.k8sDynamicClient = dynamicClient

	k8sClientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	client.k8sClientset = k8sClientset
	client.config = config

	return &client, nil
}

func monkeyPatchingToSupportInsecureConn(data []byte) []byte {
	config := map[string]interface{}{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		logrus.Warn(err)
		return data // we will skip this process
	}
	// logrus.Infof("unmarshalled config: %+#v", config)
	clusters, ok := config["clusters"].([]interface{})
	if !ok {
		logrus.Warn("unable to type cast clusters to a map array")
		return data
	}
	for _, clusterI := range clusters {
		cluster, ok := clusterI.(map[string]interface{})
		if !ok {
			logrus.Warn("unable to type case individual cluster to a map")
			continue
		}
		indCluster, ok := cluster["cluster"].(map[string]interface{})
		if !ok {
			logrus.Warn("unable to type case clusters.cluster to a map")
			continue
		}
		indCluster["insecure-skip-tls-verify"] = true // TODO: should we persist this back?
		delete(indCluster, "certificate-authority-data")
		delete(indCluster, "certificate-authority")
	}
	// logrus.Debugf("New config: %+#v", config)
	data1, err := yaml.Marshal(config)
	if err != nil {
		logrus.Warn(err)
		return data
	}
	return data1
}

func createKubeconfig(kubeconfig []byte, contextName string) error {

	yamlConfig := Kubeconfig{}
	err := yaml.Unmarshal(kubeconfig, &yamlConfig)
	if err != nil {
		return err
	}

	yamlConfig.CurrentContext = contextName

	d, err := yaml.Marshal(yamlConfig)
	if err != nil {
		return err
	}

	user, err := user.Current()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", user.HomeDir, kubeconfigpath), d, 0644)
	if err != nil {
		return err
	}

	return nil
}

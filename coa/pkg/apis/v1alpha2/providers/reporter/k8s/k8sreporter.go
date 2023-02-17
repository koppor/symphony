/*
MIT License

Copyright (c) Microsoft Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE
*/

package k8s

import (
	"context"
	"encoding/json"
	"path/filepath"

	"strconv"

	"github.com/azure/symphony/coa/pkg/apis/v1alpha2"
	"github.com/azure/symphony/coa/pkg/apis/v1alpha2/contexts"
	"github.com/azure/symphony/coa/pkg/apis/v1alpha2/providers"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8sReporterConfig struct {
	Name       string `json:"name"`
	ConfigPath string `json:"configPath"`
	InCluster  bool   `json:"inCluster"` //TODO: add context support
}

func K8sReporterConfigFromMap(properties map[string]string) (K8sReporterConfig, error) {
	ret := K8sReporterConfig{}
	if v, ok := properties["name"]; ok {
		ret.Name = v //providers.LoadEnv(v)
	}
	if v, ok := properties["configPath"]; ok {
		ret.ConfigPath = v //providers.LoadEnv(v)
	}
	if v, ok := properties["inCluster"]; ok {
		val := v //providers.LoadEnv(v)
		if val != "" {
			bVal, err := strconv.ParseBool(val)
			if err != nil {
				return ret, v1alpha2.NewCOAError(err, "invalid bool value in the 'inCluster' setting of K8s reporter", v1alpha2.BadConfig)
			}
			ret.InCluster = bVal
		}
	}
	return ret, nil
}

func (i *K8sReporter) InitWithMap(properties map[string]string) error {
	config, err := K8sReporterConfigFromMap(properties)
	if err != nil {
		return err
	}
	return i.Init(config)
}

type K8sReporter struct {
	Config        K8sReporterConfig
	Client        *kubernetes.Clientset
	DynamicClient dynamic.Interface
	Context       *contexts.ManagerContext
}

func (m *K8sReporter) ID() string {
	return m.Config.Name
}

func (a *K8sReporter) SetContext(context *contexts.ManagerContext) error {
	a.Context = context
	return nil
}

func (m *K8sReporter) Init(config providers.IProviderConfig) error {
	var err error
	aConfig, err := toK8sReporterConfig(config)
	if err != nil {
		return v1alpha2.NewCOAError(nil, "provided config is not a valid k8s reporter config", v1alpha2.BadConfig)
	}
	m.Config = aConfig
	var kConfig *rest.Config

	if m.Config.InCluster {
		kConfig, err = rest.InClusterConfig()
	} else {
		if m.Config.ConfigPath == "" {
			if home := homedir.HomeDir(); home != "" {
				m.Config.ConfigPath = filepath.Join(home, ".kube", "config")
			} else {
				return v1alpha2.NewCOAError(nil, "can't locate home direction to read default kubernetes config file, to run in cluster, set inCluster config setting to true", v1alpha2.BadConfig)
			}
		}
		kConfig, err = clientcmd.BuildConfigFromFlags("", m.Config.ConfigPath)
	}
	if err != nil {
		return err
	}
	m.Client, err = kubernetes.NewForConfig(kConfig)
	if err != nil {
		return err
	}
	m.DynamicClient, err = dynamic.NewForConfig(kConfig)
	if err != nil {
		return err
	}
	return nil
}

func toK8sReporterConfig(config providers.IProviderConfig) (K8sReporterConfig, error) {
	ret := K8sReporterConfig{}
	data, err := json.Marshal(config)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal(data, &ret)
	//ret.Name = providers.LoadEnv(ret.Name)
	//ret.ConfigPath = providers.LoadEnv(ret.ConfigPath)
	return ret, err
}

func (m *K8sReporter) Report(id string, namespace string, group string, kind string, version string, properties map[string]string, overwrtie bool) error {

	obj, err := m.DynamicClient.Resource(schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: kind,
	}).Namespace(namespace).Get(context.Background(), id, v1.GetOptions{})

	if err != nil {
		return err
	}

	propCol := make(map[string]string)

	if !overwrtie {
		if existingStatus, ok := obj.Object["status"]; ok {
			dict := existingStatus.(map[string]interface{})
			if propsElement, ok := dict["properties"]; ok {
				props := propsElement.(map[string]interface{})
				for k, v := range props {
					propCol[k] = v.(string)
				}
			}
		}
	}

	for k, v := range properties {
		propCol[k] = v
	}

	status := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       "Status",
			"metadata": map[string]interface{}{
				"name": id,
			},
			"status": map[string]interface{}{
				"properties": propCol,
			},
		},
	}

	status.SetResourceVersion(obj.GetResourceVersion())

	_, err = m.DynamicClient.Resource(schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: kind,
	}).Namespace(namespace).UpdateStatus(context.Background(), status, v1.UpdateOptions{})

	return err
}

func (a *K8sReporter) Clone(config providers.IProviderConfig) (providers.IProvider, error) {
	ret := &K8sReporter{}
	if config == nil {
		err := ret.Init(a.Config)
		if err != nil {
			return nil, err
		}
	} else {
		err := ret.Init(config)
		if err != nil {
			return nil, err
		}
	}
	if a.Context != nil {
		ret.Context = a.Context
	}
	return ret, nil
}

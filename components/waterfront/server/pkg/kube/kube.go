/*
Copyright 2018 Pax Automa Systems, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kube

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewKubeClient(url, kubeConfig string) (*kubernetes.Clientset, error) {
	var (
		err    error
		config *rest.Config
	)

	if kubeConfig == "" {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags(url, kubeConfig)
	}
	if err != nil {
		return nil, errors.Wrap(err, "could not load kube config")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "could not create kube clientset")
	}

	return clientset, nil
}

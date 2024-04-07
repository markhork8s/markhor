package main

import (
	"time"

	v1 "sops_k8s/api/types/v1"
	client_v1 "sops_k8s/clientset/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

func WatchResources(clientSet client_v1.ExampleV1Interface) cache.Store {
	sopsSecretStore, sopsSecretController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return clientSet.SopsSecrets("some-namespace").List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.SopsSecrets("some-namespace").Watch(lo)
			},
		},
		&v1.SopsSecret{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{},
	)

	go sopsSecretController.Run(wait.NeverStop)
	return sopsSecretStore
}

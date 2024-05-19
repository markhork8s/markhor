package pkg

import (
	"time"

	v1 "github.com/markhork8s/markhor/pkg/api/types/v1"
	client_v1 "github.com/markhork8s/markhor/pkg/clientset/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

func WatchResources(client client_v1.MarkhorV1Interface) cache.Store {
	markhorSecretStore, markhorSecretController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return client.MarkhorSecrets().List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return client.MarkhorSecrets().Watch(lo)
			},
		},
		&v1.MarkhorSecret{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{},
	)

	go markhorSecretController.Run(wait.NeverStop)
	return markhorSecretStore
}

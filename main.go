package main

import (
	"log"

	"github.com/civts/markhor/pkg"
	apiV1 "github.com/civts/markhor/pkg/api/types/v1"
	cs "github.com/civts/markhor/pkg/clientset"

	"k8s.io/client-go/kubernetes/scheme"
)

func init() {
	apiV1.AddToScheme(scheme.Scheme)
}

func main() {
	log.Println("Starting Markhor")
	config := pkg.ParseConfig()

	mClient, clientset := cs.GetK8sClients(config.Kubernetes.KubeconfigPath)

	go pkg.SetupHealthcheck(config.Healthcheck)

	cs.WatchMarkhorSecrets(mClient, clientset)
}

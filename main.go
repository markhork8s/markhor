package main

import (
	apiV1 "github.com/civts/markhor/pkg/api/types/v1"
	cs "github.com/civts/markhor/pkg/clientset"
	"github.com/civts/markhor/pkg/config"
	"github.com/civts/markhor/pkg/healthcheck"

	"k8s.io/client-go/kubernetes/scheme"
)

func init() {
	apiV1.AddToScheme(scheme.Scheme)
}

func main() {
	config := config.ParseConfig()

	mClient, clientset := cs.GetK8sClients(config.Kubernetes.KubeconfigPath)

	go healthcheck.SetupHealthcheck(config.Healthcheck)

	cs.WatchMarkhorSecrets(mClient, clientset, config)
}

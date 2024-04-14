package main

import (
	"flag"

	"github.com/civts/markhor/pkg"
	apiV1 "github.com/civts/markhor/pkg/api/types/v1"
	cs "github.com/civts/markhor/pkg/clientset"

	"k8s.io/client-go/kubernetes/scheme"
)

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
	flag.Parse()
}

func main() {
	apiV1.AddToScheme(scheme.Scheme)

	mClient, clientset := cs.GetK8sClients(kubeconfig)

	go pkg.SetupHealthcheck(8080)

	cs.WatchMarkhorSecrets(mClient, clientset)
}

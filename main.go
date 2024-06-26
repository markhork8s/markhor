package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/markhork8s/markhor/pkg/admission_controller"
	apiV1 "github.com/markhork8s/markhor/pkg/api/types/v1"
	cs "github.com/markhork8s/markhor/pkg/clientset"
	"github.com/markhork8s/markhor/pkg/config"
	"github.com/markhork8s/markhor/pkg/healthcheck"

	"k8s.io/client-go/kubernetes/scheme"
)

func init() {
	err := apiV1.AddToScheme(scheme.Scheme)
	if err != nil {
		fmt.Println("Could not create the Kubernetes scheme for MarkhorSecrets. Exiting")
		os.Exit(1)
	}
}

func main() {
	config, err := config.ParseConfig()
	if err != nil {
		slog.Error(fmt.Sprint("Something went wrong parsing the configuration: ", err))
		os.Exit(1)
	}

	mClient, clientset := cs.GetK8sClients(config.Kubernetes.KubeconfigPath)

	go healthcheck.SetupHealthcheck(config)
	go admission_controller.SetupAdmissionController(config)
	go cs.WatchMarkhorSecrets(mClient, clientset, config)

	setupGracefulShutdown()
}

func setupGracefulShutdown() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown

	slog.Info("Termination signal received. Shutting down...")
	slog.Info("Goodbye!")
	os.Exit(0)
}

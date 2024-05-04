package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/civts/markhor/pkg/admission_controller"
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
	config, err := config.ParseConfig()
	if err != nil {
		slog.Error("Something went wrong parsing the configuration: ", err)
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

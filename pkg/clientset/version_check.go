package clientset

import (
	"fmt"
	"log/slog"
	"os"

	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
)

func VersionCheck(clientset *kubernetes.Clientset) {
	version, err := clientset.Discovery().ServerVersion()
	if err != nil {
		slog.Error(fmt.Sprint("Could not get the version of the Kubernetes cluster we are talking to. Exiting.", err))
		os.Exit(1)
	}

	if IsCompatible(version) {
		slog.Info(fmt.Sprintf("The Kubernetes cluster is running version: %s", version.String()))
	} else {
		slog.Warn(fmt.Sprintf("The version of your Kubernetes cluster is %s. We did not test it's compatibility with this release of Markhor. Consider opening an issue to let us know if things are working as expected or if you experience some problems.", version.String()))
	}
}

func IsCompatible(version *version.Info) bool {
	if version.Major == "1" {
		if version.Minor == "28" {
			return true
		}
	}
	return false
}

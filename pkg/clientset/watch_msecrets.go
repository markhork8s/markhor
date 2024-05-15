package clientset

import (
	"log/slog"

	"fmt"
	"os"
	"strings"
	"time"

	"github.com/civts/markhor/pkg"
	apiV1 "github.com/civts/markhor/pkg/api/types/v1"
	v1 "github.com/civts/markhor/pkg/clientset/v1"
	"github.com/civts/markhor/pkg/config"
	"github.com/civts/markhor/pkg/handlers"
	"github.com/civts/markhor/pkg/healthcheck"
	"github.com/google/uuid"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

var connectedSuccessfully = false

func WatchMarkhorSecrets(mClient *v1.MarkhorV1Client, k8sClient *kubernetes.Clientset, config *config.Config) {
	go checkConnectTimeout(config.Kubernetes.ClusterTimeoutSeconds)
	markhorSecrets, err := mClient.MarkhorSecrets().Watch(metav1.ListOptions{})
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "the server could not find the requested resource") {
			slog.Error("Kubernetes does not know what a MarkhorSecret is. Did you forget to install the CRD?")
			os.Exit(1)
		}
		panic(err)
	}
	connectedSuccessfully = true
	channel := markhorSecrets.ResultChan()
	slog.Info("Started watching the events in the cluster")
	healthcheck.Healthy = true
	excludedNamespaces := arrayToSet(config.Behavior.ExcludedNamespaces)
	includedNamespaces := arrayToSet(config.Behavior.Namespaces)
	for event := range channel {
		eventId := uuid.New()
		eid := slog.String(pkg.SLOG_EVENT_ID_KEY, eventId.String())
		markhorSecret, ok := event.Object.(*apiV1.MarkhorSecret)
		namespace := markhorSecret.ObjectMeta.Namespace
		// Check that the MarkhorSecret belongs to a namespace we are supposed to manage (else, skip it)
		if _, present := excludedNamespaces[namespace]; present {
			slog.Debug(fmt.Sprint("Got an event for a MarkhorSecret in a namespace that was excluded: ", namespace))
			continue
		} else if len(includedNamespaces) != 0 {
			_, present := includedNamespaces[namespace]
			if !present {
				slog.Debug(fmt.Sprint("Got an event for a MarkhorSecret in a namespace that was not included in the ones we should watch: ", namespace))
				continue
			}
		}
		secretName := fmt.Sprintf("%s/%s", namespace, markhorSecret.ObjectMeta.Name)
		if !ok {
			slog.Debug("Failed to cast the object to type MarkhorSecret")
			continue
		}
		args := handlers.HandlerAttrs{
			MarkhorSecret: markhorSecret,
			EventId:       eid,
			Clientset:     k8sClient,
			Config:        config,
		}
		switch event.Type {
		case watch.Added:
			slog.Info(fmt.Sprint("A MarkhorSecret was added: ", secretName), eid)
			res := handlers.HandleAddition(args)
			if res {
				slog.Info(fmt.Sprint("New secret created correctly: ", secretName), eid)
			}
		case watch.Modified:
			slog.Info(fmt.Sprint("A MarkhorSecret was updated: ", secretName), eid)
			res := handlers.HandleAddition(args)
			if res {
				slog.Info(fmt.Sprint("Secret updated correctly: ", secretName), eid)
			}
		case watch.Deleted:
			slog.Info(fmt.Sprint("A MarkhorSecret was deleted: ", secretName), eid)
			handlers.HandleDeletion(args)
		}
	}
	healthcheck.Healthy = false
	slog.Warn("Finished watching the events in the cluster. Most probably, the channel was closed")
}

func checkConnectTimeout(timeout int) {
	slog.Info("Connecting to the k8s cluster")

	for i := 1; i <= timeout; i++ {
		time.Sleep(1 * time.Second)
		if connectedSuccessfully {
			return
		} else if i == 2 {
			slog.Info("No response from the k8s cluster. Will retry until the timeout")
		}
	}

	slog.Error(fmt.Sprintf("Connecting to the k8s cluster timed out after %d seconds. Check the kubeconfig file and that the cluster is up.", timeout))

	os.Exit(1)
}

func arrayToSet(stringsArray []string) map[string]struct{} {
	stringSet := make(map[string]struct{})
	for _, str := range stringsArray {
		stringSet[str] = struct{}{}
	}
	return stringSet
}

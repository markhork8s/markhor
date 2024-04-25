package handlers

import (
	"log/slog"

	"github.com/civts/markhor/pkg/config"

	v1 "github.com/civts/markhor/pkg/api/types/v1"

	"k8s.io/client-go/kubernetes"
)

// Arguments to pass to the handlers of MarkhorSecret events
type HandlerAttrs struct {
	//The Markhor Secret being processed
	MarkhorSecret *v1.MarkhorSecret
	//Identifier of this event
	EventId slog.Attr
	//Client to talk with kubernetes
	Clientset *kubernetes.Clientset
	//The configuration of Markhor
	Config *config.Config
}

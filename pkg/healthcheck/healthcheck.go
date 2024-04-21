package healthcheck

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/civts/markhor/pkg/config"
)

var Healthy bool = false

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if Healthy {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Alive and well\n"))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something is wrong\n"))
	}
}

func SetupHealthcheck(conf config.HealthcheckConfig) {
	if conf.Enabled {
		slog.Debug(fmt.Sprint("Healthcheck endpoint created on port ", conf.Port))
		http.HandleFunc("/health", healthCheckHandler)
		http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
	} else {
		slog.Debug("Skipping healthcheck -disabled in the config-")
	}
}

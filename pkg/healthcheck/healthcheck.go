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

const healthcheckEndpoint = "/health"

func SetupHealthcheck(cf *config.Config) {
	ch := cf.Healthcheck
	if ch.Enabled {
		http.HandleFunc(healthcheckEndpoint, healthCheckHandler)
		if cf.Tls.Mode == config.TLSExternalMode {
			slog.Debug(fmt.Sprint("Healthcheck http endpoint created on port ", ch.Port))
			http.ListenAndServe(fmt.Sprintf(":%d", ch.Port), nil)
		} else {
			slog.Debug(fmt.Sprint("Healthcheck https endpoint created on port ", ch.Port))
			http.ListenAndServeTLS(fmt.Sprintf(":%d", ch.Port), cf.Tls.CertPath, cf.Tls.KeyPath, nil)
		}
	} else {
		slog.Debug("Skipping healthcheck -disabled in the config-")
	}
}

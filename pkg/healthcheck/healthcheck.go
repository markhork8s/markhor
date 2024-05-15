package healthcheck

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/civts/markhor/pkg"
	"github.com/civts/markhor/pkg/config"
)

var Healthy bool = false

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	if Healthy {
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("Alive and well\n"))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("Something is wrong\n"))
	}
	if err != nil {
		slog.Warn(fmt.Sprint("Could not write healthcheck response: ", err))
	}
}

const healthcheckEndpoint = "/health"

func SetupHealthcheck(cf *config.Config) {
	ch := cf.Healthcheck
	if ch.Enabled {
		mux := http.NewServeMux()
		mux.HandleFunc(healthcheckEndpoint, healthCheckHandler)
		server := &http.Server{
			Addr:         fmt.Sprintf(":%d", ch.Port),
			Handler:      mux,
			ReadTimeout:  pkg.SERVER_READ_TIMEOUT_SECONDS * time.Second,
			WriteTimeout: pkg.SERVER_WRITE_TIMEOUT_SECONDS * time.Second,
		}
		var err error
		if cf.Tls.Enabled {
			slog.Debug(fmt.Sprint("Healthcheck https endpoint created on port ", ch.Port))
			err = server.ListenAndServeTLS(cf.Tls.CertPath, cf.Tls.KeyPath)
		} else {
			slog.Debug(fmt.Sprint("Healthcheck http endpoint created on port ", ch.Port))
			err = server.ListenAndServe()
		}
		if err != nil {
			slog.Error(fmt.Sprint("Could not start the healthcheck listener on port ", ch.Port, ": ", err))
			os.Exit(1)
		}
	} else {
		slog.Debug("Skipping healthcheck -disabled in the config-")
	}
}

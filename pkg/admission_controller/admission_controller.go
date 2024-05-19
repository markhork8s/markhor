package admission_controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/markhork8s/markhor/pkg"
	"github.com/markhork8s/markhor/pkg/config"
)

func SetupAdmissionController(conf *config.Config) {
	acConfig := conf.AdmissionController
	if acConfig.Enabled {
		mux := http.NewServeMux()
		vh := ValidateHandler{
			config: conf.MarkorSecrets,
		}
		mux.HandleFunc("/validate", vh.handler)
		server := &http.Server{
			Addr:         fmt.Sprintf(":%d", acConfig.Port),
			Handler:      mux,
			ReadTimeout:  pkg.SERVER_READ_TIMEOUT_SECONDS * time.Second,
			WriteTimeout: pkg.SERVER_WRITE_TIMEOUT_SECONDS * time.Second,
		}
		var err error
		if conf.Tls.Enabled {
			err = server.ListenAndServeTLS(conf.Tls.CertPath, conf.Tls.KeyPath)
		} else {
			err = server.ListenAndServe()
		}
		if err != nil {
			slog.Error(fmt.Sprint("Could not start the admission controller on port ", acConfig.Port, ": ", err))
			os.Exit(1)
		}
	}
}

package admission_controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/civts/markhor/pkg/config"
)

func SetupAdmissionController(conf *config.Config) {
	acConfig := conf.AdmissionController
	if acConfig.Enabled {
		http.HandleFunc("/validate", validateHandler)
		var err error
		if conf.Tls.Mode == config.TLSExternalMode {
			err = http.ListenAndServe(fmt.Sprintf(":%d", acConfig.Port), nil)
		} else {
			err = http.ListenAndServeTLS(fmt.Sprintf(":%d", acConfig.Port), conf.Tls.CertPath, conf.Tls.KeyPath, nil)
		}
		if err != nil {
			slog.Error(fmt.Sprint("Could not start the admission controller on port ", acConfig.Port, ": ", err))
			os.Exit(1)
		}
	}
}

package pkg

import (
	"fmt"
	"log"
	"net/http"
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

func SetupHealthcheck(conf HealthcheckConfig) {
	if conf.Enabled {
		log.Println("Healthcheck endpoint created on port", conf.Port)
		http.HandleFunc("/health", healthCheckHandler)
		http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
	} else {
		log.Println("Skipping healthcheck -disabled in the config-")
	}
}

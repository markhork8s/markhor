package pkg

import (
	"fmt"
	"net/http"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Alive and well\n"))
}

func SetupHealthcheck(port int) {
	fmt.Println("Starting the healthcheck")
	http.HandleFunc("/health", healthCheckHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

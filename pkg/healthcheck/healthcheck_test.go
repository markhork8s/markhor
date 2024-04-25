package healthcheck

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/civts/markhor/pkg/config"
)

func TestHealthcheckHandler_Healthy(t *testing.T) {
	req, err := http.NewRequest("GET", healthcheckEndpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthCheckHandler)

	prevHealthy := Healthy
	defer func() { Healthy = prevHealthy }()
	Healthy = true
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Alive and well\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHealthcheckHandler_Unhealthy(t *testing.T) {
	req, err := http.NewRequest("GET", healthcheckEndpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthCheckHandler)

	prevHealthy := Healthy
	defer func() { Healthy = prevHealthy }()
	Healthy = false
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := "Something is wrong\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestSetupHealthcheckEnabled_DefaultPort_Unhealthy(t *testing.T) {
	conf := config.HealthcheckConfig{
		Enabled: true,
		Port:    32714,
	}

	go SetupHealthcheck(conf)
	time.Sleep(time.Millisecond * 300)

	req, err := http.NewRequest("GET", "http://localhost:"+fmt.Sprint(conf.Port)+healthcheckEndpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code %v, got %v", http.StatusInternalServerError, resp.StatusCode)
	}
}

func TestSetupHealthcheckDisabled_RequestFail(t *testing.T) {
	conf := config.HealthcheckConfig{
		Enabled: false,
		Port:    8080,
	}

	go SetupHealthcheck(conf)
	time.Sleep(time.Millisecond * 300)

	req, err := http.NewRequest("GET", "http://localhost:"+fmt.Sprint(conf.Port)+healthcheckEndpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = http.DefaultClient.Do(req)
	if err == nil {
		t.Error("Expected request to fail when health check is disabled")
	}
}

package admission_controller

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/civts/markhor/pkg"
	apiV1 "github.com/civts/markhor/pkg/api/types/v1"
	"github.com/civts/markhor/pkg/decrypt"
	admissionV1 "k8s.io/api/admission/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func validateHandler(w http.ResponseWriter, r *http.Request) {
	const prefix = "[validate admission hook]"
	var admissionReview admissionV1.AdmissionReview
	err := json.NewDecoder(r.Body).Decode(&admissionReview)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not decode request: %v", err), http.StatusBadRequest)
		slog.Warn(fmt.Sprintf("%s could not decode request: %v", prefix, err))
		return
	}
	eventId := slog.String(pkg.SLOG_EVENT_ID_KEY, string(admissionReview.Request.UID))
	slog.Debug(prefix+" received a new request", eventId)

	var admissionResponse admissionV1.AdmissionResponse
	name, err := validateReview(&admissionReview, eventId)
	if err == nil {
		slog.Info(fmt.Sprint(prefix, " successfully validated the MarkhorSecret ", name), eventId)
		admissionResponse = admissionV1.AdmissionResponse{Allowed: true, UID: admissionReview.Request.UID}
	} else {
		if strings.Contains(err.Error(), "expected mac") {
			slog.Debug(fmt.Sprintf("%s the request for %s was not valid: invalid MAC", prefix, name), eventId)
		} else {
			slog.Debug(fmt.Sprintf("%s the request for %s was not valid: %v", prefix, name, err), eventId)
		}
		admissionResponse = admissionV1.AdmissionResponse{
			Result: &metaV1.Status{
				Status:  metaV1.StatusFailure,
				Message: err.Error(),
			},
			Allowed: false,
			UID:     admissionReview.Request.UID,
		}
	}

	response := admissionV1.AdmissionReview{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "admission.k8s.io/v1",
			Kind:       "AdmissionReview",
		},
		Response: &admissionResponse,
	}

	respBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
		slog.Error(fmt.Sprintf("%s could not encode response: %v", prefix, err), eventId)
		return
	}

	_, err = w.Write(respBytes)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
		slog.Error(fmt.Sprintf("%s could not write response: %v", prefix, err), eventId)
	}
}

func validateReview(review *admissionV1.AdmissionReview, eventId slog.Attr) (string, error) {
	var err error

	// Check that the object being validated is a markhor markhorSecret
	var markhorSecret apiV1.MarkhorSecret
	err2 := json.Unmarshal(review.Request.Object.Raw, &markhorSecret)

	if err2 != nil {
		err = field.Invalid(field.NewPath("kind"), "", "The object being validated must be a markhorsecret")
	}
	if err != nil {
		return "", err
	}

	msName := fmt.Sprint(markhorSecret.Namespace, "/", markhorSecret.Name)
	// Check that we can successfully decrypt the data in the MarkhorSecret
	_, err = decrypt.DecryptMarkhorSecretEvent(&markhorSecret, eventId)

	return msName, err
}

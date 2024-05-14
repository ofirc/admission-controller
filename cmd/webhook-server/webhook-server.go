/*
Copyright (c) 2024 Ofir Cohen.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	tlsDir      = `/run/secrets/tls`
	tlsCertFile = `tls.crt`
	tlsKeyFile  = `tls.key`
)

var debugLoggingEnabled bool

func init() {
	debugLoggingEnabled = os.Getenv("LOG_LEVEL") == "debug"
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from root endpoint!\n")
}

func handleAdmissionReview(w http.ResponseWriter, r *http.Request) {
	var review v1.AdmissionReview
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, fmt.Sprintf("could not decode request: %v", err), http.StatusBadRequest)
		return
	}

	if debugLoggingEnabled {
		reviewBytes, err := json.MarshalIndent(review, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling review: %v\n", err)
		} else {
			fmt.Printf("Received AdmissionReview: %s\n", string(reviewBytes))
		}
	}

	response := v1.AdmissionResponse{
		UID: review.Request.UID,
	}

	// Decode the Pod object from the AdmissionRequest
	var pod corev1.Pod
	if err := json.Unmarshal(review.Request.Object.Raw, &pod); err != nil {
		response.Allowed = false
		response.Result = &metav1.Status{
			Message: fmt.Sprintf("could not decode pod object: %v", err),
		}
	} else {
		// Set Allowed to false if the Pod name is "pod-with-an-invalid-name"
		if pod.Name == "pod-with-an-invalid-name" {
			response.Allowed = false
			response.Result = &metav1.Status{
				Message: "Pod with an invalid name is not allowed",
			}
		} else {
			response.Allowed = true
		}
	}

	// Respond with the same AdmissionReview structure
	review.Response = &response

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(review); err != nil {
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
	}
}

func main() {
	certPath := filepath.Join(tlsDir, tlsCertFile)
	keyPath := filepath.Join(tlsDir, tlsKeyFile)

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/validate", handleAdmissionReview)
	fmt.Println("Starting validating admission webhook server...")
	if err := http.ListenAndServeTLS(":8443", certPath, keyPath, nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

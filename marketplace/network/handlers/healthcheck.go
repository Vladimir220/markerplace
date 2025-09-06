package handlers

import (
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func HealthListener() {
	http.HandleFunc("/health", HealthHandler)
	http.ListenAndServe(":8888", nil)
}

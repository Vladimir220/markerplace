package main

import (
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func HealthListener() {
	http.HandleFunc("/health", healthHandler)
	http.ListenAndServe(":8888", nil)
}

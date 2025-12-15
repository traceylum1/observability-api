package handlers

import (
	"net/http"
)

func Live(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func Ready(w http.ResponseWriter, r *http.Request) {
	if err := checkDependencies(r); err != nil {
		http.Error(w, "not ready", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}

func checkDependencies(r *http.Request) error {
	// Example:
	// return db.Ping(r.Context())

	return nil
}
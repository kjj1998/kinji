package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		slog.Error(err.Error())
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func requireUserID(w http.ResponseWriter, r *http.Request) (string, bool) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "User ID not provided")
		return "", false
	}
	return id, true
}

func parseDate(w http.ResponseWriter, value, param string) (*time.Time, bool) {
	if value == "" {
		return nil, true
	}
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid "+param+": expected YYYY-MM-DD")
		return nil, false
	}

	return &t, true
}

package http

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		slog.Error(err.Error())
	}
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

func RequireUserId(w http.ResponseWriter, r *http.Request) (string, bool) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "User ID not provided")
		return "", false
	}
	return id, true
}

// ParseMonthYear validates the month and year query params, returning them
// normalized as MM ("01"-"12") and YYYY. If either param is empty, it
// defaults to the current month and year.
func ParseMonthYear(w http.ResponseWriter, monthVal, yearVal string) (month, year string, ok bool) {
	if monthVal == "" || yearVal == "" {
		now := time.Now()
		return now.Format("01"), now.Format("2006"), true
	}

	m, err := strconv.Atoi(monthVal)
	if err != nil || m < 1 || m > 12 {
		slog.Error("invalid month query param", "value", monthVal)
		WriteError(w, http.StatusBadRequest, "invalid month: expected 1-12")
		return "", "", false
	}

	y, err := strconv.Atoi(yearVal)
	if err != nil || y < 1000 || y > 9999 {
		slog.Error("invalid year query param", "value", yearVal)
		WriteError(w, http.StatusBadRequest, "invalid year: expected YYYY")
		return "", "", false
	}

	return fmt.Sprintf("%02d", m), fmt.Sprintf("%d", y), true
}

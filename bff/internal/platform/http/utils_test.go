package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	WriteJSON(w, http.StatusCreated, map[string]int{"a": 1})

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want 201", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	var out map[string]int
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out["a"] != 1 {
		t.Errorf("body = %+v", out)
	}
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	WriteError(w, http.StatusBadRequest, "bad thing")

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
	var out map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out["error"] != "bad thing" {
		t.Errorf("error body = %+v", out)
	}
}

func TestRequireUserId(t *testing.T) {
	t.Run("present", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/x", nil)
		r.SetPathValue("id", "u1")
		w := httptest.NewRecorder()

		id, ok := RequireUserId(w, r)
		if !ok || id != "u1" {
			t.Errorf("got (%q, %v), want (u1, true)", id, ok)
		}
	})

	t.Run("missing writes 400", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/x", nil)
		w := httptest.NewRecorder()

		id, ok := RequireUserId(w, r)
		if ok || id != "" {
			t.Errorf("got (%q, %v), want (\"\", false)", id, ok)
		}
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})
}

func TestParseMonthYear(t *testing.T) {
	cases := []struct {
		name              string
		month, year       string
		wantMonth         string
		wantYear          string
		wantOK            bool
	}{
		{"valid single digit month padded", "6", "2026", "06", "2026", true},
		{"valid two digit month", "12", "2026", "12", "2026", true},
		{"month lower bound", "1", "2026", "01", "2026", true},
		{"month too low", "0", "2026", "", "", false},
		{"month too high", "13", "2026", "", "", false},
		{"month not a number", "abc", "2026", "", "", false},
		{"year too low", "6", "999", "", "", false},
		{"year too high", "6", "10000", "", "", false},
		{"year not a number", "6", "abc", "", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			month, year, ok := ParseMonthYear(w, tc.month, tc.year)
			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOK)
			}
			if ok {
				if month != tc.wantMonth || year != tc.wantYear {
					t.Errorf("got (%q, %q), want (%q, %q)", month, year, tc.wantMonth, tc.wantYear)
				}
			} else if w.Code != http.StatusBadRequest {
				t.Errorf("expected 400 on invalid input, got %d", w.Code)
			}
		})
	}
}

func TestParseMonthYear_EmptyDefaultsToNow(t *testing.T) {
	now := time.Now()
	w := httptest.NewRecorder()

	month, year, ok := ParseMonthYear(w, "", "")
	if !ok {
		t.Fatal("expected ok for empty input")
	}
	if month != now.Format("01") || year != now.Format("2006") {
		t.Errorf("got (%q, %q), want current (%q, %q)", month, year, now.Format("01"), now.Format("2006"))
	}
	// sanity: defaulted month is a valid 1-12 value
	if m, _ := strconv.Atoi(month); m < 1 || m > 12 {
		t.Errorf("defaulted month %q out of range", month)
	}
}

func TestHealth(t *testing.T) {
	w := httptest.NewRecorder()
	Health(w, httptest.NewRequest(http.MethodGet, "/health", nil))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"status":"ok"`) {
		t.Errorf("body = %q", w.Body.String())
	}
}

package httpapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	NewRouter().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}

	gotType := response.Header().Get("Content-Type")
	if gotType != "application/json" {
		t.Fatalf("Content-Type = %q, want %q", gotType, "application/json")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	var payload map[string]string
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("unmarshal %s: %v", body, err)
	}
	if len(payload) != 1 || payload["status"] != "ok" {
		t.Fatalf("body = %#v, want map[status:ok]", payload)
	}
}

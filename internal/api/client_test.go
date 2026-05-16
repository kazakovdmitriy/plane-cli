package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClientDoSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "test-token" {
			w.WriteHeader(401)
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"id":"1","name":"Test"}`))
	}))
	defer server.Close()

	cfg := DefaultRetryConfig()
	client := NewClient("test-token", server.URL, "ws", "prj", cfg)
	body, code, err := client.Do("GET", EndpointWorkItem("1"), nil, nil)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if code != 200 {
		t.Errorf("Expected 200, got %d", code)
	}
	if string(body) != `{"id":"1","name":"Test"}` {
		t.Errorf("Got body: %s", body)
	}
}

func TestClient401NoRetry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
	}))
	defer server.Close()

	cfg := DefaultRetryConfig()
	cfg.TotalTimeout = 2 * time.Second
	client := NewClient("bad-token", server.URL, "ws", "prj", cfg)
	_, code, err := client.Do("GET", EndpointWorkItems(), nil, nil)
	if err == nil {
		t.Errorf("Expected error for 401")
	}
	if code != 401 {
		t.Errorf("Expected 401, got %d", code)
	}
}

func TestClientURLConstruction(t *testing.T) {
	var capturedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.WriteHeader(200)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	cfg := DefaultRetryConfig()
	client := NewClient("tk", server.URL, "myorg", "abc", cfg)
	client.Do("GET", EndpointWorkItem("test-1"), nil, nil)

	expected := "/api/v1/workspaces/myorg/projects/abc/work-items/test-1/"
	if capturedPath != expected {
		t.Errorf("Expected path %s, got %s", expected, capturedPath)
	}
}

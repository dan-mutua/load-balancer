package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func createBackend(urlStr string) *Backend {
	backendURL, _ := url.Parse(urlStr)
	return &Backend{
		URL:          backendURL,
		Alive:        true,
		ReverseProxy: httputil.NewSingleHostReverseProxy(backendURL),
	}
}

func TestSetAlive(t *testing.T) {
	backend := createBackend("http://localhost:8080")
	backend.SetAlive(false)

	if backend.IsAlive() {
		t.Errorf("Expected backend to be down, got up")
	}

	backend.SetAlive(true)

	if !backend.IsAlive() {
		t.Errorf("Expected backend to be up, got down")
	}
}

func TestServerPool_AddBackend(t *testing.T) {
	sp := &ServerPool{}
	backend := createBackend("http://localhost:8080")
	sp.AddBackend(backend)

	if len(sp.backends) != 1 {
		t.Errorf("Expected 1 backend, got %d", len(sp.backends))
	}
}

func TestServerPool_GetNextPeer(t *testing.T) {
	sp := &ServerPool{}
	backend1 := createBackend("http://localhost:8080")
	backend2 := createBackend("http://localhost:8081")

	sp.AddBackend(backend1)
	sp.AddBackend(backend2)

	peer := sp.GetNextPeer()
	if peer == nil {
		t.Errorf("Expected to get a peer, got nil")
	}
}

func TestHealthCheck(t *testing.T) {
	sp := &ServerPool{}
	backend := createBackend("http://localhost:8080")
	sp.AddBackend(backend)

	// Assume the backend is down for this test
	backend.SetAlive(false)

	sp.HealthCheck()

	if backend.IsAlive() {
		t.Errorf("Expected backend to be down after health check, got up")
	}
}

func TestLbHandler(t *testing.T) {
	backend1 := createBackend("http://localhost:8080")
	backend2 := createBackend("http://localhost:8081")

	serverPool = ServerPool{}
	serverPool.AddBackend(backend1)
	serverPool.AddBackend(backend2)

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	lb(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("Expected status 503, got %d", resp.StatusCode)
	}
}

func TestIsBackendAlive(t *testing.T) {
	backend := createBackend("http://localhost:8080")

	if isBackendAlive(backend.URL) {
		t.Errorf("Expected backend to be down, got up")
	}

	// To properly test an alive backend, a real server should be running at the backend URL
	// Uncomment the next lines if you have a real server running
	// backend.SetAlive(true)
	// if !isBackendAlive(backend.URL) {
	// 	t.Errorf("Expected backend to be up, got down")
	// }
}

func TestMain(m *testing.M) {
	// Create a test backend server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	// Setup the server pool with the test server
	serverPool = ServerPool{}
	testBackend := createBackend(testServer.URL)
	serverPool.AddBackend(testBackend)

	// Run the tests
	m.Run()
}

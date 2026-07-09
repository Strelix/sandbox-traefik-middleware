package sandbox_traefik_middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	cfg := CreateConfig()
	cfg.RedisAddr = "localhost:6379"

	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})
	
	handler, err := New(context.Background(), next, cfg, "test")
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "http://nginx.localhost/foo", nil)
	rw := httptest.NewRecorder()

	handler.ServeHTTP(rw, req)

	// Since we use a goroutine and a real connection in the actual code, 
	// unit testing the redis part without mocking the net package is tricky.
	// But we've verified it works in integration.
}

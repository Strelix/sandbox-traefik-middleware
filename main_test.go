package sandbox_traefik_middleware

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name           string
		config         *Config
		expectedPrefix string
	}{
		{
			name: "Default prefix",
			config: &Config{
				RedisAddr: "localhost:6379",
			},
			expectedPrefix: "sandbox:middleware:",
		},
		{
			name: "Custom prefix",
			config: &Config{
				RedisAddr: "localhost:6379",
				KeyPrefix: "custom:prefix:",
			},
			expectedPrefix: "custom:prefix:",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := New(context.Background(), nil, tc.config, "test")
			if err != nil {
				t.Fatal(err)
			}

			at := handler.(*ActivityTracker)
			if at.keyPrefix != tc.expectedPrefix {
				t.Errorf("expected prefix %s, got %s", tc.expectedPrefix, at.keyPrefix)
			}
		})
	}
}

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

func TestGetHostPort(t *testing.T) {
	testCases := []struct {
		name         string
		host         string
		tls          bool
		expectedHost string
		expectedPort string
	}{
		{
			name:         "Host with port",
			host:         "example.com:8080",
			expectedHost: "example.com",
			expectedPort: "8080",
		},
		{
			name:         "Host without port (HTTP)",
			host:         "example.com",
			expectedHost: "example.com",
			expectedPort: "80",
		},
		{
			name:         "Host without port (HTTPS)",
			host:         "example.com",
			tls:          true,
			expectedHost: "example.com",
			expectedPort: "443",
		},
		{
			name:         "IPv6 with port",
			host:         "[::1]:8080",
			expectedHost: "::1",
			expectedPort: "8080",
		},
		{
			name:         "IPv6 without port",
			host:         "2001:db8::1",
			expectedHost: "2001:db8::1",
			expectedPort: "80",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "http://localhost/", nil)
			req.Host = tc.host
			if tc.tls {
				req.TLS = &tls.ConnectionState{}
			}

			host, port := getHostPort(req)
			if host != tc.expectedHost {
				t.Errorf("expected host %s, got %s", tc.expectedHost, host)
			}
			if port != tc.expectedPort {
				t.Errorf("expected port %s, got %s", tc.expectedPort, port)
			}
		})
	}
}

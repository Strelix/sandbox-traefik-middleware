package sandbox_traefik_middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Config struct {
	RedisAddr     string `json:"redisAddr,omitempty"`
	RedisPassword string `json:"redisPassword,omitempty"`
	RedisUser     string `json:"redisUser,omitempty"`
}

func CreateConfig() *Config {
	return &Config{}
}

type ActivityTracker struct {
	next          http.Handler
	name          string
	redisAddr     string
	redisPassword string
	redisUser     string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.RedisAddr == "" {
		return nil, fmt.Errorf("redisAddr cannot be empty")
	}

	return &ActivityTracker{
		next:          next,
		name:          name,
		redisAddr:     config.RedisAddr,
		redisPassword: config.RedisPassword,
		redisUser:     config.RedisUser,
	}, nil
}

func (m *ActivityTracker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	host, port := getHostPort(req)
	key := fmt.Sprintf("sandbox:middleware:%s:%s", host, port)
	val := fmt.Sprintf("%d", time.Now().Unix())

	go func() {
		conn, err := net.DialTimeout("tcp", m.redisAddr, 2*time.Second)
		if err != nil {
			return
		}
		defer conn.Close()

		if m.redisPassword != "" {
			var authCmd string
			if m.redisUser != "" {
				// AUTH <user> <password>
				authCmd = fmt.Sprintf("*3\r\n$4\r\nAUTH\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(m.redisUser), m.redisUser, len(m.redisPassword), m.redisPassword)
			} else {
				// AUTH <password>
				authCmd = fmt.Sprintf("*2\r\n$4\r\nAUTH\r\n$%d\r\n%s\r\n", len(m.redisPassword), m.redisPassword)
			}
			_, _ = conn.Write([]byte(authCmd))
		}

		// Minimal Redis SET using RESP protocol: SET <key> <val> EX 3600
		// Format: *5\r\n$3\r\nSET\r\n$<len_key>\r\n<key>\r\n$<len_val>\r\n<val>\r\n$2\r\nEX\r\n$4\r\n3600\r\n
		cmd := fmt.Sprintf("*5\r\n$3\r\nSET\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n$2\r\nEX\r\n$4\r\n3600\r\n", len(key), key, len(val), val)
		_, _ = conn.Write([]byte(cmd))
	}()

	m.next.ServeHTTP(rw, req)
}

func getHostPort(req *http.Request) (string, string) {
	host, port, err := net.SplitHostPort(req.Host)
	if err != nil {
		// If SplitHostPort fails, it likely means the port is missing
		host = req.Host
		if req.TLS != nil {
			port = "443"
		} else {
			port = "80"
		}
	}
	return host, port
}

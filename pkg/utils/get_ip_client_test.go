package utils

import (
	"net/http"
	"testing"
)

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name       string
		headers    map[string]string
		remoteAddr string
		want       string
	}{
		{
			name: "Prioritize X-Forwarded-For",
			headers: map[string]string{
				"X-Forwarded-For": "203.0.113.195",
				"X-Real-IP":       "198.51.100.1",
			},
			remoteAddr: "192.0.2.1:12345",
			want:       "203.0.113.195",
		},
		{
			name: "Multiple IPs in X-Forwarded-For",
			headers: map[string]string{
				"X-Forwarded-For": "203.0.113.195, 10.0.0.1",
			},
			remoteAddr: "192.0.2.1:12345",
			want:       "203.0.113.195, 10.0.0.1",
		},
		{
			name: "Fallback to X-Real-IP if X-Forwarded-For is missing",
			headers: map[string]string{
				"X-Real-IP": "198.51.100.1",
			},
			remoteAddr: "192.0.2.1:12345",
			want:       "198.51.100.1",
		},
		{
			name:       "Fallback to RemoteAddr (IPv4 with port)",
			headers:    nil,
			remoteAddr: "192.0.2.1:12345",
			want:       "192.0.2.1",
		},
		{
			name:       "Fallback to RemoteAddr (IPv6 with port)",
			headers:    nil,
			remoteAddr: "[2001:db8::1]:8080",
			want:       "2001:db8::1",
		},
		{
			name:       "Fallback to RemoteAddr (Invalid format or missing port)",
			headers:    nil,
			remoteAddr: "192.0.2.1",
			want:       "192.0.2.1",
		},
		{
			name: "Empty headers but keys exist",
			headers: map[string]string{
				"X-Forwarded-For": "",
				"X-Real-IP":       "",
			},
			remoteAddr: "127.0.0.1:80",
			want:       "127.0.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Header:     make(http.Header),
				RemoteAddr: tt.remoteAddr,
			}

			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			got := GetClientIP(req)
			if got != tt.want {
				t.Errorf("GetClientIP() = %q, want %q", got, tt.want)
			}
		})
	}
}
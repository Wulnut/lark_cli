package authtest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"lark_cli/internal/auth"
)

func TestPluginTokenClient_Fetch(t *testing.T) {
	tests := []struct {
		name           string
		pluginID       string
		pluginSecret   string
		serverStatus   int
		serverResponse string
		wantErr        bool
		errContains    string
		wantToken      string
		wantExpires    time.Duration
	}{
		{
			name:           "success",
			pluginID:       "pid",
			pluginSecret:   "sec",
			serverStatus:   http.StatusOK,
			serverResponse: `{"data":{"expire_time":7200,"token":"p-1234"},"error":{"code":0,"msg":"success"}}`,
			wantToken:      "p-1234",
			wantExpires:    7200 * time.Second,
		},
		{
			name:           "api error",
			pluginID:       "pid",
			pluginSecret:   "sec",
			serverStatus:   http.StatusOK,
			serverResponse: `{"data":{},"error":{"code":10001,"msg":"invalid secret"}}`,
			wantErr:        true,
			errContains:    "API error 10001: invalid secret",
		},
		{
			name:           "http error",
			pluginID:       "pid",
			pluginSecret:   "sec",
			serverStatus:   http.StatusInternalServerError,
			serverResponse: `Internal Server Error`,
			wantErr:        true,
			errContains:    "unexpected status code: 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected method POST, got %s", r.Method)
				}
				if r.URL.Path != "/open_api/authen/plugin_token" {
					t.Errorf("expected path /open_api/authen/plugin_token, got %s", r.URL.Path)
				}
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
				}

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("failed to read body: %v", err)
				}

				var req map[string]interface{}
				if err := json.Unmarshal(body, &req); err != nil {
					t.Fatalf("failed to unmarshal body: %v", err)
				}

				if req["plugin_id"] != tt.pluginID {
					t.Errorf("expected plugin_id %s, got %v", tt.pluginID, req["plugin_id"])
				}
				if req["plugin_secret"] != tt.pluginSecret {
					t.Errorf("expected plugin_secret %s, got %v", tt.pluginSecret, req["plugin_secret"])
				}

				typeVal, ok := req["type"].(float64)
				if !ok || typeVal != 0 {
					t.Errorf("expected type 0, got %v", req["type"])
				}

				w.WriteHeader(tt.serverStatus)
				_, _ = w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			client := auth.NewPluginTokenClient(server.Client(), server.URL, tt.pluginID, tt.pluginSecret)

			token, expires, err := client.Fetch(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got none")
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if token != tt.wantToken {
				t.Errorf("expected token %q, got %q", tt.wantToken, token)
			}
			if expires != tt.wantExpires {
				t.Errorf("expected expires %v, got %v", tt.wantExpires, expires)
			}
		})
	}
}

func TestPluginTokenClient_URLSanitization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/open_api/authen/plugin_token" {
			t.Errorf("expected path /open_api/authen/plugin_token, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"expire_time":7200,"token":"p-1234"},"error":{"code":0,"msg":"success"}}`))
	}))
	defer server.Close()

	// Try with trailing slash
	client := auth.NewPluginTokenClient(server.Client(), server.URL+"/", "pid", "sec")
	_, _, err := client.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

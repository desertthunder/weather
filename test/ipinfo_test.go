package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/desertthunder/weather/internal/ipinfo"
)

type ipinfoTest struct {
	name     string
	token    string
	ipaddr   string
	mockResp string
	wantErr  bool
}

func TestGeolocate(t *testing.T) {
	tests := []ipinfoTest{
		{
			name:     "Valid IP and token",
			token:    "valid_token",
			ipaddr:   "8.8.8.8",
			mockResp: `{"city": "Austin", "region": "Texas", "country": "US", "loc": "30.2672,-97.7431"}`,
			wantErr:  false,
		},
		{
			name:    "Empty token",
			token:   "",
			ipaddr:  "8.8.8.8",
			wantErr: true, // We expect an error because the token is empty
		},
		{
			name:     "Invalid IP",
			token:    "valid_token",
			ipaddr:   "invalid_ip",
			mockResp: `{"error": "invalid IP address"}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// The testing strategy here is to use a mock server to return a response
			// that we've defined in the test case. Then we override the base URL of the
			// client to point to the mock server. This way, we can test the client's
			// behavior without making actual requests to the IPInfo API.
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(tt.mockResp))
			}))

			defer server.Close()

			client := &ipinfo.IPInfoClient{Token: tt.token}

			client.BaseURL = server.URL

			_, err := client.Geolocate(&tt.ipaddr)

			if err == nil && tt.wantErr {
				t.Errorf("Geolocate() got no error, want error")
			}

			if err != nil && !tt.wantErr {
				t.Errorf("Geolocate() got error, want no error")
			}
		})
	}
}

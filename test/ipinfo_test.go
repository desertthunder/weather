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

func TestIPInfoClient(t *testing.T) {
	t.Run("Setters", func(t *testing.T) {
		client := &ipinfo.IPInfoClient{Token: "valid_token"}

		client.SetToken("new_token")
		client.SetURL("https://new_url.com")

		if client.Token != "new_token" {
			t.Errorf("Expected token to be new_token, got %s", client.Token)
		}

		if client.BaseURL != "https://new_url.com" {
			t.Errorf("Expected base URL to be https://new_url.com, got %s", client.BaseURL)
		}

	})

	t.Run("Geolocate", func(t *testing.T) {
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
	})

	t.Run("BuildCity", func(t *testing.T) {
		r := ipinfo.IPInfoResponse{
			City:     "Austin",
			Region:   "Texas",
			Country:  "US",
			Location: "30.2672,-97.7431",
		}

		city := r.BuildCity()

		if city.Name != "Austin" {
			t.Errorf("Expected city name to be Austin, got %s", city.Name)
		}

		if city.Lat != 30.2672 {
			t.Errorf("Expected latitude to be 30.2672, got %f", city.Lat)
		}

		if city.Long != -97.7431 {
			t.Errorf("Expected longitude to be -97.7431, got %f", city.Long)
		}
	})

	t.Run("Point", func(t *testing.T) {
		r := ipinfo.IPInfoResponse{
			Location: "30.2672,-97.7431",
		}

		lat, lon := r.Point()

		if lat != 30.2672 {
			t.Errorf("Expected latitude to be 30.2672, got %f", lat)
		}

		if lon != -97.7431 {
			t.Errorf("Expected longitude to be -97.7431, got %f", lon)
		}
	})

	t.Run("Validate", func(t *testing.T) {
		r := ipinfo.IPInfoResponse{
			Location: "30.2672,-97.7431",
		}

		err := r.Validate([]byte(`{"city": "Austin", "region": "Texas", "country": "US", "loc": "30.2672,-97.7431"}`))

		if err != nil {
			t.Errorf("Expected error to be nil, got %s", err.Error())
		}

		err = r.Validate([]byte(`{"bogon": "true"}`))

		if err == nil {
			t.Errorf("Expected error to be non-nil, got %s", err)
		}
	})
}

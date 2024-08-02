package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/desertthunder/weather/cmd/cli"
	"github.com/desertthunder/weather/internal/ipinfo"
	"github.com/spf13/viper"
)

type rootInputs struct {
	ip string
}

type rootActionTest struct {
	name     string
	inputs   rootInputs
	expected string
	error    bool
}

func TestRootAction(t *testing.T) {
	// Set up environment variable
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"city": "Austin", "region": "Texas", "country": "US", "loc": "30.2672,-97.7431"}`))
	}))

	tests := []rootActionTest{
		{
			name: "Valid IP",
			inputs: rootInputs{
				ip: "127.0.0.1",
			},
			expected: "City",
			error:    false,
		},

		{
			name:     "Invalid IP",
			inputs:   rootInputs{ip: "invalid_ip"},
			expected: "Invalid IP address",
			error:    true,
		},

		{
			name:     "Empty IP",
			inputs:   rootInputs{ip: ""},
			expected: "",
			error:    false,
		},
	}

	viper.Set("IPINFO_TOKEN", "test_token")
	ipc := ipinfo.NewIPInfoClient(viper.GetString("IPINFO_TOKEN"))
	ipc.SetURL(server.URL)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Call the function
			cli.RootAction(tt.inputs.ip, ipc)

			// Restore stdout and close the pipe
			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			buf.ReadFrom(r)

			if tt.error && !bytes.Contains(buf.Bytes(), []byte(tt.expected)) {
				t.Errorf("Expected error %s not found in output %s", tt.expected, buf.String())
			}

			if !bytes.Contains(buf.Bytes(), []byte(tt.expected)) {
				t.Errorf("Expected value %s not found in output %s", tt.expected, buf.String())
			}
		})
	}
}

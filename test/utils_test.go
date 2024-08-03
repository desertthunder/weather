package test

import (
	"strings"
	"testing"

	"github.com/desertthunder/weather/internal/utils"
)

type tStruct struct {
	Some string
	Data string
	Of   string
	Any  int
}

func TestUtils(t *testing.T) {
	t.Run("IP Address Helpers", func(t *testing.T) {
		t.Run("Validate", func(t *testing.T) {
			t.Run("Valid", func(t *testing.T) {
				valid_ip := "8.8.8.8"
				valid := utils.ValidateIPAddress(valid_ip)

				if !valid {
					t.Errorf("Expected %s to be valid, got %t", valid_ip, valid)
				}
			})

			t.Run("Invalid", func(t *testing.T) {
				invalid_ip := "invalid"
				valid := utils.ValidateIPAddress(invalid_ip)

				if valid {
					t.Errorf("Expected %s to be invalid, got %t", invalid_ip, valid)
				}
			})
		})
	})

	t.Run("JSON Helpers", func(t *testing.T) {
		data := []byte(`{"key": "value"}`)

		t.Run("GetRawJSON", func(t *testing.T) {
			got := utils.GetRawJSON(data)

			if !strings.Contains(got, "key") {
				t.Errorf("Expected key not found in output")
			}

			if !strings.Contains(got, "value") {
				t.Errorf("Expected value not found in output")
			}
		})

		t.Run("PrintJSON", func(t *testing.T) {
			s := tStruct{
				Some: "some",
				Data: "data",
				Of:   "of",
				Any:  1,
			}

			output := CaptureOutput(func() {
				utils.PrintJSON(s)
			})

			want := [][]string{
				{"Some", "some"},
				{"Data", "data"},
				{"Of", "of"},
				{"Any", "1"},
			}

			if output == "" {
				t.Fatalf("Expected output to be non-empty, got %s", output)
			}

			for _, w := range want {
				if !strings.Contains(output, w[0]) {
					t.Logf("Expected %s not found in output %s", w[0], output)
					t.Errorf("Expected %s not found in output", w[0])
				}
				if !strings.Contains(output, w[1]) {
					t.Logf("Expected %s not found in output %s", w[1], output)
					t.Errorf("Expected %s not found in output", w[1])
				}
			}
		})

		t.Run("PrintRawJSON", func(t *testing.T) {
			data := []byte(`{"key": "value"}`)

			output := CaptureOutput(func() {
				utils.PrintRawJSON(data)
			})

			if !strings.Contains(output, "key") {
				t.Errorf("Expected key not found in output")
			}

			if !strings.Contains(output, "value") {
				t.Errorf("Expected value not found in output")
			}

		})
	})
}

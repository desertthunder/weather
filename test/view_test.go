package test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/desertthunder/weather/internal/view"
)

type test struct {
	name string
	data [][]string
}

func TestTable(t *testing.T) {
	headers := []string{"City", "Latitude", "Longitude"}
	tests := []test{
		{
			name: "Single row",
			data: [][]string{{"Seattle", "47.6062", "-122.3321"}},
		},
		{
			name: "Multiple rows",
			data: [][]string{
				{"Seattle", "47.6062", "-122.3321"},
				{"Portland", "45.5152", "-122.6784"},
			},
		},
		{
			name: "No data",
			data: [][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Call the function
			tbl := view.Table(headers, tt.data)
			fmt.Println(tbl.Render())

			// Restore stdout and close the pipe
			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			buf.ReadFrom(r)

			for _, header := range headers {
				if !bytes.Contains(buf.Bytes(), []byte(header)) {
					t.Errorf("Expected header %s not found in output", header)
				}
			}

			for _, row := range tt.data {
				for _, cell := range row {
					if !bytes.Contains(buf.Bytes(), []byte(cell)) {
						t.Errorf("Expected cell %s not found in output", cell)
					}
				}
			}
		})
	}
}

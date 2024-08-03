package test

import (
	"bytes"
	"os"
)

func CaptureOutput(f func()) string {
	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()

	os.Stdout = w

	f()

	w.Close()

	os.Stdout = old

	buf.ReadFrom(r)

	return buf.String()
}

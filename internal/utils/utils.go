// Package utils contains utility functions used
// throughout the project.
package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// SaveToFile takes a file path of the form
// out/filename.json and saves the data to the file.
func SaveToFile(f string, d []byte) error {
	parts := strings.Split(f, "/")
	filename := parts[0]

	if len(parts) < 1 {
		return errors.New("invalid file path")
	}

	if len(parts) > 2 {
		filename = parts[len(parts)-1]
	}

	if _, err := os.Stat("out"); os.IsNotExist(err) {
		os.Mkdir("out", 0755)

		if err != nil {
			return err
		}
	}

	file, err := os.Create(filename)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(d))

	return err
}

func PrintJSON(data interface{}) {
	s, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(s))
}

func ValidateIPAddress(ipaddr string) bool {
	parts := strings.Split(ipaddr, ".")

	if len(parts) != 4 {
		return false
	}

	for _, p := range parts {
		if len(p) < 1 || len(p) > 3 {
			return false
		}
	}

	return true
}

func PrintRawJSON(data []byte) {
	dst := &bytes.Buffer{}

	json.Indent(dst, data, "", "  ")

	fmt.Println(dst.String())
}

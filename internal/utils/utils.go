// Package utils contains utility functions used throughout the project.
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

// SaveToFile takes a file path of the form out/filename.json and saves the data
// to the file (f is the file path).
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

// PrintJSON takes any data and prints it as a JSON string. The purpose of this
// is to take a struct to pretty print it as JSON for debugging purposes.
func PrintJSON(data interface{}) {
	s, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(s))
}

// func ValidateIPAddress validates an IP address string.
//
// It checks if the IP address string is in the correct format and has four
// parts (e.g., 192.168.1.1). Each part is checked via validateIPParts.
func ValidateIPAddress(ipaddr string) bool {
	parts := strings.Split(ipaddr, ".")

	return len(parts) == 4 && validateIPParts(parts)
}

func validateIPParts(parts []string) bool {
	for _, p := range parts[0:3] {
		if len(p) < 1 || len(p) > 3 || p == "0" {
			return false
		}
	}

	return true
}

// PrintRawJSON takes a byte array and prints it as a raw JSON string. The
// purpose of this is to take a byte array pulled from a response body and
// pretty print it as JSON for debugging purposes.
func PrintRawJSON(data []byte) {
	s := GetRawJSON(data)

	fmt.Println(s)
}

// GetRawJSON takes a byte array and returns it as a raw JSON string.
func GetRawJSON(data []byte) string {
	dst := &bytes.Buffer{}

	json.Indent(dst, data, "", "  ")

	return dst.String()
}

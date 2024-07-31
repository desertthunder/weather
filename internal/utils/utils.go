// Package utils contains utility functions used
// throughout the project.
package utils

import (
	"bytes"
	"errors"
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

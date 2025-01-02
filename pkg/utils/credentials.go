package utils

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func ReadCredFile(filename string) (string, error) {
	path := filepath.Join("creds", filename)
	content, err := ioutil.ReadFile(path)
	if err != nil {
			return "", fmt.Errorf("error reading %s credential file: %v", filename, err)
	}
	return string(content), nil
}
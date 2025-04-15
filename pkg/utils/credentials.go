package utils

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func ReadCredFile(filename string) (string, error) {
	path := filepath.Join("creds", filename)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
			return "", fmt.Errorf("error reading %s credential file: %v", filename, err)
	}
	content := string(bytes)
	content = strings.TrimSpace(content)
	return string(content), nil
}
package utils

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
)

func ConvertReaderToString(body io.ReadCloser) (string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(body)

	if err != nil {
		return "", err
	}

	return buf.String(), err
}

func CreateRegExpForNamePrefix(prefix string) *regexp.Regexp {
	exp := fmt.Sprintf("^%s", prefix)
	return regexp.MustCompile(exp)
}

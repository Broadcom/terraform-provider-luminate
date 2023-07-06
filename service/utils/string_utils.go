package utils

import (
	"bytes"
	"io"
)

func ConvertReaderToString(body io.ReadCloser) (string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(body)

	if err != nil {
		return "", err
	}

	return buf.String(), err
}

package utils

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExtractIPAndPort(t *testing.T) {
	inputIP := "1.0.0.1"
	ip, port := ExtractIPAndPort(inputIP)

	require.Equal(t, inputIP, ip)
	require.Empty(t, port)

	input := fmt.Sprintf("bla://%s:126", inputIP)
	ip, port = ExtractIPAndPort(input)

	require.Equal(t, inputIP, ip)
	require.Equal(t, "126", port)

	input = ""
	ip, port = ExtractIPAndPort(input)

	require.Empty(t, ip)
	require.Empty(t, port)
}

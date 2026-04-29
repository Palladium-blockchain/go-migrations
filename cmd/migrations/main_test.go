package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun_PrintsVersionWithShortFlag(t *testing.T) {
	var out bytes.Buffer

	exitCode := run([]string{"migration", "-v"}, &out)

	require.Equal(t, 0, exitCode)
	require.NotEmpty(t, out.String())
}

func TestRun_PrintsVersionWithLongFlag(t *testing.T) {
	var out bytes.Buffer

	exitCode := run([]string{"migration", "--version"}, &out)

	require.Equal(t, 0, exitCode)
	require.NotEmpty(t, out.String())
}

func TestCurrentVersion_IsNotEmpty(t *testing.T) {
	require.NotEmpty(t, currentVersion())
}

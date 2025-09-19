package main

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionCommand(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "version")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)

	outputStr := string(output)
	assert.Contains(t, outputStr, "sz version 0.1.0")
}

func TestHelpCommand(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "help")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)

	outputStr := string(output)
	assert.Contains(t, outputStr, "sz is a CLI web browser")
	assert.Contains(t, outputStr, "Available Commands")
}

func TestHelpByDefault(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)

	outputStr := string(output)
	assert.Contains(t, outputStr, "sz is a CLI web browser")
	assert.Contains(t, outputStr, "Usage:")
}

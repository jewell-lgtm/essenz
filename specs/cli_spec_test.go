package specs

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionCommandSpec(t *testing.T) {
	t.Log("SPEC: Version Command")
	t.Log("GIVEN the sz command line tool is available")
	t.Log("WHEN the user runs `sz version`")
	t.Log("THEN the output should display the current version number")

	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", "version")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command should execute successfully")

	outputStr := string(output)
	assert.Contains(t, outputStr, "sz version", "Output should contain version information")
	assert.Contains(t, outputStr, "0.1.0", "Output should contain the correct version number")
}

func TestHelpCommandSpec(t *testing.T) {
	t.Log("SPEC: Help Command")
	t.Log("GIVEN the sz command line tool is available")
	t.Log("WHEN the user runs `sz help`")
	t.Log("THEN the output should display help information")

	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", "help")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command should execute successfully")

	outputStr := string(output)
	assert.Contains(t, outputStr, "sz is a CLI web browser", "Should contain detailed description")
	assert.Contains(t, outputStr, "Available Commands", "Should show available commands")
}

func TestDefaultBehaviorSpec(t *testing.T) {
	t.Log("SPEC: Default Behavior")
	t.Log("GIVEN the sz command line tool is available")
	t.Log("WHEN the user runs `sz` without any arguments")
	t.Log("THEN the output should display the help information by default")

	cmd := exec.Command("go", "run", "../cmd/essenz/main.go")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command should execute successfully")

	outputStr := string(output)
	assert.Contains(t, outputStr, "sz is a CLI web browser", "Should contain tool description")
	assert.Contains(t, outputStr, "Usage:", "Should show usage information")
}

func TestExecutableBinarySpec(t *testing.T) {
	t.Log("SPEC: Executable Binary")
	t.Log("GIVEN the project can be built")
	t.Log("WHEN building the sz binary")
	t.Log("THEN it should compile without errors")

	buildCmd := exec.Command("go", "build", "-o", "sz-test", "../cmd/essenz/main.go")
	buildOutput, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build should succeed: %s", string(buildOutput))

	t.Log("AND WHEN running the built binary")
	runCmd := exec.Command("./sz-test", "version")
	runOutput, err := runCmd.CombinedOutput()
	require.NoError(t, err, "Binary should execute successfully")

	outputStr := string(runOutput)
	assert.Contains(t, outputStr, "sz version 0.1.0", "Binary should output correct version")

	cleanupCmd := exec.Command("rm", "sz-test")
	_ = cleanupCmd.Run()
}

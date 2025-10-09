package specs

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// buildSzBinary builds the sz binary for testing
func buildSzBinary(t *testing.T) string {
	t.Helper()
	_ = os.MkdirAll(".bin", 0o755)
	out := "./.bin/sz-test-" + time.Now().Format("20060102-150405")
	cmd := exec.Command("go", "build", "-o", out, "./cmd/essenz")
	// Set working directory to project root for build
	cmd.Dir = ".."
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Build output: %s", string(output))
	}
	require.NoError(t, err, "Failed to build sz binary for testing")
	return "../.bin/" + filepath.Base(out)
}

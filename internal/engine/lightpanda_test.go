package engine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// cachedLightpanda returns the dev binary path if present, else "".
func cachedLightpanda(t *testing.T) string {
	t.Helper()
	p, err := filepath.Abs(filepath.Join("..", "..", ".DONOTCOMMIT", "lp", "lightpanda"))
	if err != nil {
		return ""
	}
	if _, err := os.Stat(p); err != nil {
		return ""
	}
	return p
}

func TestWriteCookieFileSchema(t *testing.T) {
	path, cleanup, err := writeCookieFile([]Cookie{
		{Name: "a", Value: "1", Domain: "example.com", Path: "/"},
		{Name: "b", Value: "2"},
	})
	require.NoError(t, err)
	defer cleanup()

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	var got []map[string]any
	require.NoError(t, json.Unmarshal(data, &got))
	require.Len(t, got, 2)
	assert.Equal(t, "a", got[0]["name"])
	assert.Equal(t, "1", got[0]["value"])
	assert.Equal(t, "example.com", got[0]["domain"])
	// Empty domain/path omitted.
	_, hasDomain := got[1]["domain"]
	assert.False(t, hasDomain, "empty domain should be omitted")
}

func TestLastLines(t *testing.T) {
	assert.Equal(t, "", lastLines("  \n  ", 3))
	assert.Equal(t, "c; d", lastLines("a\nb\nc\nd", 2))
	assert.Equal(t, "only", lastLines("only", 3))
}

func TestPinnedBuildsCoverCommonPlatforms(t *testing.T) {
	for _, key := range []string{"darwin/arm64", "darwin/amd64", "linux/arm64", "linux/amd64"} {
		b, ok := lightpandaBuilds[key]
		require.True(t, ok, "missing pinned build for %s", key)
		assert.Len(t, b.sha, 64, "sha256 should be 64 hex chars for %s", key)
		assert.Contains(t, b.url, "lightpanda-", "url should reference a lightpanda asset for %s", key)
	}
}

func TestResolveLightpandaHonoursEnvOverride(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "lp-*")
	require.NoError(t, err)
	_ = f.Close()
	t.Setenv("ESSENZ_LIGHTPANDA_PATH", f.Name())

	got, err := resolveLightpanda(context.Background(), false)
	require.NoError(t, err)
	assert.Equal(t, f.Name(), got)
}

func TestResolveLightpandaEnvOverrideMissing(t *testing.T) {
	t.Setenv("ESSENZ_LIGHTPANDA_PATH", filepath.Join(t.TempDir(), "does-not-exist"))
	_, err := resolveLightpanda(context.Background(), false)
	require.Error(t, err)
}

// TestLightpandaEngineRendersJS uses the cached dev binary (if present) to prove
// the subprocess wiring executes JavaScript and returns rendered HTML.
func TestLightpandaEngineRendersJS(t *testing.T) {
	bin := cachedLightpanda(t)
	if bin == "" {
		t.Skip("lightpanda dev binary not present; skipping integration test")
	}
	if runtime.GOOS == "windows" {
		t.Skip("not supported")
	}

	// Content injected only after JS runs — absent from the served HTML.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`<!DOCTYPE html><html><body>
<div id="app"></div>
<script>document.getElementById("app").textContent = "RENDERED_BY_JS";</script>
</body></html>`))
	}))
	defer srv.Close()

	eng := &LightpandaEngine{BinaryPath: bin}
	res, err := eng.Fetch(context.Background(), srv.URL, Options{})
	require.NoError(t, err)
	assert.Equal(t, "lightpanda", res.Engine)
	assert.Contains(t, res.HTML, "RENDERED_BY_JS", "lightpanda should execute the page's JavaScript")
}

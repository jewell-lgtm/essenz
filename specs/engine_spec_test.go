package specs

import (
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SPEC: Pluggable fetch engines
//
// `sz fetch` gains a `--engine` flag selecting how a URL is fetched/rendered:
//
//	auto        (default) escalate: http -> lightpanda -> chrome
//	http        net/http + go-readability  (no JS, lightest)
//	lightpanda  headless JS via `lightpanda fetch` subprocess
//	chrome      existing headless Chrome path (heavy fallback)
//
// All engines yield HTML; essenz's renderer/flags process it uniformly.

const articleHTML = `<!DOCTYPE html><html><head><title>The Essence Article</title></head>
<body>
<nav>Home About Contact</nav>
<article>
<h1>Understanding Lightweight Browsers</h1>
<p>Headless Chrome is powerful but heavy, consuming hundreds of megabytes of
memory for a single page render. Lighter alternatives exist that can execute
JavaScript with a fraction of the footprint, which matters a great deal for a
fast command line tool that fetches and distils web pages into clean markdown.</p>
<p>This paragraph exists so the readability extractor has enough substantive body
text to confidently identify the main content region of the document and return
it rather than discarding the page as boilerplate or navigation chrome.</p>
</article>
<footer>Copyright 2026</footer>
</body></html>`

// runFetch runs the built CLI against args and returns combined output.
func runFetch(t *testing.T, args ...string) (string, error) {
	t.Helper()
	full := append([]string{"run", "../cmd/essenz/main.go", "fetch"}, args...)
	cmd := exec.Command("go", full...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func TestEngineFlagHelpSpec(t *testing.T) {
	t.Log("SPEC: fetch --help advertises the --engine flag")

	out, err := runFetch(t, "--help")
	require.NoError(t, err)
	assert.Contains(t, out, "--engine", "help should document the --engine flag")
	assert.Contains(t, out, "auto", "help should mention the auto engine")
}

func TestHTTPEngineExtractsArticleSpec(t *testing.T) {
	t.Log("SPEC: --engine http fetches via net/http and extracts main content with readability")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(articleHTML))
	}))
	defer srv.Close()

	out, err := runFetch(t, "--engine", "http", srv.URL)
	require.NoError(t, err, "http engine should succeed on a static article")
	assert.Contains(t, out, "Understanding Lightweight Browsers", "should keep the article heading")
	assert.Contains(t, out, "fraction of the footprint", "should keep the article body")
}

func TestUnknownEngineErrorsSpec(t *testing.T) {
	t.Log("SPEC: an unknown --engine value is a clear error")

	out, err := runFetch(t, "--engine", "definitely-not-an-engine", "https://example.com")
	require.Error(t, err, "unknown engine should fail")
	assert.Contains(t, strings.ToLower(out), "engine", "error should mention the engine")
}

func TestAutoEnginePrefersHTTPForStaticPageSpec(t *testing.T) {
	t.Log("SPEC: default (auto) engine serves a server-rendered page from the http tier")
	t.Log("      A page with no <script> tags must NOT escalate to a JS engine.")

	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits++
		_, _ = w.Write([]byte(articleHTML))
	}))
	defer srv.Close()

	// No --engine => auto (default). Must succeed using only the http tier,
	// i.e. without requiring a JS browser to be installed.
	out, err := runFetch(t, srv.URL)
	require.NoError(t, err, "auto should resolve a static page via the http tier")
	assert.Contains(t, out, "Understanding Lightweight Browsers")
	assert.GreaterOrEqual(t, hits, 1, "the page should have been fetched over http")
}

package engine

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testArticle = `<!DOCTYPE html><html><head><title>T</title></head><body>
<nav>Home About</nav>
<article>
<h1>Main Heading Here</h1>
<p>This is a sufficiently long paragraph of genuine article body text so that the
readability algorithm recognises it as the main content of the page rather than
discarding everything as navigation or boilerplate chrome around the edges.</p>
<p>A second substantial paragraph reinforces that the central article region is
the dominant block of meaningful prose on this otherwise simple test document.</p>
</article>
<footer>footer junk</footer></body></html>`

func TestHTTPEngineReturnsRawHTML(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(testArticle))
	}))
	defer srv.Close()

	res, err := NewHTTPEngine().Fetch(context.Background(), srv.URL, Options{})
	require.NoError(t, err)
	assert.Equal(t, "http", res.Engine)
	// The engine returns the full page; extraction happens downstream.
	assert.Contains(t, res.HTML, "Main Heading Here")
	assert.Contains(t, res.HTML, "dominant block of meaningful prose")
	assert.Contains(t, res.HTML, "footer junk", "engine returns raw HTML, not extracted")
	assert.Equal(t, res.HTML, res.Raw, "Raw mirrors HTML for the http engine")
}

func TestHTTPEngineSendsCookiesAndUA(t *testing.T) {
	var gotUA, gotCookie string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		if c, err := r.Cookie("session"); err == nil {
			gotCookie = c.Value
		}
		_, _ = w.Write([]byte("<html><body>ok</body></html>"))
	}))
	defer srv.Close()

	_, err := NewHTTPEngine().Fetch(context.Background(), srv.URL, Options{
		Cookies: []Cookie{{Name: "session", Value: "abc123"}},
	})
	require.NoError(t, err)
	assert.True(t, strings.Contains(gotUA, "Chrome"), "should send a browser-like UA")
	assert.Equal(t, "abc123", gotCookie, "should forward cookies")
}

func TestHTTPEngineErrorsOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, err := NewHTTPEngine().Fetch(context.Background(), srv.URL, Options{})
	require.Error(t, err)
}

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomainConfigValidYAML(t *testing.T) {
	dir := t.TempDir()
	domainsDir := filepath.Join(dir, "essenz", "domains")
	require.NoError(t, os.MkdirAll(domainsDir, 0o755))

	yaml := `cookies:
  - name: session
    value: "abc123"
  - name: token
    value: "xyz"
`
	require.NoError(t, os.WriteFile(filepath.Join(domainsDir, "example.com.yaml"), []byte(yaml), 0o644))

	store := NewWithDir(dir)
	cfg, err := store.DomainConfig("example.com")
	require.NoError(t, err)
	assert.Len(t, cfg.Cookies, 2)
	assert.Equal(t, "session", cfg.Cookies[0].Name)
	assert.Equal(t, "abc123", cfg.Cookies[0].Value)
	assert.Equal(t, "token", cfg.Cookies[1].Name)
}

func TestDomainConfigMissingFile(t *testing.T) {
	store := NewWithDir(t.TempDir())
	cfg, err := store.DomainConfig("nonexistent.com")
	require.NoError(t, err)
	assert.Empty(t, cfg.Cookies)
}

func TestDomainConfigMalformedYAML(t *testing.T) {
	dir := t.TempDir()
	domainsDir := filepath.Join(dir, "essenz", "domains")
	require.NoError(t, os.MkdirAll(domainsDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(domainsDir, "bad.yaml"), []byte("{{invalid"), 0o644))

	store := NewWithDir(dir)
	_, err := store.DomainConfig("bad")
	assert.Error(t, err)
}

func TestCookiesForURL(t *testing.T) {
	dir := t.TempDir()
	domainsDir := filepath.Join(dir, "essenz", "domains")
	require.NoError(t, os.MkdirAll(domainsDir, 0o755))

	yaml := `cookies:
  - name: auth
    value: "tok"
`
	require.NoError(t, os.WriteFile(filepath.Join(domainsDir, "gitlab.example.com.yaml"), []byte(yaml), 0o644))

	store := NewWithDir(dir)
	cookies, err := store.CookiesForURL("https://gitlab.example.com/foo/bar")
	require.NoError(t, err)
	assert.Len(t, cookies, 1)
	assert.Equal(t, "auth", cookies[0].Name)
}

func TestCookiesForURLNoMatch(t *testing.T) {
	store := NewWithDir(t.TempDir())
	cookies, err := store.CookiesForURL("https://other.com/page")
	require.NoError(t, err)
	assert.Empty(t, cookies)
}

func TestCookiesForURLWithPort(t *testing.T) {
	dir := t.TempDir()
	domainsDir := filepath.Join(dir, "essenz", "domains")
	require.NoError(t, os.MkdirAll(domainsDir, 0o755))

	yaml := `cookies:
  - name: s
    value: "v"
`
	require.NoError(t, os.WriteFile(filepath.Join(domainsDir, "127.0.0.1.yaml"), []byte(yaml), 0o644))

	store := NewWithDir(dir)
	cookies, err := store.CookiesForURL("http://127.0.0.1:8080/path")
	require.NoError(t, err)
	assert.Len(t, cookies, 1)
}

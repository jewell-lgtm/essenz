package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// DomainConfig holds per-domain configuration.
type DomainConfig struct {
	Cookies []CookieConfig `yaml:"cookies"`
}

// CookieConfig represents a single cookie to send.
type CookieConfig struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// Store encapsulates config directory location and all file I/O.
type Store struct {
	dir string
}

// New creates a Store using SZ_CONFIG_DIR env var or os.UserConfigDir().
func New() (*Store, error) {
	if d := os.Getenv("SZ_CONFIG_DIR"); d != "" {
		return &Store{dir: d}, nil
	}
	d, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("cannot determine config dir: %w", err)
	}
	return &Store{dir: d}, nil
}

// NewWithDir creates a Store with an explicit base dir (for testing).
func NewWithDir(dir string) *Store {
	return &Store{dir: dir}
}

// DomainConfig returns config for a hostname. Missing file = zero value, nil error.
func (s *Store) DomainConfig(hostname string) (DomainConfig, error) {
	path := filepath.Join(s.dir, "essenz", "domains", hostname+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DomainConfig{}, nil
		}
		return DomainConfig{}, fmt.Errorf("reading domain config: %w", err)
	}

	var cfg DomainConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return DomainConfig{}, fmt.Errorf("parsing domain config for %s: %w", hostname, err)
	}
	return cfg, nil
}

// CookiesForURL parses the URL, looks up domain config, returns cookies.
func (s *Store) CookiesForURL(rawURL string) ([]CookieConfig, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("parsing URL: %w", err)
	}
	hostname := u.Hostname()
	if hostname == "" {
		return nil, nil
	}
	cfg, err := s.DomainConfig(hostname)
	if err != nil {
		return nil, err
	}
	return cfg.Cookies, nil
}

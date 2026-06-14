package engine

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// lightpandaBuild pins one platform's binary URL and its sha256.
type lightpandaBuild struct {
	url string
	sha string
}

// lightpandaBaseURL is the release channel. Lightpanda only publishes a rolling
// "nightly"; the shas below pin a specific verified snapshot. Bump url+sha
// together when updating (see .DONOTCOMMIT/lightpanda-engines-plan.md).
const lightpandaBaseURL = "https://github.com/lightpanda-io/browser/releases/download/nightly/"

// lightpandaBuilds maps GOOS/GOARCH to the pinned binary for that platform.
var lightpandaBuilds = map[string]lightpandaBuild{
	"darwin/arm64": {lightpandaBaseURL + "lightpanda-aarch64-macos", "b353e99635a34ec45d3f594be2ed9903d9a0c976ac652d6ce827ecd3f59a68a7"},
	"darwin/amd64": {lightpandaBaseURL + "lightpanda-x86_64-macos", "53567355b1067da006aa0c83c625ff55baa3151e33092d04a53c741d82d87567"},
	"linux/arm64":  {lightpandaBaseURL + "lightpanda-aarch64-linux", "0935994ca510fc317eef98000d19ab8af17d6dc4ed58fce57cd431e313370401"},
	"linux/amd64":  {lightpandaBaseURL + "lightpanda-x86_64-linux", "7b01a5ed41b6c76e071d056b4b38a39815b1f1772dcf1ade87d41722868b51a2"},
}

// resolveLightpanda finds the lightpanda binary, downloading the pinned build on
// first use. Resolution order: ESSENZ_LIGHTPANDA_PATH -> PATH -> cached download.
func resolveLightpanda(ctx context.Context, debug bool) (string, error) {
	if p := os.Getenv("ESSENZ_LIGHTPANDA_PATH"); p != "" {
		if _, err := os.Stat(p); err != nil {
			return "", fmt.Errorf("ESSENZ_LIGHTPANDA_PATH=%q: %w", p, err)
		}
		return p, nil
	}

	if p, err := exec.LookPath("lightpanda"); err == nil {
		return p, nil
	}

	key := runtime.GOOS + "/" + runtime.GOARCH
	build, ok := lightpandaBuilds[key]
	if !ok {
		return "", fmt.Errorf("no pinned lightpanda build for %s; install lightpanda and set ESSENZ_LIGHTPANDA_PATH", key)
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("locating cache dir: %w", err)
	}
	dir := filepath.Join(cacheDir, "essenz")
	// Encode the sha in the filename so bumping the pin triggers a fresh download.
	dest := filepath.Join(dir, "lightpanda-"+build.sha[:12])

	if verifyFile(dest, build.sha) == nil {
		return dest, nil
	}
	if err := downloadLightpanda(ctx, build, dir, dest, debug); err != nil {
		return "", err
	}
	return dest, nil
}

// downloadLightpanda fetches the pinned binary, verifies its sha256, and installs
// it atomically at dest.
func downloadLightpanda(ctx context.Context, b lightpandaBuild, dir, dest string, debug bool) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	if debug {
		fmt.Fprintf(os.Stderr, "[engine:lightpanda] downloading %s\n", b.url)
	} else {
		fmt.Fprintln(os.Stderr, "essenz: downloading lightpanda browser (first run, ~60MB)...")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, b.url, nil)
	if err != nil {
		return err
	}
	resp, err := (&http.Client{Timeout: 5 * time.Minute}).Do(req)
	if err != nil {
		return fmt.Errorf("downloading lightpanda: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("downloading lightpanda: unexpected status %s", resp.Status)
	}

	tmp, err := os.CreateTemp(dir, "lightpanda-dl-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()

	h := sha256.New()
	if _, err := io.Copy(io.MultiWriter(tmp, h), resp.Body); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("downloading lightpanda: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	if got := hex.EncodeToString(h.Sum(nil)); got != b.sha {
		return fmt.Errorf("lightpanda checksum mismatch: got %s want %s "+
			"(upstream nightly may have changed; update the pin or set ESSENZ_LIGHTPANDA_PATH)", got, b.sha)
	}
	if err := os.Chmod(tmpName, 0o755); err != nil {
		return err
	}
	return os.Rename(tmpName, dest)
}

// verifyFile returns nil if path exists and its sha256 matches want.
func verifyFile(path, want string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	if got := hex.EncodeToString(h.Sum(nil)); got != want {
		return fmt.Errorf("checksum mismatch for %s", path)
	}
	return nil
}

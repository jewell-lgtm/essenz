// Package main provides the sz command line tool for distilling web content into semantic markdown.
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/essenz/essenz/internal/browser"
	"github.com/essenz/essenz/internal/daemon"
	"github.com/spf13/cobra"
)

var version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "sz",
	Short: "Distill the web into semantic markdown",
	Long:  `sz is a CLI web browser that extracts the essence of web pages, reordering content by importance rather than DOM structure.`,
	Run: func(cmd *cobra.Command, _ []string) {
		_ = cmd.Help()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number of sz`,
	Run: func(cmd *cobra.Command, _ []string) {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "sz version %s\n", version)
	},
}

var fetchCmd = &cobra.Command{
	Use:   "fetch [URL or file path]",
	Short: "Fetch content from a URL or local file",
	Long: `Fetch content from an HTTP(S) URL or read from a local file.

Examples:
  sz fetch https://example.com
  sz fetch http://example.com
  sz fetch /path/to/file.html`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]

		// Check if it looks like a URL (simple heuristic)
		if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
			content, err := fetchURLWithChrome(cmd.Context(), target)
			if err != nil {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Error fetching URL: %v\n", err)
				os.Exit(1)
			}
			_, _ = fmt.Fprint(cmd.OutOrStdout(), content)
			return
		}

		// Treat as file path
		content, err := readFile(target)
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Error reading file: %v\n", err)
			os.Exit(1)
		}

		_, _ = fmt.Fprint(cmd.OutOrStdout(), content)
	},
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Manage the Chrome daemon",
	Long:  `Start, stop, or check the status of the Chrome daemon process.`,
}

var daemonStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Chrome daemon",
	Run: func(cmd *cobra.Command, _ []string) {
		server := daemon.NewServer()
		if err := server.Start(); err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Error starting daemon: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Chrome daemon started")

		// Keep the daemon running
		select {}
	},
}

var daemonStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the Chrome daemon",
	Run: func(cmd *cobra.Command, _ []string) {
		client := daemon.NewDaemonClient()
		if err := client.Shutdown(); err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Error stopping daemon: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Chrome daemon stopped")
	},
}

var daemonStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check daemon status",
	Run: func(_ *cobra.Command, _ []string) {
		if daemon.IsDaemonRunning() {
			fmt.Println("Chrome daemon is running")
		} else {
			fmt.Println("Chrome daemon is not running")
		}
	},
}

func init() {
	daemonCmd.AddCommand(daemonStartCmd)
	daemonCmd.AddCommand(daemonStopCmd)
	daemonCmd.AddCommand(daemonStatusCmd)

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(fetchCmd)
	rootCmd.AddCommand(daemonCmd)
}

// readFile reads the contents of a file and returns it as a string
func readFile(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// fetchURLWithChrome fetches content using Chrome browser automation
func fetchURLWithChrome(ctx context.Context, url string) (string, error) {
	client := browser.NewClient()
	defer client.Shutdown()

	content, err := client.FetchContent(ctx, url)
	if err != nil {
		// Fallback to simple HTTP fetch if Chrome fails
		return fetchURL(url)
	}

	return content, nil
}

// fetchURL fetches content from an HTTP or HTTPS URL (fallback method)
func fetchURL(url string) (string, error) {
	// Create HTTP client with reasonable timeout and TLS config for tests
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // For test servers with self-signed certs
			},
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

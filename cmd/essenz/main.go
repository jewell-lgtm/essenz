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

	"github.com/jewell-lgtm/essenz/internal/browser"
	"github.com/jewell-lgtm/essenz/internal/daemon"
	"github.com/jewell-lgtm/essenz/internal/extractor"
	"github.com/jewell-lgtm/essenz/internal/pageready"
	"github.com/jewell-lgtm/essenz/internal/tree"
	"github.com/spf13/cobra"
)

var version = "0.1.0"

// Command line flags
var readerView bool
var rawOutput bool

// DOM ready event flags
var waitForFrameworks bool
var domReadyTimeout string
var waitForSelector string
var debugReadiness bool

// Text node tree flags (F2)
var textNodeTree bool
var treeFormat string
var filterNavigation bool
var preserveAttributes bool

var rootCmd = &cobra.Command{
	Use:   "sz [URL or file path]",
	Short: "Distill the web into semantic markdown",
	Long: `sz is a CLI web browser that extracts the essence of web pages, reordering content by importance rather than DOM structure.

Examples:
  sz https://example.com         # Extract clean content from URL
  sz /path/to/article.html       # Extract clean content from local file
  sz --raw https://example.com   # Get raw HTML without processing
  sz                             # Show this help`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// If no arguments, show help
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}

		target := args[0]
		var content string
		var err error

		// Check if it looks like a URL (simple heuristic)
		if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
			content, err = fetchURLWithChrome(cmd.Context(), target)
			if err != nil {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Error fetching URL: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Treat as file path
			// If DOM ready flags are set, process file through Chrome for consistency
			if shouldUseChromeForFile() {
				// Convert file path to file:// URL and process through Chrome
				fileURL := "file://" + target
				content, err = fetchURLWithChrome(cmd.Context(), fileURL)
				if err != nil {
					// Fallback to direct file reading if Chrome fails
					content, err = readFile(target)
				}
			} else {
				content, err = readFile(target)
			}
			if err != nil {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Error reading file: %v\n", err)
				os.Exit(1)
			}
		}

		// Apply text node tree processing if requested
		if textNodeTree {
			treeBuilder := tree.NewTreeBuilder().
				WithFilterNavigation(filterNavigation).
				WithPreserveAttributes(preserveAttributes)

			root, err := treeBuilder.BuildTree(cmd.Context(), content)
			if err != nil {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Error building text node tree: %v\n", err)
				os.Exit(1)
			}

			// Format output based on tree format flag
			switch treeFormat {
			case "json":
				output, err := treeBuilder.ToJSON(root)
				if err != nil {
					_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Error converting tree to JSON: %v\n", err)
					os.Exit(1)
				}
				content = output
			default:
				content = treeBuilder.ToText(root)
			}

			// Skip reader view processing when text node tree is enabled
			_, _ = fmt.Fprint(cmd.OutOrStdout(), content)
			return
		}

		// Apply reader view processing by default, unless --raw flag is used
		if !rawOutput {
			ext := extractor.New()
			markdown, err := ext.ExtractContent(content)
			if err != nil {
				// Fallback to raw content on extraction error
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Warning: Reader view extraction failed, showing raw content: %v\n", err)
			} else {
				content = markdown
			}
		}

		_, _ = fmt.Fprint(cmd.OutOrStdout(), content)
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
  sz fetch /path/to/file.html
  sz fetch --reader-view https://example.com`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]

		var content string
		var err error

		// Check if it looks like a URL (simple heuristic)
		if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
			content, err = fetchURLWithChrome(cmd.Context(), target)
			if err != nil {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Error fetching URL: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Treat as file path
			// If DOM ready flags are set, process file through Chrome for consistency
			if shouldUseChromeForFile() {
				// Convert file path to file:// URL and process through Chrome
				fileURL := "file://" + target
				content, err = fetchURLWithChrome(cmd.Context(), fileURL)
				if err != nil {
					// Fallback to direct file reading if Chrome fails
					content, err = readFile(target)
				}
			} else {
				content, err = readFile(target)
			}
			if err != nil {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Error reading file: %v\n", err)
				os.Exit(1)
			}
		}

		// Apply text node tree processing if requested
		if textNodeTree {
			treeBuilder := tree.NewTreeBuilder().
				WithFilterNavigation(filterNavigation).
				WithPreserveAttributes(preserveAttributes)

			root, err := treeBuilder.BuildTree(cmd.Context(), content)
			if err != nil {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Error building text node tree: %v\n", err)
				os.Exit(1)
			}

			// Format output based on tree format flag
			switch treeFormat {
			case "json":
				output, err := treeBuilder.ToJSON(root)
				if err != nil {
					_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Error converting tree to JSON: %v\n", err)
					os.Exit(1)
				}
				content = output
			default:
				content = treeBuilder.ToText(root)
			}

			// Skip reader view processing when text node tree is enabled
			_, _ = fmt.Fprint(cmd.OutOrStdout(), content)
			return
		}

		// Apply reader view processing if requested
		if readerView {
			ext := extractor.New()
			markdown, err := ext.ExtractContent(content)
			if err != nil {
				// Fallback to raw content on extraction error
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Warning: Reader view extraction failed, showing raw content: %v\n", err)
			} else {
				content = markdown
			}
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
	// Add daemon subcommands
	daemonCmd.AddCommand(daemonStartCmd)
	daemonCmd.AddCommand(daemonStopCmd)
	daemonCmd.AddCommand(daemonStatusCmd)

	// Add flags to root command
	rootCmd.Flags().BoolVar(&rawOutput, "raw", false, "Output raw HTML without reader view processing")
	rootCmd.Flags().BoolVar(&waitForFrameworks, "wait-for-frameworks", false, "Enable framework-specific readiness detection (React, Vue, Next.js)")
	rootCmd.Flags().StringVar(&domReadyTimeout, "dom-ready-timeout", "5s", "Timeout for DOM readiness detection")
	rootCmd.Flags().StringVar(&waitForSelector, "wait-for-selector", "", "Wait for specific CSS selector to appear before extraction")
	rootCmd.Flags().BoolVar(&debugReadiness, "debug-readiness", false, "Show detailed DOM readiness detection information")

	// Text node tree flags
	rootCmd.Flags().BoolVar(&textNodeTree, "text-node-tree", false, "Build hierarchical text node tree structure")
	rootCmd.Flags().StringVar(&treeFormat, "tree-format", "text", "Output format for text node tree (text, json)")
	rootCmd.Flags().BoolVar(&filterNavigation, "filter-navigation", false, "Filter out navigation elements from tree")
	rootCmd.Flags().BoolVar(&preserveAttributes, "preserve-attributes", false, "Preserve element attributes in tree structure")

	// Add flags to fetch command
	fetchCmd.Flags().BoolVarP(&readerView, "reader-view", "r", false, "Extract main content and convert to clean markdown")
	fetchCmd.Flags().BoolVar(&waitForFrameworks, "wait-for-frameworks", false, "Enable framework-specific readiness detection (React, Vue, Next.js)")
	fetchCmd.Flags().StringVar(&domReadyTimeout, "dom-ready-timeout", "5s", "Timeout for DOM readiness detection")
	fetchCmd.Flags().StringVar(&waitForSelector, "wait-for-selector", "", "Wait for specific CSS selector to appear before extraction")
	fetchCmd.Flags().BoolVar(&debugReadiness, "debug-readiness", false, "Show detailed DOM readiness detection information")

	// Text node tree flags for fetch command
	fetchCmd.Flags().BoolVar(&textNodeTree, "text-node-tree", false, "Build hierarchical text node tree structure")
	fetchCmd.Flags().StringVar(&treeFormat, "tree-format", "text", "Output format for text node tree (text, json)")
	fetchCmd.Flags().BoolVar(&filterNavigation, "filter-navigation", false, "Filter out navigation elements from tree")
	fetchCmd.Flags().BoolVar(&preserveAttributes, "preserve-attributes", false, "Preserve element attributes in tree structure")

	// Add all commands to root
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

// shouldUseChromeForFile determines if file processing should use Chrome
func shouldUseChromeForFile() bool {
	// Use Chrome for files if any DOM ready flags or text node tree flags are set
	return waitForFrameworks || domReadyTimeout != "5s" || waitForSelector != "" || debugReadiness || textNodeTree
}

// createReadinessChecker creates a ReadinessChecker based on CLI flags
func createReadinessChecker() (*pageready.ReadinessChecker, error) {
	// Only create checker if any DOM ready flags are set
	if !waitForFrameworks && domReadyTimeout == "5s" && waitForSelector == "" && !debugReadiness {
		return nil, nil // Use default behavior
	}

	checker := pageready.NewReadinessChecker()

	// Parse timeout
	if domReadyTimeout != "5s" {
		timeout, err := time.ParseDuration(domReadyTimeout)
		if err != nil {
			return nil, fmt.Errorf("invalid timeout format: %w", err)
		}
		checker = checker.WithTimeout(timeout)
	}

	// Set framework hints
	if waitForFrameworks {
		// Enable common framework detection
		checker = checker.WithFrameworkHints([]string{"react", "vue", "angular", "nextjs"})
	}

	// Set custom selectors
	if waitForSelector != "" {
		checker = checker.WithCustomSelectors([]string{waitForSelector})
	}

	// Set debug mode
	checker = checker.WithDebug(debugReadiness)

	return checker, nil
}

// fetchURLWithChrome fetches content using Chrome browser automation
func fetchURLWithChrome(ctx context.Context, url string) (string, error) {
	client := browser.NewClient()
	defer client.Shutdown()

	// Configure DOM readiness if flags are set
	checker, err := createReadinessChecker()
	if err != nil {
		return "", fmt.Errorf("failed to configure DOM readiness: %w", err)
	}

	if checker != nil {
		client = client.WithReadinessChecker(checker)
	}

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

// Package main provides the sz command line tool for distilling web content into semantic markdown.
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

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
			// TODO: Implement HTTP(S) fetching
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "URL fetching not yet implemented: %s\n", target)
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

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(fetchCmd)
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

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

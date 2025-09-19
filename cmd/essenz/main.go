// Package main provides the sz command line tool for distilling web content into semantic markdown.
package main

import (
	"fmt"
	"os"

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
		// TODO: Implement fetch functionality
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Fetch command called with: %s\n", args[0])
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(fetchCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

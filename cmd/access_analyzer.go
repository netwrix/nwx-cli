package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var accessAnalyzerCmd = &cobra.Command{
	Use:     "access-analyzer",
	Aliases: []string{"aa"},
	Short:   "Access Analyzer commands",
	Long:    "Commands for managing Access Analyzer scanners, sources, and scans",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Access Analyzer CLI")
		fmt.Println("Available commands:")
		fmt.Println("  nwx aa config     - Configuration management")
		fmt.Println("  nwx aa scanner    - Scanner management")
		fmt.Println("  nwx aa source     - Source management")
		fmt.Println("  nwx aa scan       - Scan management")
		fmt.Println()
		fmt.Println("Use 'nwx aa <command> --help' for more information about a command.")
	},
}


func init() {
	rootCmd.AddCommand(accessAnalyzerCmd)
}
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  "Display the current version of the nwx CLI tool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("nwx version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
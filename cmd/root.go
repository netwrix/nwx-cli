package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nwx",
	Short: "Netwrix CLI tool",
	Long:  "A command-line interface tool for Netwrix operations and management.",
	Run: func(cmd *cobra.Command, args []string) {
		showIntroScreen()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func showIntroScreen() {
	// NETWRIX ASCII art logo in Vigilant Blue
	fmt.Println("\033[38;2;92;51;255m") // Vigilant Blue RGB (92, 51, 255)
	fmt.Println("███╗   ██╗ ███████╗ ████████╗ ██╗    ██╗ ██████╗  ██╗ ██╗  ██╗")
	fmt.Println("████╗  ██║ ██╔════╝ ╚══██╔══╝ ██║    ██║ ██╔══██╗ ██║ ╚██╗██╔╝")
	fmt.Println("██╔██╗ ██║ █████╗      ██║    ██║ █╗ ██║ ██████╔╝ ██║  ╚███╔╝ ")
	fmt.Println("██║╚██╗██║ ██╔══╝      ██║    ██║███╗██║ ██╔══██╗ ██║  ██╔██╗ ")
	fmt.Println("██║ ╚████║ ███████╗    ██║    ╚███╔███╔╝ ██║  ██║ ██║ ██╔╝ ██╗")
	fmt.Println("╚═╝  ╚═══╝ ╚══════╝    ╚═╝     ╚══╝╚══╝  ╚═╝  ╚═╝ ╚═╝ ╚═╝  ╚═╝")
	fmt.Println()
	fmt.Println("███████╗ ██╗      ██╗")
	fmt.Println("██╔════╝ ██║      ██║")
	fmt.Println("██║      ██║      ██║")
	fmt.Println("██║      ██║      ██║")
	fmt.Println("███████╗ ███████╗ ██║")
	fmt.Println("╚══════╝ ╚══════╝ ╚═╝")
	fmt.Println("\033[0m") // Reset color
	fmt.Println()
	fmt.Println("\033[38;2;255;198;26m🚧 Under Construction 🚧\033[0m") // Signal Yellow
	fmt.Println()
	fmt.Println("\033[38;2;252;250;245mYou've stumbled upon something that doesn't exist yet.\033[0m") // Access White
	fmt.Println("\033[38;2;252;250;245mIf you're curious about what we're building, we'd love to hear from you.\033[0m")
	fmt.Println()
	fmt.Println("\033[38;2;65;242;124mReach out: \033[4mai@netwrix.com\033[0m") // Beacon Green
	fmt.Println()
	fmt.Println("\033[38;2;35;26;64m-- The Netwrix AI Team\033[0m") // Nightwatch
	fmt.Println()
}

func init() {
	// Remove all flags and subcommands - just show the intro
}
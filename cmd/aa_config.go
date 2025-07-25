package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var aaConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Access Analyzer configuration",
	Long:  "Manage Access Analyzer configuration settings",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Access Analyzer Configuration")
		fmt.Println("Available options:")
		fmt.Println("  --endpoint    Set the Access Analyzer API endpoint")
		fmt.Println("  --show        Show current configuration")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  nwx aa config --endpoint=\"http://localhost:3020\"")
		fmt.Println("  nwx aa config --show")
	},
}

var (
	endpointFlag string
	showFlag     bool
)

func init() {
	aaConfigCmd.Flags().StringVar(&endpointFlag, "endpoint", "", "Set the Access Analyzer API endpoint")
	aaConfigCmd.Flags().BoolVar(&showFlag, "show", false, "Show current configuration")
	
	aaConfigCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if endpointFlag != "" {
			if err := setAAEndpoint(endpointFlag); err != nil {
				fmt.Fprintf(os.Stderr, "Error setting endpoint: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Access Analyzer endpoint set to: %s\n", endpointFlag)
			// TODO: Test connection
			fmt.Println("⚠️  Connection test not implemented yet")
		}
		
		if showFlag {
			showAAConfig()
		}
		
		// If no flags provided, show help
		if endpointFlag == "" && !showFlag {
			cmd.Help()
		}
	}
	
	accessAnalyzerCmd.AddCommand(aaConfigCmd)
}

func setAAEndpoint(endpoint string) error {
	configDir, err := getAAConfigDir()
	if err != nil {
		return err
	}
	
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	
	configFile := filepath.Join(configDir, "endpoint")
	return os.WriteFile(configFile, []byte(endpoint), 0644)
}

func getAAEndpoint() (string, error) {
	configDir, err := getAAConfigDir()
	if err != nil {
		return "", err
	}
	
	configFile := filepath.Join(configDir, "endpoint")
	data, err := os.ReadFile(configFile)
	if os.IsNotExist(err) {
		return "", nil // No config file exists yet
	}
	if err != nil {
		return "", err
	}
	
	return string(data), nil
}

func getAAConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".nwx", "access-analyzer"), nil
}

func showAAConfig() {
	fmt.Println("Access Analyzer Configuration:")
	
	endpoint, err := getAAEndpoint()
	if err != nil {
		fmt.Printf("  endpoint: <error: %v>\n", err)
	} else if endpoint == "" {
		fmt.Println("  endpoint: <not configured>")
	} else {
		fmt.Printf("  endpoint: %s\n", endpoint)
	}
}
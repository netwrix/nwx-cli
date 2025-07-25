package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management",
	Long:  "Manage nwx CLI configuration settings",
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long:  "Set a configuration value. Available keys: endpoint",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]
		
		switch key {
		case "endpoint":
			if err := setEndpoint(value); err != nil {
				fmt.Fprintf(os.Stderr, "Error setting endpoint: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Endpoint set to: %s\n", value)
			// TODO: Test connection
			fmt.Println("⚠️  Connection test not implemented yet")
		default:
			fmt.Fprintf(os.Stderr, "Unknown configuration key: %s\n", key)
			fmt.Println("Available keys: endpoint")
			os.Exit(1)
		}
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Long:  "Get a configuration value. Available keys: endpoint",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		
		switch key {
		case "endpoint":
			endpoint, err := getEndpoint()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting endpoint: %v\n", err)
				os.Exit(1)
			}
			if endpoint == "" {
				fmt.Println("No endpoint configured")
			} else {
				fmt.Printf("Current endpoint: %s\n", endpoint)
			}
		default:
			fmt.Fprintf(os.Stderr, "Unknown configuration key: %s\n", key)
			fmt.Println("Available keys: endpoint")
			os.Exit(1)
		}
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration",
	Long:  "Display all current configuration settings",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Current configuration:")
		
		endpoint, err := getEndpoint()
		if err != nil {
			fmt.Printf("  endpoint: <error: %v>\n", err)
		} else if endpoint == "" {
			fmt.Println("  endpoint: <not configured>")
		} else {
			fmt.Printf("  endpoint: %s\n", endpoint)
		}
	},
}

func setEndpoint(endpoint string) error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}
	
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	
	configFile := filepath.Join(configDir, "config")
	return os.WriteFile(configFile, []byte(endpoint), 0644)
}

func getEndpoint() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	
	configFile := filepath.Join(configDir, "config")
	data, err := os.ReadFile(configFile)
	if os.IsNotExist(err) {
		return "", nil // No config file exists yet
	}
	if err != nil {
		return "", err
	}
	
	return string(data), nil
}

func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".nwx"), nil
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)
	rootCmd.AddCommand(configCmd)
}
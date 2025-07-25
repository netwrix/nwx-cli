package cmd

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// Styles for the interactive interface
var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7C3AED")).
		Padding(0, 1)

	menuStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#7C3AED")).
		Padding(0, 1).
		MarginTop(1)

	selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7C3AED")).
		Background(lipgloss.Color("#E0E7FF")).
		Padding(0, 1).
		Bold(true)

	normalStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#374151")).
		Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		MarginTop(1)

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#EF4444")).
		Bold(true)

	successStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#10B981")).
		Bold(true)
)

// InteractiveModel represents the main interactive CLI model
type InteractiveModel struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
	mode     string
	done     bool
	err      error
}

// MenuItem represents a menu item with description
type MenuItem struct {
	Title       string
	Description string
	Action      func() error
}

// Menu item generators to avoid initialization cycles
func getMainMenuItems() []MenuItem {
	return []MenuItem{
		{
			Title:       "Access Analyzer",
			Description: "Manage Access Analyzer configuration and scanners",
			Action: func() error {
				return runAccessAnalyzerMenu()
			},
		},
		{
			Title:       "Configuration",
			Description: "Configure CLI settings and endpoints",
			Action: func() error {
				return runConfigMenu()
			},
		},
		{
			Title:       "Help",
			Description: "Show help information",
			Action: func() error {
				return runHelpMenu()
			},
		},
		{
			Title:       "Exit",
			Description: "Exit the CLI",
			Action: func() error {
				os.Exit(0)
				return nil
			},
		},
	}
}

func getAccessAnalyzerMenuItems() []MenuItem {
	return []MenuItem{
		{
			Title:       "Scanner Management",
			Description: "Create, list, and manage scanners",
			Action: func() error {
				return runScannerMenu()
			},
		},
		{
			Title:       "← Back to Main Menu",
			Description: "Return to main menu",
			Action: func() error {
				return runMainMenu()
			},
		},
	}
}

func getScannerMenuItems() []MenuItem {
	return []MenuItem{
		{
			Title:       "Create Scanner",
			Description: "Interactive scanner creation workflow",
			Action: func() error {
				return runScannerCreation()
			},
		},
		{
			Title:       "← Back to Access Analyzer",
			Description: "Return to Access Analyzer menu",
			Action: func() error {
				return runAccessAnalyzerMenu()
			},
		},
	}
}

// InteractiveCommand creates the main interactive command
var InteractiveCommand = &cobra.Command{
	Use:   "interactive",
	Short: "Start interactive CLI mode",
	Long:  "Start the interactive CLI interface for a guided experience",
	Run: func(cmd *cobra.Command, args []string) {
		// Show intro logo
		showIntroLogo()
		
		// Try interactive mode first, fall back to simple menu if not available
		if err := runMainMenu(); err != nil {
			if strings.Contains(err.Error(), "TTY") || strings.Contains(err.Error(), "device not configured") {
				// Fall back to simple text-based menu
				runSimpleMenu()
			} else {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		}
	},
}

// initialModel creates the initial model for the interactive CLI
func initialModel(items []MenuItem) InteractiveModel {
	choices := make([]string, len(items))
	for i, item := range items {
		choices[i] = item.Title
	}

	return InteractiveModel{
		choices:  choices,
		selected: make(map[int]struct{}),
		mode:     "menu",
	}
}

// Init initializes the model
func (m InteractiveModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m InteractiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ":
			if m.mode == "menu" {
				m.done = true
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

// View renders the model
func (m InteractiveModel) View() string {
	if m.done {
		return ""
	}

	var s strings.Builder

	// Title with styled header
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#5C33FF")).
		Bold(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#5C33FF")).
		Padding(0, 2).
		MarginBottom(1)
	
	s.WriteString(headerStyle.Render("NETWRIX CLI"))
	s.WriteString("\n\n")

	// Menu items
	for i, choice := range m.choices {
		if m.cursor == i {
			s.WriteString(selectedStyle.Render("→ " + choice))
		} else {
			s.WriteString(normalStyle.Render("  " + choice))
		}
		s.WriteString("\n")
	}

	// Help text
	s.WriteString(helpStyle.Render("\nNavigation: ↑/↓ or j/k to move, enter to select, q to quit"))

	return s.String()
}

// Menu runner functions
func runMainMenu() error {
	return runMenu(getMainMenuItems(), "Main Menu")
}

func runAccessAnalyzerMenu() error {
	return runMenu(getAccessAnalyzerMenuItems(), "Access Analyzer")
}

func runScannerMenu() error {
	return runMenu(getScannerMenuItems(), "Scanner Management")
}

func runConfigMenu() error {
	fmt.Println(menuStyle.Render("📋 Configuration Menu"))
	fmt.Println("Configuration options coming soon...")
	fmt.Println(helpStyle.Render("Press any key to continue..."))
	fmt.Scanln()
	return runMainMenu()
}

func runAAConfigMenu() error {
	fmt.Println(menuStyle.Render("⚙️  Access Analyzer Configuration"))
	
	// Show current endpoint
	endpoint, err := getAAEndpoint()
	if err != nil || endpoint == "" {
		fmt.Println(errorStyle.Render("No Access Analyzer endpoint configured"))
	} else {
		fmt.Printf("Current endpoint: %s\n", successStyle.Render(endpoint))
	}
	
	fmt.Print("\nEnter new endpoint (or press Enter to keep current): ")
	var newEndpoint string
	fmt.Scanln(&newEndpoint)
	
	if newEndpoint != "" {
		if err := setAAEndpoint(newEndpoint); err != nil {
			fmt.Printf("Error setting endpoint: %v\n", err)
		} else {
			fmt.Println(successStyle.Render("✅ Endpoint updated successfully"))
		}
	}
	
	fmt.Println(helpStyle.Render("Press any key to continue..."))
	fmt.Scanln()
	return runAccessAnalyzerMenu()
}

func runHelpMenu() error {
	fmt.Println(menuStyle.Render("📚 Help"))
	fmt.Println(`
Welcome to NWX CLI - Interactive Mode!

This CLI provides tools for managing Access Analyzer scanners and configuration.

Main Features:
• 🔍 Access Analyzer - Scanner management and configuration
• ⚙️  Configuration - CLI settings and endpoints
• 📋 Interactive workflows - Guided scanner creation

Navigation:
• Use ↑/↓ or j/k to navigate menus
• Press Enter to select an item
• Press 'q' to quit at any time

For more information, visit: https://github.com/netwrix/nwx-cli
`)
	fmt.Println(helpStyle.Render("Press any key to continue..."))
	fmt.Scanln()
	return runMainMenu()
}

func runStatusCommand() error {
	fmt.Println(menuStyle.Render("🔍 Access Analyzer Status"))
	
	client, err := getAPIClient()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		fmt.Println(helpStyle.Render("Press any key to continue..."))
		fmt.Scanln()
		return runAccessAnalyzerMenu()
	}
	
	fmt.Printf("🔗 Endpoint: %s\n", client.BaseURL)
	
	if err := client.TestConnection(); err != nil {
		fmt.Printf("❌ Connection failed: %v\n", err)
	} else {
		fmt.Println(successStyle.Render("✅ Connection successful"))
	}
	
	fmt.Println(helpStyle.Render("Press any key to continue..."))
	fmt.Scanln()
	return runAccessAnalyzerMenu()
}

func runScannerCreation() error {
	fmt.Println(menuStyle.Render("🚀 Scanner Creation"))
	fmt.Println("Starting interactive scanner creation workflow...")
	
	// Check if endpoint is configured
	client, err := getAPIClient()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		fmt.Println(helpStyle.Render("Press any key to continue..."))
		fmt.Scanln()
		return runScannerMenu()
	}
	
	fmt.Printf("🔍 Connecting to Access Analyzer at: %s\n", client.BaseURL)
	
	// Test connection first
	if err := client.TestConnection(); err != nil {
		fmt.Printf("❌ Connection failed: %v\n", err)
		fmt.Println(helpStyle.Render("Press any key to continue..."))
		fmt.Scanln()
		return runScannerMenu()
	}
	
	// Get existing scanners
	response, err := client.GetSourceTypes()
	if err != nil {
		fmt.Printf("⚠️  Could not fetch existing scanners: %v\n", err)
	} else {
		fmt.Printf("✅ Found %d existing scanners\n", len(response.Data))
	}
	
	fmt.Println()
	
	// Start interactive scanner creation workflow
	if err := runInteractiveScannerCreation(response); err != nil {
		if err.Error() == "EOF" {
			fmt.Println("❌ Scanner creation cancelled by user")
		} else {
			fmt.Printf("❌ Scanner creation failed: %v\n", err)
		}
	}
	
	fmt.Println(helpStyle.Render("Press any key to continue..."))
	fmt.Scanln()
	return runScannerMenu()
}

// runMenu runs a generic menu with the given items
func runMenu(items []MenuItem, title string) error {
	model := initialModel(items)
	
	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return err
	}
	
	// Get the final model and execute the selected action
	if m, ok := finalModel.(InteractiveModel); ok && m.done {
		if m.cursor < len(items) {
			return items[m.cursor].Action()
		}
	}
	
	return nil
}

// showIntroLogo displays the NETWRIX CLI intro logo
func showIntroLogo() {
	// Clear screen
	fmt.Print("\033[2J\033[H")
	
	// NETWRIX logo in Vigilant Blue
	logoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#5C33FF")).
		Bold(true)
	
	fmt.Println(logoStyle.Render(""))
	fmt.Println(logoStyle.Render("███╗   ██╗███████╗████████╗██╗    ██╗██████╗ ██╗██╗  ██╗"))
	fmt.Println(logoStyle.Render("████╗  ██║██╔════╝╚══██╔══╝██║    ██║██╔══██╗██║╚██╗██╔╝"))
	fmt.Println(logoStyle.Render("██╔██╗ ██║█████╗     ██║   ██║ █╗ ██║██████╔╝██║ ╚███╔╝ "))
	fmt.Println(logoStyle.Render("██║╚██╗██║██╔══╝     ██║   ██║███╗██║██╔══██╗██║ ██╔██╗ "))
	fmt.Println(logoStyle.Render("██║ ╚████║███████╗   ██║   ╚███╔███╔╝██║  ██║██║██╔╝ ██╗"))
	fmt.Println(logoStyle.Render("╚═╝  ╚═══╝╚══════╝   ╚═╝    ╚══╝╚══╝ ╚═╝  ╚═╝╚═╝╚═╝  ╚═╝"))
	fmt.Println(logoStyle.Render(""))
	fmt.Println(logoStyle.Render("         ██████╗██╗     ██╗                            "))
	fmt.Println(logoStyle.Render("        ██╔════╝██║     ██║                            "))
	fmt.Println(logoStyle.Render("        ██║     ██║     ██║                            "))
	fmt.Println(logoStyle.Render("        ██║     ██║     ██║                            "))
	fmt.Println(logoStyle.Render("        ╚██████╗███████╗██║                            "))
	fmt.Println(logoStyle.Render("         ╚═════╝╚══════╝╚═╝                            "))
	fmt.Println()
	
	// Subtitle
	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		Italic(true)
	
	fmt.Println(subtitleStyle.Render("              Interactive Command Line Interface"))
	fmt.Println()
	
	// Version or tagline
	taglineStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#10B981")).
		Bold(true)
	
	fmt.Println(taglineStyle.Render("         🚀 Access Analyzer Scanner Management"))
	fmt.Println()
	
	// Brief pause for dramatic effect
	fmt.Print("Press any key to continue...")
	fmt.Scanln()
	fmt.Print("\033[2J\033[H") // Clear screen again
}

// runSimpleMenu provides a fallback text-based menu when TTY is not available
func runSimpleMenu() {
	fmt.Println("🚀 NWX CLI - Interactive Mode")
	fmt.Println("=============================")
	fmt.Println()
	
	items := getMainMenuItems()
	
	for {
		fmt.Println("Main Menu:")
		for i, item := range items {
			fmt.Printf("%d. %s - %s\n", i+1, item.Title, item.Description)
		}
		fmt.Println()
		fmt.Print("Select an option (1-4): ")
		
		var choice int
		if _, err := fmt.Scanln(&choice); err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}
		
		if choice < 1 || choice > len(items) {
			fmt.Println("Invalid choice. Please select a valid option.")
			continue
		}
		
		selectedItem := items[choice-1]
		fmt.Printf("Selected: %s\n\n", selectedItem.Title)
		
		if selectedItem.Title == "Exit" {
			fmt.Println("Goodbye!")
			return
		}
		
		// For now, just show what would happen
		fmt.Printf("This would execute: %s\n", selectedItem.Title)
		fmt.Println("(Full implementation available in proper TTY environment)")
		fmt.Println()
	}
}

func init() {
	// Add the interactive command to the root command
	rootCmd.AddCommand(InteractiveCommand)
}
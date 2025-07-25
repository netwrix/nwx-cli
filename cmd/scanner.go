package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var scannerCmd = &cobra.Command{
	Use:   "scanner",
	Short: "Scanner management",
	Long:  "Manage Access Analyzer scanners",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Scanner Management")
		fmt.Println("Available commands:")
		fmt.Println("  nwx aa scanner --create    - Create a new scanner interactively")
		fmt.Println()
		fmt.Println("Use 'nwx aa scanner <command> --help' for more information.")
	},
}

var scannerCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new scanner",
	Long:  "Interactive scanner creation workflow",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Interactive Scanner Creation")
		fmt.Println("=" + strings.Repeat("=", 35))
		fmt.Println()
		
		// Check if endpoint is configured
		client, err := getAPIClient()
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		
		fmt.Printf("🔍 Connecting to Access Analyzer at: %s\n", client.BaseURL)
		
		// Test connection first
		if err := client.TestConnection(); err != nil {
			fmt.Printf("❌ Connection failed: %v\n", err)
			return
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
				return
			}
			fmt.Printf("❌ Scanner creation failed: %v\n", err)
			return
		}
	},
}



// Add --create flag to scanner command
var createFlag bool

// ScannerCreationData holds the data collected during scanner creation
type ScannerCreationData struct {
	// Basic Information
	Name        string
	DisplayName string
	Description string
	Version     string
	Icon        string
	Language    string
	
	// Scan Types
	SupportedScanTypes []string
	
	// Connection Configuration
	AuthMethods []string
	
	// File Generation
	GenerateFiles bool
	OutputDir     string
}

// runInteractiveScannerCreation runs the interactive scanner creation workflow
func runInteractiveScannerCreation(existingScanners *SourceTypeListResponse) error {
	scanner := &ScannerCreationData{}
	
	// Step 1: Basic Information
	if err := collectBasicInfo(scanner, existingScanners); err != nil {
		return err
	}
	
	// Step 2: Language Selection
	if err := collectLanguage(scanner); err != nil {
		return err
	}
	
	// Step 3: Scan Types
	if err := collectScanTypes(scanner); err != nil {
		return err
	}
	
	// Step 4: Authentication Methods
	if err := collectAuthMethods(scanner); err != nil {
		return err
	}
	
	// Step 5: File Generation Options
	if err := collectFileGeneration(scanner); err != nil {
		return err
	}
	
	// Step 6: Summary and Confirmation
	if err := showSummaryAndConfirm(scanner); err != nil {
		return err
	}
	
	// Step 7: Generate Files
	if scanner.GenerateFiles {
		return generateScannerFiles(scanner)
	}
	
	fmt.Println("✅ Scanner configuration completed!")
	return nil
}

// collectBasicInfo collects basic scanner information
func collectBasicInfo(scanner *ScannerCreationData, existing *SourceTypeListResponse) error {
	fmt.Println("📋 Step 1: Basic Information")
	fmt.Println()
	
	// Get existing scanner names for validation
	existingNames := make(map[string]bool)
	if existing != nil {
		for _, s := range existing.Data {
			existingNames[s.TypeName] = true
		}
	}
	
	// Scanner name (kebab-case)
	namePrompt := &survey.Input{
		Message: "Scanner name (kebab-case, e.g., 'my-scanner'):",
		Help:    "This will be used as the technical identifier",
	}
	if err := survey.AskOne(namePrompt, &scanner.Name, survey.WithValidator(func(val interface{}) error {
		if str := val.(string); str != "" {
			if existingNames[str] {
				return fmt.Errorf("scanner name '%s' already exists", str)
			}
			// Basic kebab-case validation
			if !strings.Contains(str, "-") && len(str) > 0 {
				return fmt.Errorf("scanner name should be kebab-case (e.g., 'my-scanner')")
			}
		}
		return nil
	})); err != nil {
		return err
	}
	
	// Display name
	displayPrompt := &survey.Input{
		Message: "Display name:",
		Help:    "Human-readable name shown in the UI",
		Default: strings.Title(strings.ReplaceAll(scanner.Name, "-", " ")),
	}
	if err := survey.AskOne(displayPrompt, &scanner.DisplayName); err != nil {
		return err
	}
	
	// Description
	descPrompt := &survey.Input{
		Message: "Description:",
		Help:    "Brief description of what this scanner does",
	}
	if err := survey.AskOne(descPrompt, &scanner.Description); err != nil {
		return err
	}
	
	// Version
	versionPrompt := &survey.Input{
		Message: "Version:",
		Default: "1.0.0",
		Help:    "Semantic version (e.g., 1.0.0)",
	}
	if err := survey.AskOne(versionPrompt, &scanner.Version); err != nil {
		return err
	}
	
	// Icon
	iconOptions := []string{"folder", "database", "cloud", "server", "lock", "file", "network", "other"}
	iconPrompt := &survey.Select{
		Message: "Choose an icon:",
		Options: iconOptions,
		Default: "folder",
	}
	if err := survey.AskOne(iconPrompt, &scanner.Icon); err != nil {
		return err
	}
	
	fmt.Println()
	return nil
}

// collectLanguage collects the programming language for the scanner
func collectLanguage(scanner *ScannerCreationData) error {
	fmt.Println("💻 Step 2: Programming Language")
	fmt.Println()
	
	languageOptions := []string{"python", "javascript", "go", "java", "c#"}
	languagePrompt := &survey.Select{
		Message: "Select programming language:",
		Options: languageOptions,
		Default: "python",
		Help:    "Choose the programming language for your scanner implementation",
	}
	
	if err := survey.AskOne(languagePrompt, &scanner.Language); err != nil {
		return err
	}
	
	fmt.Println()
	return nil
}

// collectScanTypes collects supported scan types
func collectScanTypes(scanner *ScannerCreationData) error {
	fmt.Println("🔍 Step 3: Scan Types")
	fmt.Println()
	
	scanTypeOptions := []string{"access", "sensitive_data"}
	scanTypePrompt := &survey.MultiSelect{
		Message: "Select supported scan types:",
		Options: scanTypeOptions,
		Default: []string{"access"},
		Help:    "Use space to select/deselect, enter to confirm",
	}
	
	if err := survey.AskOne(scanTypePrompt, &scanner.SupportedScanTypes); err != nil {
		return err
	}
	
	fmt.Println()
	return nil
}

// collectAuthMethods collects authentication methods
func collectAuthMethods(scanner *ScannerCreationData) error {
	fmt.Println("🔐 Step 4: Authentication Methods")
	fmt.Println()
	
	authOptions := []string{
		"Username/Password",
		"API Key",
		"OAuth2",
		"Certificate",
		"Service Account",
		"Windows Authentication",
		"Custom",
	}
	
	authPrompt := &survey.MultiSelect{
		Message: "Select authentication methods:",
		Options: authOptions,
		Default: []string{"Username/Password"},
		Help:    "Use space to select/deselect, enter to confirm",
	}
	
	if err := survey.AskOne(authPrompt, &scanner.AuthMethods); err != nil {
		return err
	}
	
	fmt.Println()
	return nil
}

// collectFileGeneration collects file generation options
func collectFileGeneration(scanner *ScannerCreationData) error {
	fmt.Println("📁 Step 5: File Generation")
	fmt.Println()
	
	generatePrompt := &survey.Confirm{
		Message: "Generate scanner files in current directory?",
		Default: true,
		Help:    "This will create the scanner structure, Dockerfile, and example code",
	}
	
	if err := survey.AskOne(generatePrompt, &scanner.GenerateFiles); err != nil {
		return err
	}
	
	if scanner.GenerateFiles {
		dirPrompt := &survey.Input{
			Message: "Output directory:",
			Default: "./" + scanner.Name,
			Help:    "Directory where scanner files will be generated",
		}
		if err := survey.AskOne(dirPrompt, &scanner.OutputDir); err != nil {
			return err
		}
	}
	
	fmt.Println()
	return nil
}

// showSummaryAndConfirm shows a summary and asks for confirmation
func showSummaryAndConfirm(scanner *ScannerCreationData) error {
	fmt.Println("📊 Step 6: Summary")
	fmt.Println()
	
	fmt.Printf("Name:          %s\n", scanner.Name)
	fmt.Printf("Display Name:  %s\n", scanner.DisplayName)
	fmt.Printf("Description:   %s\n", scanner.Description)
	fmt.Printf("Version:       %s\n", scanner.Version)
	fmt.Printf("Icon:          %s\n", scanner.Icon)
	fmt.Printf("Language:      %s\n", scanner.Language)
	fmt.Printf("Scan Types:    %s\n", strings.Join(scanner.SupportedScanTypes, ", "))
	fmt.Printf("Auth Methods:  %s\n", strings.Join(scanner.AuthMethods, ", "))
	
	if scanner.GenerateFiles {
		fmt.Printf("Output Dir:    %s\n", scanner.OutputDir)
	}
	
	fmt.Println()
	
	confirmPrompt := &survey.Confirm{
		Message: "Create scanner with these settings?",
		Default: true,
	}
	
	var confirmed bool
	if err := survey.AskOne(confirmPrompt, &confirmed); err != nil {
		return err
	}
	
	if !confirmed {
		return fmt.Errorf("scanner creation cancelled")
	}
	
	return nil
}

// generateScannerFiles generates the scanner files
func generateScannerFiles(scanner *ScannerCreationData) error {
	fmt.Printf("🚀 Generating scanner files in: %s\n", scanner.OutputDir)
	
	// Create output directory
	if err := os.MkdirAll(scanner.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Create config directory
	configDir := filepath.Join(scanner.OutputDir, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	// Generate all files
	files := []struct {
		name    string
		content string
	}{
		{"scannerSpecification.json", generateScannerSpecification(scanner)},
		{"Dockerfile", generateDockerfile(scanner)},
		{"README.md", generateReadme(scanner)},
		{"config/config.example.json", generateConfigExample(scanner)},
		{fmt.Sprintf("%s-source-type.json", scanner.Name), generateSourceType(scanner)},
	}
	
	// Add language-specific files
	switch scanner.Language {
	case "python":
		files = append(files, 
			struct{name, content string}{"requirements.txt", generateRequirements(scanner)},
			struct{name, content string}{"scanner.py", generateScannerPython(scanner)},
		)
	case "javascript":
		files = append(files, 
			struct{name, content string}{"package.json", generatePackageJson(scanner)},
			struct{name, content string}{"scanner.js", generateScannerJavaScript(scanner)},
		)
	case "go":
		files = append(files, 
			struct{name, content string}{"go.mod", generateGoMod(scanner)},
			struct{name, content string}{"scanner.go", generateScannerGo(scanner)},
		)
	case "java":
		files = append(files, 
			struct{name, content string}{"pom.xml", generatePomXml(scanner)},
			struct{name, content string}{"Scanner.java", generateScannerJava(scanner)},
		)
	case "c#":
		files = append(files, 
			struct{name, content string}{"Scanner.csproj", generateCsProj(scanner)},
			struct{name, content string}{"Scanner.cs", generateScannerCSharp(scanner)},
		)
	}
	
	for _, file := range files {
		filePath := filepath.Join(scanner.OutputDir, file.name)
		if err := os.WriteFile(filePath, []byte(file.content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", file.name, err)
		}
		fmt.Printf("  ✅ Created %s\n", file.name)
	}
	
	fmt.Println()
	fmt.Println("✅ Scanner files generated successfully!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  1. cd %s\n", scanner.OutputDir)
	fmt.Println("  2. Review and customize the generated files")
	fmt.Println("  3. Update scanner.py with your specific implementation")
	fmt.Println("  4. Test your scanner: docker build -t my-scanner .")
	fmt.Println("  5. Deploy to Access Analyzer")
	
	return nil
}

// generateScannerSpecification generates the scannerSpecification.json file
func generateScannerSpecification(scanner *ScannerCreationData) string {
	spec := map[string]interface{}{
		"name":    strings.ToUpper(strings.ReplaceAll(scanner.Name, "-", "_")),
		"version": scanner.Version,
		"connectionConfig": map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"key":         "host",
					"label":       "Host",
					"type":        "text",
					"required":    true,
					"placeholder": "example.com",
					"description": "Host to connect to",
				},
			},
		},
		"outputSchema": generateMinimalOutputSchema(scanner),
	}
	
	// Add minimal access scan config for access scanners
	if contains(scanner.SupportedScanTypes, "access") {
		spec["accessScanConfig"] = map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"key":         "scanDepth",
					"label":       "Scan Depth",
					"type":        "number",
					"required":    false,
					"default":     10,
					"min":         1,
					"max":         100,
					"description": "Maximum scan depth",
				},
			},
		}
	}
	
	// Add minimal sensitive data scan config for sensitive data scanners
	if contains(scanner.SupportedScanTypes, "sensitive_data") {
		spec["sensitiveDataScanConfig"] = map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"key":         "maxFileSize",
					"label":       "Max File Size (MB)",
					"type":        "number",
					"required":    false,
					"default":     100,
					"min":         1,
					"max":         1000,
					"description": "Maximum file size to scan",
				},
			},
		}
	}
	
	data, _ := json.MarshalIndent(spec, "", "  ")
	return string(data)
}


// generateMinimalOutputSchema generates a minimal output schema
func generateMinimalOutputSchema(scanner *ScannerCreationData) map[string]interface{} {
	schema := map[string]interface{}{}
	
	if contains(scanner.SupportedScanTypes, "access") {
		schema["access"] = map[string]interface{}{
			"columns": []map[string]interface{}{
				{
					"name":        "scan_id",
					"type":        "string",
					"maxLength":   36,
					"nullable":    false,
					"primaryKey":  true,
					"description": "Unique identifier for the scan run",
				},
				{
					"name":        "resource_id",
					"type":        "string",
					"maxLength":   255,
					"nullable":    false,
					"primaryKey":  true,
					"description": "Unique identifier for the resource",
				},
				{
					"name":        "scan_timestamp",
					"type":        "timestamp",
					"nullable":    false,
					"defaultValue": "CURRENT_TIMESTAMP",
					"description": "When this scan record was created",
				},
			},
		}
	}
	
	if contains(scanner.SupportedScanTypes, "sensitive_data") {
		schema["sensitiveData"] = map[string]interface{}{
			"columns": []map[string]interface{}{
				{
					"name":        "scan_id",
					"type":        "string",
					"maxLength":   36,
					"nullable":    false,
					"primaryKey":  true,
					"description": "Unique identifier for the scan run",
				},
				{
					"name":        "match_id",
					"type":        "string",
					"maxLength":   36,
					"nullable":    false,
					"primaryKey":  true,
					"description": "Unique identifier for the match",
				},
				{
					"name":        "scan_timestamp",
					"type":        "timestamp",
					"nullable":    false,
					"defaultValue": "CURRENT_TIMESTAMP",
					"description": "When this scan record was created",
				},
			},
		}
	}
	
	return schema
}


// generateDockerfile generates the Dockerfile
func generateDockerfile(scanner *ScannerCreationData) string {
	switch scanner.Language {
	case "python":
		return `FROM python:3.11-slim

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    python3-dev \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy and install requirements
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy scanner files
COPY scanner.py .
COPY scannerSpecification.json .

# Set default environment variables
ENV RABBITMQ_HOST=rabbitmq
ENV RABBITMQ_PORT=5672
ENV RABBITMQ_USER=guest
ENV RABBITMQ_PASSWORD=guest

ENV APP_DB_HOST=postgres-app
ENV APP_DB_PORT=5432
ENV APP_DB_NAME=app
ENV APP_DB_USER=appuser
ENV APP_DB_PASSWORD=app_password

# Collection database is ClickHouse (NOT PostgreSQL)
ENV COLLECTION_DB_HOST=clickhouse
ENV COLLECTION_DB_PORT=9000
ENV COLLECTION_DB_NAME=default
ENV COLLECTION_DB_USER=default
ENV COLLECTION_DB_PASSWORD=

CMD ["python", "scanner.py"]
`
	case "javascript":
		return `FROM node:18-alpine

# Install system dependencies
RUN apk add --no-cache \
    gcc \
    musl-dev \
    postgresql-dev

WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm install

# Copy scanner files
COPY scanner.js .
COPY scannerSpecification.json .

# Set default environment variables
ENV RABBITMQ_HOST=rabbitmq
ENV RABBITMQ_PORT=5672
ENV RABBITMQ_USER=guest
ENV RABBITMQ_PASSWORD=guest

ENV APP_DB_HOST=postgres-app
ENV APP_DB_PORT=5432
ENV APP_DB_NAME=app
ENV APP_DB_USER=appuser
ENV APP_DB_PASSWORD=app_password

ENV COLLECTION_DB_HOST=clickhouse
ENV COLLECTION_DB_PORT=9000
ENV COLLECTION_DB_NAME=default
ENV COLLECTION_DB_USER=default
ENV COLLECTION_DB_PASSWORD=

CMD ["node", "scanner.js"]
`
	case "go":
		return `FROM golang:1.21-alpine AS builder

# Install system dependencies
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy scanner files
COPY scanner.go .
COPY scannerSpecification.json .

# Build the application
RUN go build -o scanner scanner.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy binary and config
COPY --from=builder /app/scanner .
COPY --from=builder /app/scannerSpecification.json .

# Set default environment variables
ENV RABBITMQ_HOST=rabbitmq
ENV RABBITMQ_PORT=5672
ENV RABBITMQ_USER=guest
ENV RABBITMQ_PASSWORD=guest

ENV APP_DB_HOST=postgres-app
ENV APP_DB_PORT=5432
ENV APP_DB_NAME=app
ENV APP_DB_USER=appuser
ENV APP_DB_PASSWORD=app_password

ENV COLLECTION_DB_HOST=clickhouse
ENV COLLECTION_DB_PORT=9000
ENV COLLECTION_DB_NAME=default
ENV COLLECTION_DB_USER=default
ENV COLLECTION_DB_PASSWORD=

CMD ["./scanner"]
`
	case "java":
		return `FROM openjdk:17-jdk-slim

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy Maven files
COPY pom.xml .
COPY src ./src

# Copy scanner config
COPY scannerSpecification.json .

# Build the application
RUN ./mvnw clean package -DskipTests

# Set default environment variables
ENV RABBITMQ_HOST=rabbitmq
ENV RABBITMQ_PORT=5672
ENV RABBITMQ_USER=guest
ENV RABBITMQ_PASSWORD=guest

ENV APP_DB_HOST=postgres-app
ENV APP_DB_PORT=5432
ENV APP_DB_NAME=app
ENV APP_DB_USER=appuser
ENV APP_DB_PASSWORD=app_password

ENV COLLECTION_DB_HOST=clickhouse
ENV COLLECTION_DB_PORT=9000
ENV COLLECTION_DB_NAME=default
ENV COLLECTION_DB_USER=default
ENV COLLECTION_DB_PASSWORD=

CMD ["java", "-jar", "target/scanner.jar"]
`
	case "c#":
		return `FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build

WORKDIR /app

# Copy project files
COPY *.csproj ./
RUN dotnet restore

# Copy source code
COPY . .
RUN dotnet publish -c Release -o out

# Runtime stage
FROM mcr.microsoft.com/dotnet/aspnet:8.0

WORKDIR /app

# Copy built application
COPY --from=build /app/out .

# Set default environment variables
ENV RABBITMQ_HOST=rabbitmq
ENV RABBITMQ_PORT=5672
ENV RABBITMQ_USER=guest
ENV RABBITMQ_PASSWORD=guest

ENV APP_DB_HOST=postgres-app
ENV APP_DB_PORT=5432
ENV APP_DB_NAME=app
ENV APP_DB_USER=appuser
ENV APP_DB_PASSWORD=app_password

ENV COLLECTION_DB_HOST=clickhouse
ENV COLLECTION_DB_PORT=9000
ENV COLLECTION_DB_NAME=default
ENV COLLECTION_DB_USER=default
ENV COLLECTION_DB_PASSWORD=

CMD ["dotnet", "Scanner.dll"]
`
	default:
		// Default to Python
		return `FROM python:3.11-slim

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    python3-dev \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy and install requirements
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy scanner files
COPY scanner.py .
COPY scannerSpecification.json .

# Set default environment variables
ENV RABBITMQ_HOST=rabbitmq
ENV RABBITMQ_PORT=5672
ENV RABBITMQ_USER=guest
ENV RABBITMQ_PASSWORD=guest

ENV APP_DB_HOST=postgres-app
ENV APP_DB_PORT=5432
ENV APP_DB_NAME=app
ENV APP_DB_USER=appuser
ENV APP_DB_PASSWORD=app_password

ENV COLLECTION_DB_HOST=clickhouse
ENV COLLECTION_DB_PORT=9000
ENV COLLECTION_DB_NAME=default
ENV COLLECTION_DB_USER=default
ENV COLLECTION_DB_PASSWORD=

CMD ["python", "scanner.py"]
`
	}
}

// generateRequirements generates requirements.txt
func generateRequirements(scanner *ScannerCreationData) string {
	return `# Core dependencies for Access Analyzer scanners
pika>=1.3.0
psycopg2-binary>=2.9.0
clickhouse-driver>=0.2.6
requests>=2.31.0

# Add your scanner-specific dependencies here
# For example:
# boto3>=1.26.0  # for AWS scanners
# azure-identity>=1.15.0  # for Azure scanners
# pysmb>=1.2.9  # for SMB/CIFS scanners
`
}

// generateScannerPython generates the minimal Python scanner template
func generateScannerPython(scanner *ScannerCreationData) string {
	return fmt.Sprintf(`#!/usr/bin/env python3
"""
%s Scanner
%s

Minimal scanner scaffolding for Access Analyzer.
"""

import os
import json
import logging
import pika
import psycopg2
from clickhouse_driver import Client
from datetime import datetime

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class %sScanner:
    def __init__(self, config):
        self.config = config
        self.connection_config = config.get('connectionConfig', {})
        self.scan_config = config.get('accessScanConfig', {})
        
        # Results storage
        self.results = []
        
        # Scan metadata (set by queue scanner)
        self.scan_id = None
        self.source_id = None
        self.scan_table_name = None
    
    def connect(self):
        """Establish connection to the data source"""
        # TODO: Implement connection logic
        host = self.connection_config.get('host')
        logger.info(f"Connecting to {host}...")
        
        # TODO: Add your connection implementation here
        pass
    
    def scan(self):
        """Perform the actual scan"""
        logger.info(f"Starting scan with ID: {self.scan_id}")
        
        try:
            # Connect to the data source
            self.connect()
            
            # TODO: Implement your scanning logic here
            # Example:
            # for resource in enumerate_resources():
            #     result = {
            #         'scan_id': self.scan_id,
            #         'resource_id': resource.id,
            #         'scan_timestamp': datetime.now()
            #     }
            #     self.results.append(result)
            
            logger.info(f"Scan completed. Found {len(self.results)} resources.")
            
        except Exception as e:
            logger.error(f"Scan failed: {e}")
            raise
    
    def save_results_to_db(self, db_config):
        """Save scan results to ClickHouse collection database"""
        logger.info(f"Saving {len(self.results)} results to database")
        
        # TODO: Implement database saving logic
        # See scanner framework documentation for examples
        pass

class QueueScanner:
    def __init__(self):
        # Load scanner specification
        with open('scannerSpecification.json', 'r') as f:
            spec = json.load(f)
            self.scanner_name = spec['name']
            self.scanner_version = spec['version']
            
        # Queue and table names
        self.scan_queue_name = f'{self.scanner_name}-{self.scanner_version}-scan-access'
        self.test_queue_name = f'{self.scanner_name}-{self.scanner_version}-test'
        
        version_with_underscores = self.scanner_version.replace('.', '_')
        self.table_name = f'{self.scanner_name}_{version_with_underscores}_access'.lower()
        
        # Database configurations from environment
        self.app_db_config = {
            'host': os.environ.get('APP_DB_HOST'),
            'port': os.environ.get('APP_DB_PORT'),
            'database': os.environ.get('APP_DB_NAME'),
            'user': os.environ.get('APP_DB_USER'),
            'password': os.environ.get('APP_DB_PASSWORD')
        }
        
        self.collection_db_config = {
            'host': os.environ.get('COLLECTION_DB_HOST', 'clickhouse'),
            'port': os.environ.get('COLLECTION_DB_PORT', '9000'),
            'database': os.environ.get('COLLECTION_DB_NAME', 'default'),
            'user': os.environ.get('COLLECTION_DB_USER', 'default'),
            'password': os.environ.get('COLLECTION_DB_PASSWORD', '')
        }
    
    def connect_rabbitmq(self):
        """Connect to RabbitMQ"""
        # TODO: Implement RabbitMQ connection
        # See scanner framework documentation for examples
        pass
    
    def update_scan_status(self, scan_id, status, error_message=None):
        """Update scan status in app database"""
        # TODO: Implement status update logic
        pass
    
    def process_scan_job(self, ch, method, properties, body):
        """Process a scan job from the queue"""
        # TODO: Implement scan job processing
        pass
    
    def process_test_job(self, ch, method, properties, body):
        """Process a test connection job from the queue"""
        # TODO: Implement test connection logic
        pass
    
    def run(self):
        """Start the scanner"""
        logger.info("Starting %s scanner...")
        
        # TODO: Implement scanner startup logic
        pass

if __name__ == "__main__":
    scanner = QueueScanner()
    scanner.run()
`,
		scanner.DisplayName,
		scanner.Description,
		toPascalCase(scanner.Name),
		scanner.DisplayName,
	)
}

// generateReadme generates the README.md file
func generateReadme(scanner *ScannerCreationData) string {
	return fmt.Sprintf(`# %s Scanner

%s

## Getting Started

This is a minimal scanner scaffolding for Access Analyzer. You'll need to implement the actual scanning logic.

## Files Generated

- `+"`scannerSpecification.json`"+` - Scanner configuration schema
- `+"`scanner.py`"+` - Main scanner implementation (minimal scaffolding)
- `+"`Dockerfile`"+` - Docker container configuration
- `+"`requirements.txt`"+` - Python dependencies
- `+"`config/config.example.json`"+` - Configuration example

## Next Steps

1. **Review the scanner specification** in `+"`scannerSpecification.json`"+`
2. **Implement your scanning logic** in `+"`scanner.py`"+`
3. **Add your specific dependencies** to `+"`requirements.txt`"+`
4. **Update the configuration** in `+"`config/config.example.json`"+`
5. **Test your implementation** with `+"`python scanner.py`"+`
6. **Build and deploy** with `+"`docker build -t %s-scanner .`"+`

## Documentation

See the scanner framework documentation for detailed implementation guidance:
- Connection handling
- Database integration
- Queue processing
- Error handling
- Testing

## TODO

- [ ] Implement connection logic in `+"`scanner.py`"+`
- [ ] Add scanning logic for your data source
- [ ] Implement result processing
- [ ] Add error handling
- [ ] Test with real data
- [ ] Add logging and monitoring
`,
		scanner.DisplayName,
		scanner.Description,
		scanner.Name,
	)
}

// generateConfigExample generates a minimal configuration example
func generateConfigExample(scanner *ScannerCreationData) string {
	config := map[string]interface{}{
		"connectionConfig": map[string]interface{}{
			"host": "example.com",
		},
		"accessScanConfig": map[string]interface{}{
			"scanDepth": 10,
		},
	}
	
	data, _ := json.MarshalIndent(config, "", "  ")
	return string(data)
}

// generateSourceType generates the source type definition
func generateSourceType(scanner *ScannerCreationData) string {
	sourceType := map[string]interface{}{
		"displayName":        scanner.DisplayName,
		"description":        scanner.Description,
		"icon":               scanner.Icon,
		"scannerImage":       fmt.Sprintf("access-analyzer/%s-scanner:latest", scanner.Name),
		"supportedScanTypes": scanner.SupportedScanTypes,
		"scannerSpecification": map[string]string{
			"$ref": "scannerSpecification.json",
		},
	}
	
	data, _ := json.MarshalIndent(sourceType, "", "  ")
	return string(data)
}

// Helper functions
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func toPascalCase(s string) string {
	words := strings.Split(s, "-")
	result := ""
	for _, word := range words {
		if len(word) > 0 {
			result += strings.Title(word)
		}
	}
	return result
}

// generatePackageJson generates package.json for JavaScript
func generatePackageJson(scanner *ScannerCreationData) string {
	return fmt.Sprintf(`{
  "name": "%s-scanner",
  "version": "%s",
  "description": "%s",
  "main": "scanner.js",
  "dependencies": {
    "amqplib": "^0.10.3",
    "pg": "^8.11.3",
    "@clickhouse/client": "^0.2.7"
  },
  "scripts": {
    "start": "node scanner.js"
  }
}`, scanner.Name, scanner.Version, scanner.Description)
}

// generateGoMod generates go.mod for Go
func generateGoMod(scanner *ScannerCreationData) string {
	return fmt.Sprintf(`module %s-scanner

go 1.21

require (
	github.com/lib/pq v1.10.9
	github.com/streadway/amqp v1.1.0
	github.com/ClickHouse/clickhouse-go/v2 v2.15.0
)
`, scanner.Name)
}

// generatePomXml generates pom.xml for Java
func generatePomXml(scanner *ScannerCreationData) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>
    
    <groupId>com.accessanalyzer</groupId>
    <artifactId>%s-scanner</artifactId>
    <version>%s</version>
    <packaging>jar</packaging>
    
    <properties>
        <maven.compiler.source>17</maven.compiler.source>
        <maven.compiler.target>17</maven.compiler.target>
        <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
    </properties>
    
    <dependencies>
        <dependency>
            <groupId>com.rabbitmq</groupId>
            <artifactId>amqp-client</artifactId>
            <version>5.19.0</version>
        </dependency>
        <dependency>
            <groupId>org.postgresql</groupId>
            <artifactId>postgresql</artifactId>
            <version>42.6.0</version>
        </dependency>
        <dependency>
            <groupId>com.clickhouse</groupId>
            <artifactId>clickhouse-jdbc</artifactId>
            <version>0.4.6</version>
        </dependency>
    </dependencies>
    
    <build>
        <plugins>
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-compiler-plugin</artifactId>
                <version>3.11.0</version>
                <configuration>
                    <source>17</source>
                    <target>17</target>
                </configuration>
            </plugin>
        </plugins>
    </build>
</project>`, scanner.Name, scanner.Version)
}

// generateCsProj generates Scanner.csproj for C#
func generateCsProj(scanner *ScannerCreationData) string {
	return fmt.Sprintf(`<Project Sdk="Microsoft.NET.Sdk">

  <PropertyGroup>
    <OutputType>Exe</OutputType>
    <TargetFramework>net8.0</TargetFramework>
    <ImplicitUsings>enable</ImplicitUsings>
    <Nullable>enable</Nullable>
    <AssemblyName>%s-Scanner</AssemblyName>
    <AssemblyVersion>%s</AssemblyVersion>
  </PropertyGroup>

  <ItemGroup>
    <PackageReference Include="RabbitMQ.Client" Version="6.6.0" />
    <PackageReference Include="Npgsql" Version="7.0.6" />
    <PackageReference Include="ClickHouse.Client" Version="7.1.0" />
    <PackageReference Include="Newtonsoft.Json" Version="13.0.3" />
  </ItemGroup>

</Project>`, scanner.Name, scanner.Version)
}

// generateScannerJavaScript generates scanner.js for JavaScript
func generateScannerJavaScript(scanner *ScannerCreationData) string {
	return fmt.Sprintf(`/**
 * %s Scanner
 * %s
 * 
 * Minimal scanner scaffolding for Access Analyzer.
 */

const amqp = require('amqplib');
const { Client } = require('pg');
const { createClient } = require('@clickhouse/client');
const fs = require('fs');

class %sScanner {
    constructor(config) {
        this.config = config;
        this.connectionConfig = config.connectionConfig || {};
        this.scanConfig = config.accessScanConfig || {};
        
        // Results storage
        this.results = [];
        
        // Scan metadata (set by queue scanner)
        this.scanId = null;
        this.sourceId = null;
        this.scanTableName = null;
    }
    
    async connect() {
        // TODO: Implement connection logic
        const host = this.connectionConfig.host;
        console.log('Connecting to', host);
        
        // TODO: Add your connection implementation here
    }
    
    async scan() {
        console.log('Starting scan with ID:', this.scanId);
        
        try {
            // Connect to the data source
            await this.connect();
            
            // TODO: Implement your scanning logic here
            // Example:
            // for (const resource of await this.enumerateResources()) {
            //     const result = {
            //         scan_id: this.scanId,
            //         resource_id: resource.id,
            //         scan_timestamp: new Date()
            //     };
            //     this.results.push(result);
            // }
            
            console.log('Scan completed. Found', this.results.length, 'resources.');
            
        } catch (error) {
            console.error('Scan failed:', error);
            throw error;
        }
    }
    
    async saveResultsToDb(dbConfig) {
        console.log('Saving', this.results.length, 'results to database');
        
        // TODO: Implement database saving logic
        // See scanner framework documentation for examples
    }
}

class QueueScanner {
    constructor() {
        // Load scanner specification
        const spec = JSON.parse(fs.readFileSync('scannerSpecification.json', 'utf8'));
        this.scannerName = spec.name;
        this.scannerVersion = spec.version;
        
        // Queue and table names
        this.scanQueueName = this.scannerName + '-' + this.scannerVersion + '-scan-access';
        this.testQueueName = this.scannerName + '-' + this.scannerVersion + '-test';
        
        const versionWithUnderscores = this.scannerVersion.replace(/\./g, '_');
        this.tableName = (this.scannerName + '_' + versionWithUnderscores + '_access').toLowerCase();
        
        // Database configurations from environment
        this.appDbConfig = {
            host: process.env.APP_DB_HOST,
            port: process.env.APP_DB_PORT,
            database: process.env.APP_DB_NAME,
            user: process.env.APP_DB_USER,
            password: process.env.APP_DB_PASSWORD
        };
        
        this.collectionDbConfig = {
            host: process.env.COLLECTION_DB_HOST || 'clickhouse',
            port: process.env.COLLECTION_DB_PORT || '9000',
            database: process.env.COLLECTION_DB_NAME || 'default',
            username: process.env.COLLECTION_DB_USER || 'default',
            password: process.env.COLLECTION_DB_PASSWORD || ''
        };
    }
    
    async connectRabbitMQ() {
        // TODO: Implement RabbitMQ connection
        // See scanner framework documentation for examples
    }
    
    async updateScanStatus(scanId, status, errorMessage = null) {
        // TODO: Implement status update logic
    }
    
    async processScanJob(msg) {
        // TODO: Implement scan job processing
    }
    
    async processTestJob(msg) {
        // TODO: Implement test connection logic
    }
    
    async run() {
        console.log('Starting %s scanner...');
        
        // TODO: Implement scanner startup logic
    }
}

// Start the scanner
if (require.main === module) {
    const scanner = new QueueScanner();
    scanner.run().catch(console.error);
}
`, scanner.DisplayName, scanner.Description, toPascalCase(scanner.Name), scanner.DisplayName)
}

// generateScannerGo generates scanner.go for Go
func generateScannerGo(scanner *ScannerCreationData) string {
	return fmt.Sprintf(`package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

// %s Scanner
// %s
// 
// Minimal scanner scaffolding for Access Analyzer.

type %sScanner struct {
	config           map[string]interface{}
	connectionConfig map[string]interface{}
	scanConfig       map[string]interface{}
	
	// Results storage
	results []map[string]interface{}
	
	// Scan metadata (set by queue scanner)
	scanID        string
	sourceID      string
	scanTableName string
}

func New%sScanner(config map[string]interface{}) *%sScanner {
	return &%sScanner{
		config:           config,
		connectionConfig: getMapFromConfig(config, "connectionConfig"),
		scanConfig:       getMapFromConfig(config, "accessScanConfig"),
		results:          make([]map[string]interface{}, 0),
	}
}

func (s *%sScanner) Connect() error {
	// TODO: Implement connection logic
	host, _ := s.connectionConfig["host"].(string)
	fmt.Printf("Connecting to %%s...\\n", host)
	
	// TODO: Add your connection implementation here
	return nil
}

func (s *%sScanner) Scan() error {
	fmt.Printf("Starting scan with ID: %%s\\n", s.scanID)
	
	// Connect to the data source
	if err := s.Connect(); err != nil {
		return err
	}
	
	// TODO: Implement your scanning logic here
	// Example:
	// for _, resource := range s.enumerateResources() {
	//     result := map[string]interface{}{
	//         "scan_id":        s.scanID,
	//         "resource_id":    resource.ID,
	//         "scan_timestamp": time.Now(),
	//     }
	//     s.results = append(s.results, result)
	// }
	
	fmt.Printf("Scan completed. Found %%d resources.\\n", len(s.results))
	return nil
}

func (s *%sScanner) SaveResultsToDB(dbConfig map[string]interface{}) error {
	fmt.Printf("Saving %%d results to database\\n", len(s.results))
	
	// TODO: Implement database saving logic
	// See scanner framework documentation for examples
	return nil
}

type QueueScanner struct {
	scannerName    string
	scannerVersion string
	scanQueueName  string
	testQueueName  string
	tableName      string
	
	appDbConfig        map[string]interface{}
	collectionDbConfig map[string]interface{}
}

func NewQueueScanner() *QueueScanner {
	// Load scanner specification
	specData, err := ioutil.ReadFile("scannerSpecification.json")
	if err != nil {
		log.Fatal("Error reading scannerSpecification.json:", err)
	}
	
	var spec map[string]interface{}
	if err := json.Unmarshal(specData, &spec); err != nil {
		log.Fatal("Error parsing scannerSpecification.json:", err)
	}
	
	name := spec["name"].(string)
	version := spec["version"].(string)
	
	// Queue and table names
	scanQueueName := name + "-" + version + "-scan-access"
	testQueueName := name + "-" + version + "-test"
	
	versionWithUnderscores := strings.ReplaceAll(version, ".", "_")
	tableName := strings.ToLower(name + "_" + versionWithUnderscores + "_access")
	
	return &QueueScanner{
		scannerName:    name,
		scannerVersion: version,
		scanQueueName:  scanQueueName,
		testQueueName:  testQueueName,
		tableName:      tableName,
		
		appDbConfig: map[string]interface{}{
			"host":     os.Getenv("APP_DB_HOST"),
			"port":     os.Getenv("APP_DB_PORT"),
			"database": os.Getenv("APP_DB_NAME"),
			"user":     os.Getenv("APP_DB_USER"),
			"password": os.Getenv("APP_DB_PASSWORD"),
		},
		
		collectionDbConfig: map[string]interface{}{
			"host":     getEnvWithDefault("COLLECTION_DB_HOST", "clickhouse"),
			"port":     getEnvWithDefault("COLLECTION_DB_PORT", "9000"),
			"database": getEnvWithDefault("COLLECTION_DB_NAME", "default"),
			"username": getEnvWithDefault("COLLECTION_DB_USER", "default"),
			"password": getEnvWithDefault("COLLECTION_DB_PASSWORD", ""),
		},
	}
}

func (qs *QueueScanner) ConnectRabbitMQ() error {
	// TODO: Implement RabbitMQ connection
	// See scanner framework documentation for examples
	return nil
}

func (qs *QueueScanner) UpdateScanStatus(scanID, status, errorMessage string) error {
	// TODO: Implement status update logic
	return nil
}

func (qs *QueueScanner) ProcessScanJob(message []byte) error {
	// TODO: Implement scan job processing
	return nil
}

func (qs *QueueScanner) ProcessTestJob(message []byte) error {
	// TODO: Implement test connection logic
	return nil
}

func (qs *QueueScanner) Run() error {
	fmt.Printf("Starting %%s scanner...\\n", qs.scannerName)
	
	// TODO: Implement scanner startup logic
	return nil
}

// Helper functions
func getMapFromConfig(config map[string]interface{}, key string) map[string]interface{} {
	if value, exists := config[key]; exists {
		if mapValue, ok := value.(map[string]interface{}); ok {
			return mapValue
		}
	}
	return make(map[string]interface{})
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	scanner := NewQueueScanner()
	if err := scanner.Run(); err != nil {
		log.Fatal("Scanner failed:", err)
	}
}
`, scanner.DisplayName, scanner.Description, toPascalCase(scanner.Name), toPascalCase(scanner.Name), toPascalCase(scanner.Name), toPascalCase(scanner.Name), toPascalCase(scanner.Name), toPascalCase(scanner.Name), toPascalCase(scanner.Name))
}

// generateScannerJava generates Scanner.java for Java
func generateScannerJava(scanner *ScannerCreationData) string {
	return fmt.Sprintf(`/**
 * %s Scanner
 * %s
 * 
 * Minimal scanner scaffolding for Access Analyzer.
 */

import java.util.*;
import java.sql.*;
import java.io.*;
import java.time.LocalDateTime;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.rabbitmq.client.*;

public class %sScanner {
    private Map<String, Object> config;
    private Map<String, Object> connectionConfig;
    private Map<String, Object> scanConfig;
    
    // Results storage
    private List<Map<String, Object>> results;
    
    // Scan metadata (set by queue scanner)
    private String scanId;
    private String sourceId;
    private String scanTableName;
    
    public %sScanner(Map<String, Object> config) {
        this.config = config;
        this.connectionConfig = getMapFromConfig(config, "connectionConfig");
        this.scanConfig = getMapFromConfig(config, "accessScanConfig");
        this.results = new ArrayList<>();
    }
    
    public void connect() throws Exception {
        // TODO: Implement connection logic
        String host = (String) connectionConfig.get("host");
        System.out.println("Connecting to " + host + "...");
        
        // TODO: Add your connection implementation here
    }
    
    public void scan() throws Exception {
        System.out.println("Starting scan with ID: " + scanId);
        
        try {
            // Connect to the data source
            connect();
            
            // TODO: Implement your scanning logic here
            // Example:
            // for (Resource resource : enumerateResources()) {
            //     Map<String, Object> result = new HashMap<>();
            //     result.put("scan_id", scanId);
            //     result.put("resource_id", resource.getId());
            //     result.put("scan_timestamp", LocalDateTime.now());
            //     results.add(result);
            // }
            
            System.out.println("Scan completed. Found " + results.size() + " resources.");
            
        } catch (Exception e) {
            System.err.println("Scan failed: " + e.getMessage());
            throw e;
        }
    }
    
    public void saveResultsToDb(Map<String, Object> dbConfig) throws Exception {
        System.out.println("Saving " + results.size() + " results to database");
        
        // TODO: Implement database saving logic
        // See scanner framework documentation for examples
    }
    
    // Getters and setters
    public void setScanId(String scanId) { this.scanId = scanId; }
    public void setSourceId(String sourceId) { this.sourceId = sourceId; }
    public void setScanTableName(String scanTableName) { this.scanTableName = scanTableName; }
    
    @SuppressWarnings("unchecked")
    private Map<String, Object> getMapFromConfig(Map<String, Object> config, String key) {
        Object value = config.get(key);
        if (value instanceof Map) {
            return (Map<String, Object>) value;
        }
        return new HashMap<>();
    }
}

class QueueScanner {
    private String scannerName;
    private String scannerVersion;
    private String scanQueueName;
    private String testQueueName;
    private String tableName;
    
    private Map<String, Object> appDbConfig;
    private Map<String, Object> collectionDbConfig;
    
    public QueueScanner() throws Exception {
        // Load scanner specification
        ObjectMapper mapper = new ObjectMapper();
        Map<String, Object> spec = mapper.readValue(new File("scannerSpecification.json"), Map.class);
        
        this.scannerName = (String) spec.get("name");
        this.scannerVersion = (String) spec.get("version");
        
        // Queue and table names
        this.scanQueueName = scannerName + "-" + scannerVersion + "-scan-access";
        this.testQueueName = scannerName + "-" + scannerVersion + "-test";
        
        String versionWithUnderscores = scannerVersion.replace(".", "_");
        this.tableName = (scannerName + "_" + versionWithUnderscores + "_access").toLowerCase();
        
        // Database configurations from environment
        this.appDbConfig = new HashMap<>();
        appDbConfig.put("host", System.getenv("APP_DB_HOST"));
        appDbConfig.put("port", System.getenv("APP_DB_PORT"));
        appDbConfig.put("database", System.getenv("APP_DB_NAME"));
        appDbConfig.put("user", System.getenv("APP_DB_USER"));
        appDbConfig.put("password", System.getenv("APP_DB_PASSWORD"));
        
        this.collectionDbConfig = new HashMap<>();
        collectionDbConfig.put("host", getEnvWithDefault("COLLECTION_DB_HOST", "clickhouse"));
        collectionDbConfig.put("port", getEnvWithDefault("COLLECTION_DB_PORT", "9000"));
        collectionDbConfig.put("database", getEnvWithDefault("COLLECTION_DB_NAME", "default"));
        collectionDbConfig.put("username", getEnvWithDefault("COLLECTION_DB_USER", "default"));
        collectionDbConfig.put("password", getEnvWithDefault("COLLECTION_DB_PASSWORD", ""));
    }
    
    public void connectRabbitMQ() throws Exception {
        // TODO: Implement RabbitMQ connection
        // See scanner framework documentation for examples
    }
    
    public void updateScanStatus(String scanId, String status, String errorMessage) throws Exception {
        // TODO: Implement status update logic
    }
    
    public void processScanJob(byte[] message) throws Exception {
        // TODO: Implement scan job processing
    }
    
    public void processTestJob(byte[] message) throws Exception {
        // TODO: Implement test connection logic
    }
    
    public void run() throws Exception {
        System.out.println("Starting " + scannerName + " scanner...");
        
        // TODO: Implement scanner startup logic
    }
    
    private String getEnvWithDefault(String key, String defaultValue) {
        String value = System.getenv(key);
        return value != null ? value : defaultValue;
    }
    
    public static void main(String[] args) {
        try {
            QueueScanner scanner = new QueueScanner();
            scanner.run();
        } catch (Exception e) {
            System.err.println("Scanner failed: " + e.getMessage());
            e.printStackTrace();
        }
    }
}
`, scanner.DisplayName, scanner.Description, toPascalCase(scanner.Name), toPascalCase(scanner.Name))
}

// generateScannerCSharp generates Scanner.cs for C#
func generateScannerCSharp(scanner *ScannerCreationData) string {
	return fmt.Sprintf(`using System;
using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;
using Newtonsoft.Json;
using RabbitMQ.Client;
using Npgsql;

/// <summary>
/// %s Scanner
/// %s
/// 
/// Minimal scanner scaffolding for Access Analyzer.
/// </summary>
namespace AccessAnalyzer.Scanner
{
    public class %sScanner
    {
        private readonly Dictionary<string, object> _config;
        private readonly Dictionary<string, object> _connectionConfig;
        private readonly Dictionary<string, object> _scanConfig;
        
        // Results storage
        private readonly List<Dictionary<string, object>> _results;
        
        // Scan metadata (set by queue scanner)
        public string ScanId { get; set; }
        public string SourceId { get; set; }
        public string ScanTableName { get; set; }
        
        public %sScanner(Dictionary<string, object> config)
        {
            _config = config;
            _connectionConfig = GetMapFromConfig(config, "connectionConfig");
            _scanConfig = GetMapFromConfig(config, "accessScanConfig");
            _results = new List<Dictionary<string, object>>();
        }
        
        public async Task ConnectAsync()
        {
            // TODO: Implement connection logic
            var host = _connectionConfig.GetValueOrDefault("host", "").ToString();
            Console.WriteLine($"Connecting to {host}...");
            
            // TODO: Add your connection implementation here
        }
        
        public async Task ScanAsync()
        {
            Console.WriteLine($"Starting scan with ID: {ScanId}");
            
            try
            {
                // Connect to the data source
                await ConnectAsync();
                
                // TODO: Implement your scanning logic here
                // Example:
                // foreach (var resource in await EnumerateResourcesAsync())
                // {
                //     var result = new Dictionary<string, object>
                //     {
                //         ["scan_id"] = ScanId,
                //         ["resource_id"] = resource.Id,
                //         ["scan_timestamp"] = DateTime.Now
                //     };
                //     _results.Add(result);
                // }
                
                Console.WriteLine($"Scan completed. Found {_results.Count} resources.");
            }
            catch (Exception ex)
            {
                Console.Error.WriteLine($"Scan failed: {ex.Message}");
                throw;
            }
        }
        
        public async Task SaveResultsToDbAsync(Dictionary<string, object> dbConfig)
        {
            Console.WriteLine($"Saving {_results.Count} results to database");
            
            // TODO: Implement database saving logic
            // See scanner framework documentation for examples
        }
        
        private Dictionary<string, object> GetMapFromConfig(Dictionary<string, object> config, string key)
        {
            if (config.TryGetValue(key, out var value) && value is Dictionary<string, object> map)
            {
                return map;
            }
            return new Dictionary<string, object>();
        }
    }
    
    public class QueueScanner
    {
        private readonly string _scannerName;
        private readonly string _scannerVersion;
        private readonly string _scanQueueName;
        private readonly string _testQueueName;
        private readonly string _tableName;
        
        private readonly Dictionary<string, object> _appDbConfig;
        private readonly Dictionary<string, object> _collectionDbConfig;
        
        public QueueScanner()
        {
            // Load scanner specification
            var specJson = File.ReadAllText("scannerSpecification.json");
            var spec = JsonConvert.DeserializeObject<Dictionary<string, object>>(specJson);
            
            _scannerName = spec["name"].ToString();
            _scannerVersion = spec["version"].ToString();
            
            // Queue and table names
            _scanQueueName = $"{_scannerName}-{_scannerVersion}-scan-access";
            _testQueueName = $"{_scannerName}-{_scannerVersion}-test";
            
            var versionWithUnderscores = _scannerVersion.Replace(".", "_");
            _tableName = $"{_scannerName}_{versionWithUnderscores}_access".ToLower();
            
            // Database configurations from environment
            _appDbConfig = new Dictionary<string, object>
            {
                ["host"] = Environment.GetEnvironmentVariable("APP_DB_HOST"),
                ["port"] = Environment.GetEnvironmentVariable("APP_DB_PORT"),
                ["database"] = Environment.GetEnvironmentVariable("APP_DB_NAME"),
                ["user"] = Environment.GetEnvironmentVariable("APP_DB_USER"),
                ["password"] = Environment.GetEnvironmentVariable("APP_DB_PASSWORD")
            };
            
            _collectionDbConfig = new Dictionary<string, object>
            {
                ["host"] = GetEnvWithDefault("COLLECTION_DB_HOST", "clickhouse"),
                ["port"] = GetEnvWithDefault("COLLECTION_DB_PORT", "9000"),
                ["database"] = GetEnvWithDefault("COLLECTION_DB_NAME", "default"),
                ["username"] = GetEnvWithDefault("COLLECTION_DB_USER", "default"),
                ["password"] = GetEnvWithDefault("COLLECTION_DB_PASSWORD", "")
            };
        }
        
        public async Task ConnectRabbitMQAsync()
        {
            // TODO: Implement RabbitMQ connection
            // See scanner framework documentation for examples
        }
        
        public async Task UpdateScanStatusAsync(string scanId, string status, string errorMessage = null)
        {
            // TODO: Implement status update logic
        }
        
        public async Task ProcessScanJobAsync(byte[] message)
        {
            // TODO: Implement scan job processing
        }
        
        public async Task ProcessTestJobAsync(byte[] message)
        {
            // TODO: Implement test connection logic
        }
        
        public async Task RunAsync()
        {
            Console.WriteLine($"Starting {_scannerName} scanner...");
            
            // TODO: Implement scanner startup logic
        }
        
        private string GetEnvWithDefault(string key, string defaultValue)
        {
            return Environment.GetEnvironmentVariable(key) ?? defaultValue;
        }
        
        public static async Task Main(string[] args)
        {
            try
            {
                var scanner = new QueueScanner();
                await scanner.RunAsync();
            }
            catch (Exception ex)
            {
                Console.Error.WriteLine($"Scanner failed: {ex.Message}");
                Environment.Exit(1);
            }
        }
    }
}
`, scanner.DisplayName, scanner.Description, toPascalCase(scanner.Name), toPascalCase(scanner.Name))
}

// truncateString truncates a string to a maximum length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func init() {
	scannerCmd.Flags().BoolVar(&createFlag, "create", false, "Create a new scanner interactively")
	
	// Handle --create flag
	scannerCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if createFlag {
			scannerCreateCmd.Run(cmd, args)
		}
	}
	
	accessAnalyzerCmd.AddCommand(scannerCmd)
}
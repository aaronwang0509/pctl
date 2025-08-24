package main

import (
	"fmt"
	"log"

	"github.com/aaronwang/pctl/internal/token"
	pkgtoken "github.com/aaronwang/pctl/pkg/token"
)

// ExampleInternalTokenUsage demonstrates how other PCTL commands would use token generation internally
func ExampleInternalTokenUsage() {
	fmt.Println("=== PCTL Internal Token API Usage Example ===\n")
	
	// 1. Load token configuration (as ELK command would do)
	fmt.Println("1. Loading token configuration from file...")
	tokenConfig, err := pkgtoken.LoadConfig("configs/token/examples/service-account.yaml")
	if err != nil {
		log.Printf("Failed to load config: %v\n", err)
		return
	}
	fmt.Printf("✓ Loaded config for service account: %s\n", tokenConfig.ServiceAccountID)
	
	// 2. Create token client options
	fmt.Println("\n2. Creating token client with options...")
	options := pkgtoken.GeneratorOptions{
		Config:       *tokenConfig,
		OutputFormat: pkgtoken.OutputFormatJSON,
		Verbose:      false, // ELK would typically run quietly
	}
	
	_ = pkgtoken.NewClient(options)
	fmt.Printf("✓ Token client created\n")
	
	// 3. Validate configuration before use
	fmt.Println("\n3. Validating configuration...")
	if err := pkgtoken.Validate(tokenConfig); err != nil {
		log.Printf("⚠ Config validation failed: %v\n", err)
		fmt.Println("  (This is expected with example config - real config would pass)")
	} else {
		fmt.Printf("✓ Configuration is valid\n")
	}
	
	// 4. Demonstrate creating config programmatically (as ELK might do)
	fmt.Println("\n4. Creating configuration programmatically...")
	elkConfig := pkgtoken.DefaultConfig()
	elkConfig.ServiceAccountID = "elk-log-streamer-service-account"
	elkConfig.Platform = "https://openam-elk.forgerock.io"
	elkConfig.Scope = "fr:am:* fr:idm:*"
	elkConfig.ExpSeconds = 3600 // 1 hour for ELK operations
	
	fmt.Printf("✓ Created ELK-specific config:\n")
	fmt.Printf("  - Service Account: %s\n", elkConfig.ServiceAccountID)
	fmt.Printf("  - Platform: %s\n", elkConfig.Platform)
	fmt.Printf("  - Scope: %s\n", elkConfig.Scope)
	fmt.Printf("  - Expires: %d seconds\n", elkConfig.ExpSeconds)
	
	// 5. Demonstrate different output formats for internal use
	fmt.Println("\n5. Testing different output formats...")
	
	// Create a mock token result (what a real token would look like)
	mockResult := &token.TokenResult{
		AccessToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9...[truncated-for-example]",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		Scope:      "fr:am:* fr:idm:*",
		Metadata: map[string]interface{}{
			"service_account_id": elkConfig.ServiceAccountID,
			"platform":          elkConfig.Platform,
			"source":            "pctl-elk-command",
		},
	}
	
	// Test JSON format (most common for internal API usage)
	options.OutputFormat = pkgtoken.OutputFormatJSON
	jsonClient := pkgtoken.NewClient(options)
	jsonOutput, err := jsonClient.FormatOutput(mockResult)
	if err != nil {
		log.Printf("Failed to format JSON: %v", err)
	} else {
		fmt.Printf("✓ JSON format ready for API consumption\n")
		fmt.Printf("  Length: %d characters\n", len(jsonOutput))
	}
	
	// Test YAML format (useful for configuration files)
	options.OutputFormat = pkgtoken.OutputFormatYAML
	yamlClient := pkgtoken.NewClient(options)
	yamlOutput, err := yamlClient.FormatOutput(mockResult)
	if err != nil {
		log.Printf("Failed to format YAML: %v", err)
	} else {
		fmt.Printf("✓ YAML format ready for config files\n")
		fmt.Printf("  Length: %d characters\n", len(yamlOutput))
	}
	
	// 6. Demonstrate error handling for internal use
	fmt.Println("\n6. Testing error handling...")
	
	invalidConfig := &token.TokenConfig{
		Type: token.TokenTypeServiceAccount,
		// Missing required fields
	}
	
	if err := pkgtoken.Validate(invalidConfig); err != nil {
		fmt.Printf("✓ Validation correctly caught error: %s\n", err.Error())
	}
	
	// 7. Demonstrate metadata access (useful for logging and monitoring)
	fmt.Println("\n7. Accessing token metadata...")
	
	serviceAccountID := mockResult.Metadata["service_account_id"].(string)
	platform := mockResult.Metadata["platform"].(string)
	source := mockResult.Metadata["source"].(string)
	
	fmt.Printf("✓ Metadata extraction:\n")
	fmt.Printf("  - Service Account: %s\n", serviceAccountID)
	fmt.Printf("  - Platform: %s\n", platform)
	fmt.Printf("  - Source: %s\n", source)
	fmt.Printf("  - Token Length: %d chars\n", len(mockResult.AccessToken))
	fmt.Printf("  - Expires In: %d seconds\n", mockResult.ExpiresIn)
	
	fmt.Println("\n=== Internal Token API Usage Complete ===")
	fmt.Println("This demonstrates how the ELK command would use token generation internally.")
}

func runExample() {
	ExampleInternalTokenUsage()
}
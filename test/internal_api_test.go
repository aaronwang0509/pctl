package test

import (
	"testing"
	"time"

	"github.com/aaronwang/pctl/internal/token"
	pkgtoken "github.com/aaronwang/pctl/pkg/token"
)

// TestInternalTokenAPI tests using the token functionality as an internal API
// This demonstrates how other parts of PCTL (like ELK) would use token generation
func TestInternalTokenAPI(t *testing.T) {
	// Test loading config programmatically (as ELK would do)
	config := &token.TokenConfig{
		Type:             token.TokenTypeServiceAccount,
		ServiceAccountID: "internal-test-id",
		JWKJson:         `{"kty":"RSA","n":"test","e":"AQAB","d":"test"}`,
		Platform:        "https://internal.test.com",
		Scope:           "fr:am:* fr:idm:*",
		ExpSeconds:      3600,
		Verbose:         false,
	}

	// Validate config
	if err := pkgtoken.Validate(config); err != nil {
		// This error is expected since we're using test JWK
		t.Logf("Config validation failed as expected: %v", err)
	}

	// Test creating options for internal use
	options := pkgtoken.GeneratorOptions{
		Config:       *config,
		OutputFormat: pkgtoken.OutputFormatJSON,
		Verbose:      false,
	}

	client := pkgtoken.NewClient(options)
	if client == nil {
		t.Fatal("Failed to create token client")
	}

	// Test formatting - this would be used internally to get different formats
	mockResult := &token.TokenResult{
		AccessToken: "internal-test-token-12345",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		ExpiresAt:   time.Now().Add(time.Hour),
		Scope:      "fr:am:* fr:idm:*",
		Metadata: map[string]interface{}{
			"service_account_id": "internal-test-id",
			"generated_at":       time.Now().Unix(),
			"source":            "internal-api",
		},
	}

	// Test different output formats for internal use
	formats := []pkgtoken.OutputFormat{
		pkgtoken.OutputFormatText,
		pkgtoken.OutputFormatJSON,
		pkgtoken.OutputFormatYAML,
	}

	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			options.OutputFormat = format
			testClient := pkgtoken.NewClient(options)
			
			output, err := testClient.FormatOutput(mockResult)
			if err != nil {
				t.Errorf("Failed to format output as %s: %v", format, err)
			}
			
			if output == "" {
				t.Errorf("Empty output for format %s", format)
			}
			
			// Verify the output contains the token
			if !containsString(output, "internal-test-token-12345") {
				t.Errorf("Output doesn't contain expected token for format %s", format)
			}
		})
	}
}

// TestInternalConfigCreation tests creating configs programmatically
func TestInternalConfigCreation(t *testing.T) {
	// Test creating default config and modifying it (as ELK might do)
	config := pkgtoken.DefaultConfig()
	
	// Modify for internal use
	config.ServiceAccountID = "elk-service-account"
	config.Platform = "https://elk.test.com"
	config.Scope = "fr:am:* fr:idm:*"
	config.ExpSeconds = 3600
	
	if config.ServiceAccountID != "elk-service-account" {
		t.Errorf("Expected service account ID 'elk-service-account', got %s", config.ServiceAccountID)
	}
	
	if config.Platform != "https://elk.test.com" {
		t.Errorf("Expected platform 'https://elk.test.com', got %s", config.Platform)
	}
}

// TestInternalErrorHandling tests error scenarios for internal usage
func TestInternalErrorHandling(t *testing.T) {
	// Test handling validation errors internally
	invalidConfig := &token.TokenConfig{
		Type: token.TokenTypeServiceAccount,
		// Missing required fields
	}

	err := pkgtoken.Validate(invalidConfig)
	if err == nil {
		t.Error("Expected validation error for incomplete config")
	}

	// Test handling generation errors
	options := pkgtoken.GeneratorOptions{
		Config:       *invalidConfig,
		OutputFormat: pkgtoken.OutputFormatJSON,
		Verbose:      false,
	}

	client := pkgtoken.NewClient(options)
	_, err = client.Generate()
	
	if err == nil {
		t.Error("Expected generation error for invalid config")
	}

	// Verify error message is meaningful for internal consumers
	if !containsString(err.Error(), "configuration validation failed") {
		t.Errorf("Expected meaningful error message, got: %s", err.Error())
	}
}

// TestTokenMetadata tests accessing token metadata for internal use
func TestTokenMetadata(t *testing.T) {
	result := &token.TokenResult{
		AccessToken: "test-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		ExpiresAt:   time.Now().Add(time.Hour),
		Metadata: map[string]interface{}{
			"service_account_id": "test-id",
			"platform":          "https://test.com",
			"generated_at":       time.Now().Unix(),
			"custom_field":       "custom_value",
		},
	}

	// Test accessing metadata fields (as ELK might need)
	serviceAccountID, ok := result.Metadata["service_account_id"].(string)
	if !ok || serviceAccountID != "test-id" {
		t.Errorf("Expected service_account_id 'test-id', got %v", serviceAccountID)
	}

	platform, ok := result.Metadata["platform"].(string)  
	if !ok || platform != "https://test.com" {
		t.Errorf("Expected platform 'https://test.com', got %v", platform)
	}

	// Test token expiration checking (useful for token refresh logic)
	if result.ExpiresAt.Before(time.Now()) {
		t.Error("Token appears to be expired")
	}

	// Test token length validation (real tokens should be substantial)
	if len(result.AccessToken) < 10 {
		t.Errorf("Token seems too short for production use: %d chars", len(result.AccessToken))
	}
}

// Helper function for string contains check
func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
package token

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aaronwang/pctl/internal/token"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name      string
		yamlContent string
		wantErr   bool
		validate  func(t *testing.T, config *token.TokenConfig)
	}{
		{
			name: "valid service account config",
			yamlContent: `
type: "service-account"
service_account_id: "test-id"
jwk_json: '{"kty":"RSA","n":"test","e":"AQAB","d":"test"}'
platform: "https://test.forgerock.com"
scope: "fr:am:* fr:idm:*"
exp_seconds: 3600
verify_ssl: true
`,
			wantErr: false,
			validate: func(t *testing.T, config *token.TokenConfig) {
				if config.Type != token.TokenTypeServiceAccount {
					t.Errorf("Expected type %s, got %s", token.TokenTypeServiceAccount, config.Type)
				}
				if config.ServiceAccountID != "test-id" {
					t.Errorf("Expected service_account_id 'test-id', got %s", config.ServiceAccountID)
				}
				if config.Platform != "https://test.forgerock.com" {
					t.Errorf("Expected platform 'https://test.forgerock.com', got %s", config.Platform)
				}
				if config.ExpSeconds != 3600 {
					t.Errorf("Expected exp_seconds 3600, got %d", config.ExpSeconds)
				}
				if len(config.Scopes) != 2 {
					t.Errorf("Expected 2 scopes, got %d", len(config.Scopes))
				}
			},
		},
		{
			name: "config with platform field",
			yamlContent: `
service_account_id: "test-id"
jwk_json: '{"kty":"RSA"}'
platform: "https://platform.forgerock.com"
exp_seconds: 899
`,
			wantErr: false,
			validate: func(t *testing.T, config *token.TokenConfig) {
				if config.BaseURL != "https://platform.forgerock.com" {
					t.Errorf("Expected baseURL to be set from platform, got %s", config.BaseURL)
				}
			},
		},
		{
			name: "invalid yaml",
			yamlContent: `
invalid: yaml: content:
  - malformed
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary config file
			tempDir := t.TempDir()
			configPath := filepath.Join(tempDir, "test-config.yaml")
			
			if err := os.WriteFile(configPath, []byte(tt.yamlContent), 0644); err != nil {
				t.Fatalf("Failed to create temp config file: %v", err)
			}

			// Test LoadConfig
			config, err := LoadConfig(configPath)
			
			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, config)
			}
		})
	}
}

func TestLoadConfigErrors(t *testing.T) {
	// Test empty path
	_, err := LoadConfig("")
	if err == nil {
		t.Error("Expected error for empty config path")
	}

	// Test non-existent file
	_, err = LoadConfig("/non/existent/path.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *token.TokenConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid service account config",
			config: &token.TokenConfig{
				Type:             token.TokenTypeServiceAccount,
				ServiceAccountID: "test-id",
				JWKJson:         `{"kty":"RSA"}`,
				Platform:        "https://test.forgerock.com",
			},
			wantErr: false,
		},
		{
			name: "missing service account ID",
			config: &token.TokenConfig{
				Type:     token.TokenTypeServiceAccount,
				JWKJson: `{"kty":"RSA"}`,
				Platform: "https://test.forgerock.com",
			},
			wantErr: true,
			errMsg:  "service_account_id is required",
		},
		{
			name: "missing JWK",
			config: &token.TokenConfig{
				Type:             token.TokenTypeServiceAccount,
				ServiceAccountID: "test-id",
				Platform:        "https://test.forgerock.com",
			},
			wantErr: true,
			errMsg:  "jwk_json or privateKey is required",
		},
		{
			name: "missing platform",
			config: &token.TokenConfig{
				Type:             token.TokenTypeServiceAccount,
				ServiceAccountID: "test-id",
				JWKJson:         `{"kty":"RSA"}`,
			},
			wantErr: true,
			errMsg:  "baseUrl or platform is required",
		},
		{
			name: "valid user config",
			config: &token.TokenConfig{
				Type:     token.TokenTypeUser,
				Username: "testuser",
				Password: "testpass",
				Platform: "https://test.forgerock.com",
			},
			wantErr: false,
		},
		{
			name: "user config missing username",
			config: &token.TokenConfig{
				Type:     token.TokenTypeUser,
				Password: "testpass",
				Platform: "https://test.forgerock.com",
			},
			wantErr: true,
			errMsg:  "username is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.config)
			
			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			} else if tt.wantErr && err != nil {
				if tt.errMsg != "" && !containsString(err.Error(), tt.errMsg) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.errMsg, err.Error())
				}
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	if config.Type != token.TokenTypeServiceAccount {
		t.Errorf("Expected default type %s, got %s", token.TokenTypeServiceAccount, config.Type)
	}
	
	if config.ExpiresIn != 60*time.Minute {
		t.Errorf("Expected default ExpiresIn 60m, got %v", config.ExpiresIn)
	}
	
	if len(config.Scopes) != 2 {
		t.Errorf("Expected default scopes length 2, got %d", len(config.Scopes))
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
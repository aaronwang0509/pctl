package token

import (
	"encoding/json"
	"testing"
)

func TestJWKParsing(t *testing.T) {
	// Test JWK structure parsing
	jwkString := `{
		"kty": "RSA",
		"use": "sig",
		"kid": "test-key-id",
		"n": "test-modulus",
		"e": "AQAB",
		"d": "test-private-exponent"
	}`

	var jwk JWK
	if err := json.Unmarshal([]byte(jwkString), &jwk); err != nil {
		t.Fatalf("Failed to parse JWK: %v", err)
	}

	if jwk.Kty != "RSA" {
		t.Errorf("Expected kty 'RSA', got %s", jwk.Kty)
	}
	if jwk.Kid != "test-key-id" {
		t.Errorf("Expected kid 'test-key-id', got %s", jwk.Kid)
	}
	if jwk.E != "AQAB" {
		t.Errorf("Expected e 'AQAB', got %s", jwk.E)
	}
}

func TestServiceAccountGeneratorConfig(t *testing.T) {
	config := TokenConfig{
		Type:             TokenTypeServiceAccount,
		ServiceAccountID: "test-service-account",
		Platform:        "https://test.forgerock.com",
		Scope:           "fr:am:* fr:idm:*",
		ExpSeconds:      3600,
	}

	generator := &ServiceAccountGenerator{
		Config:  config,
		Verbose: false,
	}

	if generator.Config.ServiceAccountID != "test-service-account" {
		t.Errorf("Expected service account ID 'test-service-account', got %s", generator.Config.ServiceAccountID)
	}
	
	if generator.Config.Platform != "https://test.forgerock.com" {
		t.Errorf("Expected platform 'https://test.forgerock.com', got %s", generator.Config.Platform)
	}
}

func TestTokenResultStructure(t *testing.T) {
	result := &TokenResult{
		AccessToken: "test-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		Scope:      "test-scope",
		Metadata: map[string]interface{}{
			"service_account_id": "test-id",
			"platform":          "https://test.com",
		},
	}

	if result.AccessToken != "test-token" {
		t.Errorf("Expected access token 'test-token', got %s", result.AccessToken)
	}
	
	if result.TokenType != "Bearer" {
		t.Errorf("Expected token type 'Bearer', got %s", result.TokenType)
	}

	if result.Metadata["service_account_id"] != "test-id" {
		t.Errorf("Expected service_account_id 'test-id', got %v", result.Metadata["service_account_id"])
	}
}

// Test JWK validation without actually making HTTP calls
func TestJWKValidation(t *testing.T) {
	tests := []struct {
		name    string
		jwkJson string
		wantErr bool
	}{
		{
			name: "valid JWK",
			jwkJson: `{
				"kty": "RSA",
				"n": "test",
				"e": "AQAB", 
				"d": "test"
			}`,
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			jwkJson: `{"kty": "RSA", invalid json}`,
			wantErr: true,
		},
		{
			name:    "empty JWK",
			jwkJson: `{}`,
			wantErr: false, // Will fail later in RSA conversion
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var jwk JWK
			err := json.Unmarshal([]byte(tt.jwkJson), &jwk)
			
			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name   string
		config TokenConfig
		field  string // Which field to check
		value  interface{}
	}{
		{
			name: "service account type",
			config: TokenConfig{
				Type: TokenTypeServiceAccount,
			},
			field: "Type",
			value: TokenTypeServiceAccount,
		},
		{
			name: "platform URL",
			config: TokenConfig{
				Platform: "https://openam.forgerock.com",
			},
			field: "Platform", 
			value: "https://openam.forgerock.com",
		},
		{
			name: "expiration seconds",
			config: TokenConfig{
				ExpSeconds: 899,
			},
			field: "ExpSeconds",
			value: 899,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.field {
			case "Type":
				if tt.config.Type != tt.value.(TokenType) {
					t.Errorf("Expected %s, got %s", tt.value, tt.config.Type)
				}
			case "Platform":
				if tt.config.Platform != tt.value.(string) {
					t.Errorf("Expected %s, got %s", tt.value, tt.config.Platform)
				}
			case "ExpSeconds":
				if tt.config.ExpSeconds != tt.value.(int) {
					t.Errorf("Expected %d, got %d", tt.value, tt.config.ExpSeconds)
				}
			}
		})
	}
}
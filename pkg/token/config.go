package token

import (
	"fmt"
	"os"
	"time"
	"strings"

	"gopkg.in/yaml.v3"
	"github.com/aaronwang/pctl/internal/token"
)

// LoadConfig loads token configuration from a YAML file
func LoadConfig(configPath string) (*token.TokenConfig, error) {
	if configPath == "" {
		return nil, fmt.Errorf("config path is required")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config token.TokenConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults and normalize fields
	if config.Type == "" {
		config.Type = token.TokenTypeServiceAccount
	}
	
	// Handle alternative field names from authflow format
	if config.Platform != "" && config.BaseURL == "" {
		config.BaseURL = config.Platform
	}
	
	// Convert exp_seconds to ExpiresIn duration
	if config.ExpSeconds > 0 && config.ExpiresIn == 0 {
		config.ExpiresIn = time.Duration(config.ExpSeconds) * time.Second
	}
	
	// Set default expiry if none specified
	if config.ExpiresIn == 0 {
		config.ExpiresIn = 60 * time.Minute // Default 1 hour
	}

	// Convert single scope string to scopes array
	if config.Scope != "" && len(config.Scopes) == 0 {
		config.Scopes = strings.Split(config.Scope, " ")
	}

	return &config, nil
}

// Validate validates the token configuration
func Validate(c *token.TokenConfig) error {
	if c.BaseURL == "" && c.Platform == "" {
		return fmt.Errorf("baseUrl or platform is required")
	}

	switch c.Type {
	case token.TokenTypeServiceAccount:
		if c.ServiceAccountID == "" {
			return fmt.Errorf("service_account_id is required for service account tokens")
		}
		if c.JWKJson == "" && c.PrivateKey == "" {
			return fmt.Errorf("jwk_json or privateKey is required for service account tokens")
		}
	case token.TokenTypeUser:
		if c.Username == "" {
			return fmt.Errorf("username is required for user tokens")
		}
		if c.Password == "" {
			return fmt.Errorf("password is required for user tokens")
		}
	case token.TokenTypeCustom:
		if c.ClientID == "" {
			return fmt.Errorf("clientId is required for custom tokens")
		}
		if c.ClientSecret == "" {
			return fmt.Errorf("clientSecret is required for custom tokens")
		}
	default:
		return fmt.Errorf("invalid token type: %s", c.Type)
	}

	return nil
}

// DefaultConfig returns a default token configuration
func DefaultConfig() *token.TokenConfig {
	return &token.TokenConfig{
		Type:         token.TokenTypeServiceAccount,
		ExpiresIn:    60 * time.Minute,
		Scopes:       []string{"openid", "profile"},
		CustomClaims: make(map[string]interface{}),
	}
}
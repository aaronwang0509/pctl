package token

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
	"github.com/aaronwang/pctl/internal/token"
)

// Generator is the main token generator interface
type Generator interface {
	Generate() (*token.TokenResult, error)
}

// GeneratorOptions represents options for token generation
type GeneratorOptions struct {
	Config       token.TokenConfig
	OutputFormat OutputFormat
	Verbose      bool
}

// Client is the main entry point for token operations
type Client struct {
	options GeneratorOptions
}

// NewClient creates a new token client
func NewClient(options GeneratorOptions) *Client {
	return &Client{
		options: options,
	}
}

// Generate generates a token based on the configuration
func (c *Client) Generate() (*token.TokenResult, error) {
	// Validate configuration
	if err := Validate(&c.options.Config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Create appropriate generator based on token type
	var generator Generator
	switch c.options.Config.Type {
	case token.TokenTypeServiceAccount:
		generator = &token.ServiceAccountGenerator{Config: c.options.Config, Verbose: c.options.Verbose}
	case token.TokenTypeUser:
		generator = &token.UserTokenGenerator{Config: c.options.Config, Verbose: c.options.Verbose}
	case token.TokenTypeCustom:
		generator = &token.CustomTokenGenerator{Config: c.options.Config, Verbose: c.options.Verbose}
	default:
		return nil, fmt.Errorf("unsupported token type: %s", c.options.Config.Type)
	}

	return generator.Generate()
}

// FormatOutput formats the token result according to the specified format
func (c *Client) FormatOutput(result *token.TokenResult) (string, error) {
	switch c.options.OutputFormat {
	case OutputFormatJSON:
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal JSON: %w", err)
		}
		return string(data), nil

	case OutputFormatYAML:
		data, err := yaml.Marshal(result)
		if err != nil {
			return "", fmt.Errorf("failed to marshal YAML: %w", err)
		}
		return string(data), nil

	case OutputFormatText:
		fallthrough
	default:
		var output strings.Builder
		output.WriteString("Token Generation Result:\n")
		output.WriteString("=======================\n")
		output.WriteString(fmt.Sprintf("Access Token: %s\n", result.AccessToken))
		output.WriteString(fmt.Sprintf("Token Type: %s\n", result.TokenType))
		output.WriteString(fmt.Sprintf("Expires In: %d seconds\n", result.ExpiresIn))
		output.WriteString(fmt.Sprintf("Expires At: %s\n", result.ExpiresAt.Format("2006-01-02 15:04:05 MST")))
		if result.Scope != "" {
			output.WriteString(fmt.Sprintf("Scope: %s\n", result.Scope))
		}
		if result.RefreshToken != "" {
			output.WriteString(fmt.Sprintf("Refresh Token: %s\n", result.RefreshToken))
		}
		return output.String(), nil
	}
}
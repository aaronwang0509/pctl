package token

import (
	"testing"

	"github.com/aaronwang/pctl/internal/token"
)

func TestNewClient(t *testing.T) {
	options := GeneratorOptions{
		Config: token.TokenConfig{
			Type:             token.TokenTypeServiceAccount,
			ServiceAccountID: "test-id",
		},
		OutputFormat: OutputFormatJSON,
		Verbose:      true,
	}

	client := NewClient(options)
	if client == nil {
		t.Error("Expected client to be created, got nil")
	}
	
	if client.options.OutputFormat != OutputFormatJSON {
		t.Errorf("Expected output format %s, got %s", OutputFormatJSON, client.options.OutputFormat)
	}
}

func TestFormatOutput(t *testing.T) {
	client := &Client{}
	
	result := &token.TokenResult{
		AccessToken: "test-token",
		TokenType:   "Bearer", 
		ExpiresIn:   3600,
		Scope:      "test-scope",
	}

	tests := []struct {
		name         string
		outputFormat OutputFormat
		wantContains []string
		wantErr      bool
	}{
		{
			name:         "text format",
			outputFormat: OutputFormatText,
			wantContains: []string{"Token Generation Result", "Access Token: test-token", "Token Type: Bearer"},
			wantErr:      false,
		},
		{
			name:         "json format",
			outputFormat: OutputFormatJSON,
			wantContains: []string{`"access_token": "test-token"`, `"token_type": "Bearer"`},
			wantErr:      false,
		},
		{
			name:         "yaml format", 
			outputFormat: OutputFormatYAML,
			wantContains: []string{"access_token: test-token", "token_type: Bearer"},
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.options.OutputFormat = tt.outputFormat
			
			output, err := client.FormatOutput(result)
			
			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			for _, want := range tt.wantContains {
				if !containsString(output, want) {
					t.Errorf("Expected output to contain '%s', got:\n%s", want, output)
				}
			}
		})
	}
}

func TestGenerateValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		config  token.TokenConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "missing service account ID",
			config: token.TokenConfig{
				Type:     token.TokenTypeServiceAccount,
				Platform: "https://test.com",
			},
			wantErr: true,
			errMsg:  "service_account_id is required",
		},
		{
			name: "missing platform",
			config: token.TokenConfig{
				Type:             token.TokenTypeServiceAccount,
				ServiceAccountID: "test-id",
				JWKJson:         `{"kty":"RSA"}`,
			},
			wantErr: true,
			errMsg:  "baseUrl or platform is required",
		},
		{
			name: "invalid token type",
			config: token.TokenConfig{
				Type:     "invalid-type",
				Platform: "https://test.com",
			},
			wantErr: true,
			errMsg:  "invalid token type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := GeneratorOptions{
				Config:       tt.config,
				OutputFormat: OutputFormatText,
				Verbose:      false,
			}

			client := NewClient(options)
			_, err := client.Generate()

			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			} else if tt.wantErr && err != nil {
				if !containsString(err.Error(), tt.errMsg) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.errMsg, err.Error())
				}
			}
		})
	}
}
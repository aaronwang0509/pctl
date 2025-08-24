package token

import (
	"time"
)

// TokenType represents the type of token to generate
type TokenType string

const (
	TokenTypeServiceAccount TokenType = "service-account"
	TokenTypeUser           TokenType = "user"
	TokenTypeCustom         TokenType = "custom"
)

// TokenConfig represents the configuration for token generation
type TokenConfig struct {
	// Token type
	Type TokenType `yaml:"type" json:"type"`
	
	// PAIC connection details
	BaseURL      string `yaml:"baseUrl" json:"baseUrl"`
	Platform     string `yaml:"platform" json:"platform"` // Alternative name for baseUrl
	Username     string `yaml:"username" json:"username"`
	Password     string `yaml:"password" json:"password"`
	ClientID     string `yaml:"clientId" json:"clientId"`
	ClientSecret string `yaml:"clientSecret" json:"clientSecret"`
	
	// Service Account specific
	ServiceAccountID   string `yaml:"service_account_id" json:"service_account_id"`
	ServiceAccountName string `yaml:"serviceAccountName" json:"serviceAccountName"`
	PrivateKey         string `yaml:"privateKey" json:"privateKey"`
	KeyID              string `yaml:"keyId" json:"keyId"`
	JWKJson            string `yaml:"jwk_json" json:"jwk_json"` // JWK as JSON string
	
	// Token properties
	Audience  string        `yaml:"audience" json:"audience"`
	Issuer    string        `yaml:"issuer" json:"issuer"`
	Subject   string        `yaml:"subject" json:"subject"`
	ExpiresIn time.Duration `yaml:"expiresIn" json:"expiresIn"`
	ExpSeconds int          `yaml:"exp_seconds" json:"exp_seconds"` // Alternative expiry format
	Scopes    []string      `yaml:"scopes" json:"scopes"`
	Scope     string        `yaml:"scope" json:"scope"` // Alternative single scope format
	
	// Output and behavior
	OutputFormat string `yaml:"output_format" json:"output_format"`
	Verbose      bool   `yaml:"verbose" json:"verbose"`
	VerifySSL    bool   `yaml:"verify_ssl" json:"verify_ssl"`
	Proxy        string `yaml:"proxy" json:"proxy"`
	
	// Custom claims
	CustomClaims map[string]interface{} `yaml:"customClaims" json:"customClaims"`
}

// TokenResult represents the result of token generation
type TokenResult struct {
	AccessToken  string                 `json:"access_token" yaml:"access_token"`
	TokenType    string                 `json:"token_type" yaml:"token_type"`
	ExpiresIn    int64                  `json:"expires_in" yaml:"expires_in"`
	ExpiresAt    time.Time              `json:"expires_at" yaml:"expires_at"`
	Scope        string                 `json:"scope,omitempty" yaml:"scope,omitempty"`
	RefreshToken string                 `json:"refresh_token,omitempty" yaml:"refresh_token,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}
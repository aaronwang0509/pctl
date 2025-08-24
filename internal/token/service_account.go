package token

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ServiceAccountGenerator handles service account token generation
type ServiceAccountGenerator struct {
	Config  TokenConfig
	Verbose bool
}

// JWK represents a JSON Web Key structure
type JWK struct {
	Kty string `json:"kty"` // Key Type
	Use string `json:"use"` // Public Key Use
	Kid string `json:"kid"` // Key ID
	N   string `json:"n"`   // Modulus
	E   string `json:"e"`   // Exponent
	D   string `json:"d"`   // Private Exponent
	P   string `json:"p"`   // First Prime Factor
	Q   string `json:"q"`   // Second Prime Factor
	DP  string `json:"dp"`  // First Factor CRT Exponent
	DQ  string `json:"dq"`  // Second Factor CRT Exponent
	QI  string `json:"qi"`  // First CRT Coefficient
}

// PaicTokenResponse represents the response from PAIC token endpoint
type PaicTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
	Scope       string `json:"scope,omitempty"`
}

// Generate generates a service account token
func (g *ServiceAccountGenerator) Generate() (*TokenResult, error) {
	if g.Verbose {
		fmt.Printf("Generating service account token for: %s\n", g.Config.ServiceAccountID)
	}

	// Parse JWK from JSON string
	var jwk JWK
	if err := json.Unmarshal([]byte(g.Config.JWKJson), &jwk); err != nil {
		return nil, fmt.Errorf("failed to parse JWK: %w", err)
	}

	// Create RSA private key from JWK
	privateKey, err := g.jwkToRSAPrivateKey(&jwk)
	if err != nil {
		return nil, fmt.Errorf("failed to convert JWK to RSA private key: %w", err)
	}

	// Create JWT assertion
	jwtAssertion, err := g.createJWTAssertion(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT assertion: %w", err)
	}

	if g.Verbose {
		fmt.Printf("JWT assertion created successfully\n")
	}

	// Exchange JWT assertion for access token
	tokenResponse, err := g.exchangeJWTForToken(jwtAssertion)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange JWT for token: %w", err)
	}

	// Build result
	now := time.Now()
	expiresAt := now.Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)
	
	result := &TokenResult{
		AccessToken: tokenResponse.AccessToken,
		TokenType:   tokenResponse.TokenType,
		ExpiresIn:   tokenResponse.ExpiresIn,
		ExpiresAt:   expiresAt,
		Scope:       tokenResponse.Scope,
		Metadata: map[string]interface{}{
			"service_account_id": g.Config.ServiceAccountID,
			"generated_at":       now.Unix(),
			"platform":          g.Config.Platform,
		},
	}

	if g.Verbose {
		fmt.Printf("Token generated successfully, expires at: %s\n", result.ExpiresAt.Format(time.RFC3339))
	}

	return result, nil
}

// jwkToRSAPrivateKey converts JWK to RSA private key
func (g *ServiceAccountGenerator) jwkToRSAPrivateKey(jwk *JWK) (*rsa.PrivateKey, error) {
	// Decode base64url components
	n, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %w", err)
	}
	
	d, err := base64.RawURLEncoding.DecodeString(jwk.D)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private exponent: %w", err)
	}

	p, err := base64.RawURLEncoding.DecodeString(jwk.P)
	if err != nil {
		return nil, fmt.Errorf("failed to decode first prime: %w", err)
	}

	q, err := base64.RawURLEncoding.DecodeString(jwk.Q)
	if err != nil {
		return nil, fmt.Errorf("failed to decode second prime: %w", err)
	}

	// Create big integers from byte arrays
	nInt := new(big.Int).SetBytes(n)
	dInt := new(big.Int).SetBytes(d)
	pInt := new(big.Int).SetBytes(p)
	qInt := new(big.Int).SetBytes(q)

	// Create RSA private key
	key := &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: nInt,
			E: 65537, // Standard RSA exponent (AQAB in base64)
		},
		D:      dInt,
		Primes: []*big.Int{pInt, qInt},
	}

	// Precompute values for faster operations
	key.Precompute()

	return key, nil
}

// createJWTAssertion creates a JWT assertion for service account authentication
func (g *ServiceAccountGenerator) createJWTAssertion(privateKey *rsa.PrivateKey) (string, error) {
	now := time.Now()
	
	// Generate random JWT ID
	jtiBytes := make([]byte, 16)
	if _, err := rand.Read(jtiBytes); err != nil {
		return "", fmt.Errorf("failed to generate JWT ID: %w", err)
	}
	jti := base64.RawURLEncoding.EncodeToString(jtiBytes)

	// Build audience URL
	baseURL := strings.TrimRight(g.Config.BaseURL, "/")
	if baseURL == "" {
		baseURL = strings.TrimRight(g.Config.Platform, "/")
	}
	audience := baseURL + "/am/oauth2/access_token"

	// Determine expiration
	expSeconds := g.Config.ExpSeconds
	if expSeconds == 0 {
		expSeconds = int(g.Config.ExpiresIn.Seconds())
	}
	if expSeconds == 0 {
		expSeconds = 899 // Default to 899 seconds
	}

	// Create JWT claims
	claims := jwt.MapClaims{
		"iss": g.Config.ServiceAccountID,
		"sub": g.Config.ServiceAccountID,
		"aud": audience,
		"exp": now.Unix() + int64(expSeconds),
		"jti": jti,
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign token
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	if g.Verbose {
		fmt.Printf("JWT assertion created for audience: %s\n", audience)
		fmt.Printf("JWT expiration: %s\n", time.Unix(now.Unix()+int64(expSeconds), 0).Format(time.RFC3339))
	}

	return tokenString, nil
}

// exchangeJWTForToken exchanges JWT assertion for access token
func (g *ServiceAccountGenerator) exchangeJWTForToken(jwtAssertion string) (*PaicTokenResponse, error) {
	// Build token endpoint URL
	baseURL := strings.TrimRight(g.Config.BaseURL, "/")
	if baseURL == "" {
		baseURL = strings.TrimRight(g.Config.Platform, "/")
	}
	tokenURL := baseURL + "/am/oauth2/access_token"

	// Prepare form data
	data := url.Values{
		"client_id":   {"service-account"},
		"grant_type":  {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
		"assertion":   {jwtAssertion},
		"scope":       {g.Config.Scope},
	}

	if g.Verbose {
		fmt.Printf("Making token request to: %s\n", tokenURL)
		fmt.Printf("Grant type: %s\n", "urn:ietf:params:oauth:grant-type:jwt-bearer")
		fmt.Printf("Scope: %s\n", g.Config.Scope)
	}

	// Create HTTP client
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create request
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "pctl/0.1.0")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make token request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if g.Verbose {
		fmt.Printf("Response status: %d %s\n", resp.StatusCode, resp.Status)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		if g.Verbose {
			fmt.Printf("Response body: %s\n", string(body))
		}
		return nil, fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var tokenResponse PaicTokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	if g.Verbose {
		fmt.Printf("Access token received (length: %d chars)\n", len(tokenResponse.AccessToken))
		fmt.Printf("Token type: %s\n", tokenResponse.TokenType)
		fmt.Printf("Expires in: %d seconds\n", tokenResponse.ExpiresIn)
	}

	return &tokenResponse, nil
}
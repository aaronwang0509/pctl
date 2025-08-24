package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/aaronwang/pctl/pkg/token"
)

var (
	tokenConfigFile string
	tokenOutput     string
	tokenType       string
)

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Generate and manage JWT tokens for PAIC authentication",
	Long: `Generate JWT tokens for PAIC authentication including:
- Service account tokens
- User authentication tokens
- Custom JWT tokens with specific claims

Examples:
  pctl token -c config.yaml
  pctl token --type service-account --output json
  pctl token --config token-config.yaml --verbose`,
	RunE: runToken,
}

func runToken(cmd *cobra.Command, args []string) error {
	// Load token configuration
	tokenConfig, err := token.LoadConfig(tokenConfigFile)
	if err != nil {
		return fmt.Errorf("failed to load token config: %w", err)
	}

	// Override token type from CLI flag if different  
	if tokenType != "service-account" {
		switch tokenType {
		case "user":
			tokenConfig.Type = "user"
		case "custom":
			tokenConfig.Type = "custom" 
		}
	}

	// Create token client options
	options := token.GeneratorOptions{
		Config:       *tokenConfig,
		OutputFormat: token.OutputFormat(tokenOutput),
		Verbose:      viper.GetBool("verbose"),
	}

	// Create token client and generate token
	client := token.NewClient(options)
	result, err := client.Generate()
	if err != nil {
		return fmt.Errorf("token generation failed: %w", err)
	}

	// Format and output the result
	output, err := client.FormatOutput(result)
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	fmt.Print(output)
	return nil
}

func init() {
	rootCmd.AddCommand(tokenCmd)

	// Token-specific flags
	tokenCmd.Flags().StringVarP(&tokenConfigFile, "config", "c", "", "token configuration file (required)")
	tokenCmd.Flags().StringVarP(&tokenOutput, "output", "o", "text", "output format (text, json, yaml)")
	tokenCmd.Flags().StringVarP(&tokenType, "type", "t", "service-account", "token type (service-account, user, custom)")

	// Mark config as required
	tokenCmd.MarkFlagRequired("config")

	// Bind flags to viper
	viper.BindPFlag("token.config", tokenCmd.Flags().Lookup("config"))
	viper.BindPFlag("token.output", tokenCmd.Flags().Lookup("output"))
	viper.BindPFlag("token.type", tokenCmd.Flags().Lookup("type"))
}
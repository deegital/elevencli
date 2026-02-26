package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	elevenlabs "github.com/haguro/elevenlabs-go"
	"github.com/spf13/cobra"

	"github.com/deegital/elevencli/internal/config"
)

var (
	version       = "0.1.0"
	apiKey        string
	resolvedKey   string
	client        *elevenlabs.Client
)

var rootCmd = &cobra.Command{
	Use:     "elevencli",
	Short:   "ElevenLabs CLI — text-to-speech and sound effects from the command line",
	Version: version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		config.Init()
		if cmd.Annotations["noAuth"] == "true" {
			return nil
		}
		key, err := config.ResolveAPIKey(apiKey)
		if err != nil {
			return err
		}
		resolvedKey = key
		client = elevenlabs.NewClient(context.Background(), key, 120*time.Second)
		return nil
	},
}

// resolveAPIKeyValue returns the resolved API key for direct HTTP calls.
func resolveAPIKeyValue() (string, error) {
	if resolvedKey != "" {
		return resolvedKey, nil
	}
	return "", fmt.Errorf("API key not resolved")
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "ElevenLabs API key")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

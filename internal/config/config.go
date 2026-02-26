package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Init sets up viper to read from ~/.elevencli.yaml and the ELEVENLABS_API_KEY env var.
func Init() {
	viper.SetConfigName(".elevencli")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.SetEnvPrefix("")
	_ = viper.BindEnv("api_key", "ELEVENLABS_API_KEY")
	_ = viper.ReadInConfig() // ok if missing
}

// ResolveAPIKey returns the API key using priority: flag > env > config file.
// The flagValue should be the value of the --api-key flag (empty if not set).
func ResolveAPIKey(flagValue string) (string, error) {
	if flagValue != "" {
		return flagValue, nil
	}
	if key := viper.GetString("api_key"); key != "" {
		return key, nil
	}
	return "", fmt.Errorf(`API key not found. Provide it via one of:
  1. --api-key flag:            elevencli --api-key <key> <command>
  2. Environment variable:      export ELEVENLABS_API_KEY=<key>
  3. Config file (~/.elevencli.yaml):
       api_key: <key>`)
}

package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"wax/pkg/output"
)

const defaultConfigTemplate = `# wax config
discogs:
  token: <your_discogs_personal_access_token>
  username: <your_discogs_username>
`

// DiscogsConfig holds Discogs API credentials.
type DiscogsConfig struct {
	Token    string `mapstructure:"token"`
	Username string `mapstructure:"username"`
}

// Config is the root config struct.
type Config struct {
	Discogs DiscogsConfig `mapstructure:"discogs"`
}

var AppConfig Config

// Load reads config from the default config directory.
func Load() error {
	return LoadWithDir(DefaultConfigDir())
}

// LoadWithDir reads config from the given directory, creating a default if missing.
func LoadWithDir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("cannot create config dir %s: %w", dir, err)
	}

	configFile := filepath.Join(dir, "config.yaml")

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := os.WriteFile(configFile, []byte(defaultConfigTemplate), 0644); err != nil {
			return fmt.Errorf("cannot write default config: %w", err)
		}
		output.Info("Default config created at: %s", configFile)
		output.Warning("Edit the config and set your Discogs token and username.")
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	if err := viper.MergeInConfig(); err != nil {
		return fmt.Errorf("cannot read config: %w", err)
	}

	return nil
}

// DefaultConfigDir returns ~/.config/wax.
func DefaultConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("cannot get home dir: %s", err)
	}
	return filepath.Join(home, ".config", "wax")
}

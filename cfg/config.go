package cfg

import (
	"os"
	"path/filepath"

	"github.com/gocd-contrib/gocd-cli/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	CONFIG_DIRNAME  = ".gocd"
	CONFIG_FILENAME = "settings"
	CONFIG_ENV_PFX  = "gocdcli"
)

type Config struct {
	native *viper.Viper
}

func (c *Config) SetServerUrl(url string) error {
	// to be implemented
	return nil
}

func (c *Config) GetServerUrl() string {
	// to be implemented
	return ""
}

func (c *Config) SetBasicAuth(user, pass string) error {
	// to be implemented
	return nil
}

func (c *Config) GetAuth() map[string]string {
	// to be implemented
	return nil
}

func (c *Config) Consume(configFile string) {
	if configFile != "" {
		// Use config file from the flag.
		c.native.SetConfigFile(configFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			utils.AbortLoudly(err)
		}

		cfgDir := filepath.Join(home, CONFIG_DIRNAME)
		os.MkdirAll(cfgDir, os.ModePerm)

		c.native.AddConfigPath(cfgDir)
		c.native.SetConfigName(CONFIG_FILENAME)
	}

	// If a config file is found, read it in.
	if err := c.native.ReadInConfig(); err == nil {
		utils.Echofln("Using config file:", c.native.ConfigFileUsed())
	}
}

func NewConfig() *Config {
	native := viper.New()
	native.SetEnvPrefix(CONFIG_ENV_PFX)
	native.AutomaticEnv() // read in environment variables that match

	return &Config{native: native}
}

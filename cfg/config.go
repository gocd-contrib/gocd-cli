package cfg

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gocd-contrib/gocd-cli/utils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	CONFIG_DIRNAME  = ".gocd"
	CONFIG_FILENAME = "settings"
	CONFIG_FILETYPE = "yaml"
	CONFIG_ENV_PFX  = "gocdcli"
)

type Config struct {
	native *viper.Viper
}

var onlyNumeric, _ = regexp.Compile(`^\d+$`)

func (c *Config) SetServerUrl(urlArg string) error {
	if "" == urlArg {
		return errors.New("Must specify a url")
	}

	if u, err := url.Parse(urlArg); err != nil {
		return err
	} else {
		if !u.IsAbs() || u.Hostname() == "" {
			return errors.New("URL must include protocol and hostname")
		}

		if u.Port() != "" && !onlyNumeric.MatchString(u.Port()) {
			return errors.New("Port must be numeric")
		}

		c.native.Set("server.url", u.String())
		return c.native.WriteConfig()
	}
}

func (c *Config) GetServerUrl() string {
	return c.native.GetString("server.url")
}

func (c *Config) SetBasicAuth(user, pass string) error {
	if "" == user || "" == pass {
		return errors.New("Must specify user and password")
	}

	c.native.Set("auth.type", "basic")
	c.native.Set("auth.user", user)
	c.native.Set("auth.password", pass)
	return c.native.WriteConfig()
}

func (c *Config) GetAuth() map[string]string {
	return c.native.GetStringMapString("auth")
}

func (c *Config) Consume(configFile string) error {
	if configFile != "" {
		// Use config file from the flag.
		c.native.SetConfigFile(configFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		cfgDir := filepath.Join(home, CONFIG_DIRNAME)

		if err = os.MkdirAll(cfgDir, os.ModePerm); err != nil {
			return err
		}

		c.native.AddConfigPath(cfgDir)
		c.native.SetConfigName(CONFIG_FILENAME)
		c.native.SetConfigType(CONFIG_FILETYPE)
		configFile = filepath.Join(cfgDir, CONFIG_FILENAME+"."+CONFIG_FILETYPE)
	}

	// If a config file is found, read it in.
	if err := c.native.ReadInConfig(); err == nil {
		utils.Echofln("Using config file: %s", c.native.ConfigFileUsed())
		return c.native.WriteConfigAs(configFile)
	} else {
		return err
	}
}

func NewConfig() *Config {
	native := viper.New()
	native.SetEnvPrefix(CONFIG_ENV_PFX)
	native.AutomaticEnv() // read in environment variables that match

	return &Config{native: native}
}

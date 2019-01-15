package cfg

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gocd-contrib/gocd-cli/utils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

const (
	CONFIG_DIRNAME  = `.gocd`
	CONFIG_FILENAME = `settings`
	CONFIG_FILETYPE = `yaml`
	CONFIG_ENV_PFX  = `gocdcli`
	CONFIG_VERSION  = `config_version`
	CURRENT_VERSION = 1 // format version
)

type dict map[string]interface{}

type Config struct {
	native *viper.Viper
	fs     afero.Fs
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

func (c *Config) ConfigFile() string {
	return c.native.ConfigFileUsed()
}

func (c *Config) Bootstrap(configFile string, migrations []*migration) (err error) {
	if err = c.Consume(configFile); err == nil {
		if err = c.Migrate(migrations); err == nil {
			c.LayerConfigs()
		}
	}
	return
}

func (c *Config) Consume(configFile string) error {
	if configFile != "" {
		// Use config file from the flag.
		c.native.SetConfigFile(configFile)
		return c.native.ReadInConfig()
	}
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	cfgDir := filepath.Join(home, CONFIG_DIRNAME)

	if err = c.fs.MkdirAll(cfgDir, os.ModePerm); err != nil {
		return err
	}

	c.native.AddConfigPath(cfgDir)
	c.native.SetConfigName(CONFIG_FILENAME)
	configFile = filepath.Join(cfgDir, CONFIG_FILENAME+"."+CONFIG_FILETYPE)

	// If a config file is found, read it in.
	if err := c.native.ReadInConfig(); err == nil {
		return nil
	} else {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			utils.Debug(`No config file found; creating a new one at: %q`, configFile)
			return c.native.WriteConfigAs(configFile)
		default:
			return err
		}
	}
}

func (c *Config) Migrate(migrations []*migration) error {
	utils.Debug(`Running any necessary config migrations...`)

	if !c.native.InConfig(CONFIG_VERSION) {
		utils.Debug(`%q is missing in config file; assuming current version`, CONFIG_VERSION)
		c.ensureCurrentVersion()
		if err := c.native.WriteConfig(); err != nil {
			return err
		}
	}

	if _, ok := c.native.Get(CONFIG_VERSION).(int); !ok {
		return fmt.Errorf(`The %q key in %q must be numeric`, CONFIG_VERSION, c.ConfigFile())
	}

	if ver := c.native.GetInt(CONFIG_VERSION); ver > CURRENT_VERSION {
		return fmt.Errorf(`%q: %d is not supported by the configuration file in this CLI version; max supported version: %d`, CONFIG_VERSION, ver, CURRENT_VERSION)
	} else {
		if ver != CURRENT_VERSION {
			utils.Debug(`Config file version is out of date; migrating...`)

			if migratedConf, err := applyMigrations(c.native.AllSettings(), migrations); err == nil {
				migrated := newViper(c.fs)
				migrated.SetConfigFile(c.ConfigFile())

				utils.Debug(`Committing changes...`)
				if err = migrated.MergeConfigMap(migratedConf); err != nil {
					return err
				}

				utils.Debug(`Flushing updated config to disk...`)
				if err = migrated.WriteConfig(); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}

func (c *Config) Unset(key string) error {
	if "" == key {
		return errors.New(`Missing key`)
	}

	switch key {
	case `auth-basic`: // we only support a single user profile, so make this an alias
		fallthrough
	case `auth`:
		return c.writeConfigExcludingKey(`auth`)
	case `server-url`:
		return c.writeConfigExcludingKey(`server.url`)
	default:
		return fmt.Errorf(`Unknown key %q`, key)
	}
}

/**
 * This effectively removes a key from a config file, and accepts a
 * nested key syntax (e.g., foo.bar.baz).
 *
 * Viper does not have a facility to delete a key from the config file.
 * One might think we could get away with some obvious alternatives, but
 * they fail horribly:
 *
 *   `viper.Set(key, nil)` clears the override register, but the config
 *      register values still "show through".
 *
 *   `viper.Set(key, "")` sets values to empty string, and is thus not a
 *      generic solution as one needs to type-check Get() for string values,
 *      even if the value you expect is not a string. Also, it leaves the
 *      config file riddled with empty-string literals, which is especially
 *      ugly on nested keys.
 *
 *   `viper.MergeConfigMap(map[string]interface{} {"key": nil})` panics
 *      on nil pointer for nested keys.
 *
 * So, as luck would have it, this seemingly simple operation turned out
 * to be much more complicated than anticipated, as Viper has not yet built
 * this into its core functionality. Quelle suprise!
 *
 * https://giphy.com/gifs/askreddit-surprise-FbfNWx3LPoy2I
 */
func (c *Config) writeConfigExcludingKey(key string) error {
	// `_readonly` is a viper instance that is only for reading from
	// disk, ignoring any override registers, defaults, and environment
	// variable overrides. We only want on-disk values.
	var _readonly *viper.Viper
	if ro, err := c.fsOnlyViper(c.fs); err != nil {
		return err
	} else {
		_readonly = ro
	}

	cfg := make(dict)

	// add all keys+vals from the config file except for the specified key
	for _, k := range _readonly.AllKeys() {
		if k != key && !strings.HasPrefix(k, key+`.`) {
			if recursiveHasKey(_readonly, k) {
				// construct map of on-disk values, excluding specified key
				cfg[k] = _readonly.Get(k)
			}
		} else {
			// Clear viper's override register for the specified key to
			// be deleted.
			//
			// Note that this will not clear the config register, which
			// represents what is on disk. This will also *NOT* clobber
			// any environment variable overrides, which is the correct
			// and intended behavior.
			c.native.Set(k, nil)
		}
	}

	// `_tmp` is a viper instance that is only for writing to disk.
	//
	// It is important to not swap out `c.native` for this instance
	// as `c.native` may contain entries in its override register
	// that would otherwise be destroyed.
	_tmp := newViper(c.fs)

	// writes to config register, and not the override register
	if err := _tmp.MergeConfigMap(cfg); err != nil {
		return err
	}

	if err := _tmp.WriteConfigAs(c.ConfigFile()); err != nil {
		return err
	}

	if f, err := c.fs.Open(c.ConfigFile()); err != nil {
		return err
	} else {
		return c.native.ReadConfig(f)
	}
}

/**
 * Returns a new Viper with values only sourced from the config file
 *
 * Intentionally excludes any defaults and environment variable overrides
 * as we will use this to read flattened keys and values declared only
 * on disk.
 */
func (c *Config) fsOnlyViper(fs afero.Fs) (*viper.Viper, error) {
	// clean viper with no defaults and AutomaticeEnv() unset
	v := newViper(fs)
	v.SetConfigFile(c.ConfigFile())
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	return v, nil
}

/**
 * Viper currently does not have a way to test if nested keys are
 * present in the config file. Viper.IsSet() tests against config
 * after merging with all sources (e.g., remote, environment, defaults,
 * flags, etc), which is not what we want. This basically applies
 * Viper.InConfig() over the nested path segments.
 */
func recursiveHasKey(v *viper.Viper, key string) (found bool) {
	key = strings.ToLower(key)
	paths := strings.Split(key, `.`)

	if len(paths) > 1 {
		for _, seg := range paths {
			if found = v.InConfig(seg); found {
				v = v.Sub(seg)
			} else {
				return
			}
		}
	} else {
		found = v.InConfig(key)
	}
	return
}

func (c *Config) ensureCurrentVersion() {
	c.native.MergeConfigMap(map[string]interface{}{
		CONFIG_VERSION: CURRENT_VERSION,
	})
}

func newViper(fs afero.Fs) *viper.Viper {
	v := viper.New()
	v.SetConfigType(CONFIG_FILETYPE)
	v.SetFs(fs)
	return v
}

// Configures any default values and enables environment variable
// overrides
func (c *Config) LayerConfigs() {
	v := c.native
	v.SetDefault(CONFIG_VERSION, CURRENT_VERSION)
	v.SetEnvPrefix(CONFIG_ENV_PFX)
	v.AutomaticEnv() // read in environment variables that match
}

func NewConfig(fs afero.Fs) *Config {
	return &Config{native: newViper(fs), fs: fs}
}

package cfg

import "github.com/spf13/afero"

var conf = NewConfig(afero.NewOsFs())

func Conf() *Config {
	return conf
}

func Setup(configFile string) (err error) {
	return Conf().Bootstrap(configFile, migrations)
}

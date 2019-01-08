package cfg

var conf = NewConfig()

func Conf() *Config {
	return conf
}

func Setup(configFile string) (err error) {
	return Conf().Bootstrap(configFile, migrations)
}

package cfg

var CfgFile string
var conf = NewConfig()

func Conf() *Config {
	return conf
}

func Setup() error {
	return Conf().Consume(CfgFile)
}

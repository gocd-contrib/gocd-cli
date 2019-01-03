package cfg

var CfgFile string
var conf = NewConfig()

func Conf() *Config {
	return conf
}

func Setup() {
	Conf().Consume(CfgFile)
}

package config

type Config struct {
	Logger     LoggerConf
	Connection DatabaseConf
	Store      string
}

type LoggerConf struct {
	Level string
	File  string
}

type DatabaseConf struct {
	Host     string
	User     string
	Password string
}

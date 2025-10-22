package model

type DBConfig struct {
	User           string `mapstructure:"user"`
	Passwd         string `mapstructure:"password"`
	Net            string `mapstructure:"net"`
	Addr           string `mapstructure:"addr"`
	DBName         string `mapstructure:"dbname"`
	MaxOpenConns   int    `mapstructure:"max_open_conns"`
	MaxIdleConns   int    `mapstructure:"max_idle_conns"`
}
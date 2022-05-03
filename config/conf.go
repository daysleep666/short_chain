package config

import (
	"github.com/BurntSushi/toml"
)

type ReisConfig struct {
	Addr string `toml:"addr"`
	Pwd  string `toml:"pwd"`
}

type MysqlConfig struct {
	Addr   string `toml:"addr"`
	User   string `toml:"user"`
	Pwd    string `toml:"pwd"`
	DBName string `toml:"dbname"`
}

type LoggerConfig struct {
	FilePath string `toml:"file_path"`
}

type ServerConfig struct {
	LoggerConfig LoggerConfig `toml:"logger"`
	RedisConfig  ReisConfig   `toml:"redis"`
	MysqlConfig  MysqlConfig  `toml:"short_url"`
}

var CONFIG_INSTANCE ServerConfig

func InitConfig(filePath string) (err error) {
	if _, err = toml.DecodeFile(filePath, &CONFIG_INSTANCE); err != nil {
		return
	}
	return
}

package config

import (
	"github.com/BurntSushi/toml"
)

type ReisConfig struct {
	Addr string `toml:"addr"`
	Pwd  string `toml:"pwd"`
}

type MysqlConfig struct {
	Addr     string `toml:"addr"`
	User     string `toml:"user"`
	Pwd      string `toml:"pwd"`
	DBName   string `toml:"dbname"`
	TableCnt int16  `toml:"table_cnt"`
}

type LoggerConfig struct {
	FilePath string `toml:"file_path"`
	Level    uint8  `toml:"level"`
}

type ServerConfig struct {
	Domain          string `toml:"domain"`
	MachineID       int64  `toml:"machine_id"`
	LoggerMachineID int64  `toml:"logger_machine_id"`
}

type Config struct {
	LoggerConfig        LoggerConfig `toml:"logger"`
	RedisConfig         ReisConfig   `toml:"redis"`
	ShortURLMysqlConfig MysqlConfig  `toml:"short_url"`
	ServerConfig        ServerConfig `toml:"server"`
}

var CONFIG_INSTANCE Config

func InitConfig(filePath string) (err error) {
	if _, err = toml.DecodeFile(filePath, &CONFIG_INSTANCE); err != nil {
		return
	}
	return
}

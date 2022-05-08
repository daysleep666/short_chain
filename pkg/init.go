package pkg

import (
	"log"

	"github.com/daysleep666/short_chain/config"
	"github.com/daysleep666/short_chain/pkg/repo"
)

func MustInit(configPath string) {
	if err := config.InitConfig(configPath); err != nil {
		log.Fatalf("init config failed err:%v", err)
	}
	if err := repo.InitShortUrlRecordDB(config.CONFIG_INSTANCE.ShortURLMysqlConfig.Addr, config.CONFIG_INSTANCE.ShortURLMysqlConfig.User, config.CONFIG_INSTANCE.ShortURLMysqlConfig.Pwd, config.CONFIG_INSTANCE.ShortURLMysqlConfig.DBName); err != nil {
		log.Fatalf("init short_url_record_db failed err:%v", err)
	}
	if err := repo.InitRedis(config.CONFIG_INSTANCE.RedisConfig.Addr, config.CONFIG_INSTANCE.RedisConfig.Pwd); err != nil {
		log.Fatalf("init redis failed err:%v", err)
	}
	if err := InitLogger(config.CONFIG_INSTANCE.LoggerConfig.FilePath); err != nil {
		log.Fatalf("init redis failed err:%v", err)
	}

	InitUniqueIDSerInstance(config.CONFIG_INSTANCE.ServerConfig.MachineID)
	InitUniqueIDSerForLoggerInstance(config.CONFIG_INSTANCE.ServerConfig.LoggerMachineID)
}

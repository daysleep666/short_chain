package pkg

import (
	"context"

	"github.com/daysleep666/short_chain/config"
	"github.com/daysleep666/short_chain/pkg/repo"
)

const (
	UNIQUE_ID_KEY = "unique_id"
)

type UniqueIDService interface {
	Generate(ctx context.Context) (uint64, error)
}

var UNIQUE_ID_SERVICE_INSTANCE UniqueIDService

type uniqueIDServiceRedis struct {
	logger Logger
}

func NewUniqueIDService(logger Logger) (s UniqueIDService) {
	return &uniqueIDServiceRedis{
		logger: logger,
	}
}

func (r *uniqueIDServiceRedis) Generate(ctx context.Context) (uint64, error) {
	id, err := repo.REDIS_INSTANCE.Incr(ctx, UNIQUE_ID_KEY).Uint64()
	if err != nil {
		r.logger.Errorf("[generate id failed] [err:%v]", err)
		return 0, config.DB_ERROR
	}
	return id, nil
}

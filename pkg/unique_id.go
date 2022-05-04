package pkg

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/daysleep666/short_chain/config"
	"github.com/daysleep666/short_chain/pkg/repo"

	"go.uber.org/atomic"
)

const (
	UNIQUE_ID_KEY = "unique_id"
	// 基准时间戳 2022-01-01 00:00:00
	INITIAL_MS = 1640966400000
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

type uniqueIDServiceSnowflake struct {
	machineID int64
	number    atomic.Int64
}

func NewUniqueIDSnowflakeService(machineID int64) (s UniqueIDService) {
	return &uniqueIDServiceSnowflake{
		machineID: machineID,
	}
}

func (r *uniqueIDServiceSnowflake) Generate(ctx context.Context) (uint64, error) {
	str := fmt.Sprintf("0%041b%010b%012b",
		(time.Now().UnixMilli()-INITIAL_MS)&2199023255551,
		r.machineID,
		r.number.Inc()&4095)
	uniqueID, _ := strconv.ParseUint(str, 2, 64)
	return uniqueID, nil
}

package pkg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/daysleep666/short_chain/config"
	"github.com/daysleep666/short_chain/pkg/repo"
)

type ShortURLStorageSaveParam struct {
	LongURL  string
	ShortURL string
	UniqueID uint64
}

func (p *ShortURLStorageSaveParam) Check() error {
	if len(p.LongURL) == 0 || len(p.ShortURL) == 0 || p.UniqueID == 0 {
		return config.PARAM_ERROR
	}
	return nil
}

type ShortURLStorageService interface {
	Save(ctx context.Context, param *ShortURLStorageSaveParam) error
	Search(ctx context.Context, uniqueID uint64) (string, error)
}

type shortURLStorageService struct {
	logger   Logger
	tableCnt int16
}

func NewShortURLStorageService(logger Logger, tableCnt int16) ShortURLStorageService {
	return &shortURLStorageService{
		logger:   logger,
		tableCnt: tableCnt,
	}
}

func (sss *shortURLStorageService) Save(ctx context.Context, param *ShortURLStorageSaveParam) (err error) {
	if param == nil {
		sss.logger.Warnf("[save failed] [none param]")
		return config.PARAM_ERROR
	}
	if err = param.Check(); err != nil {
		sss.logger.Warnf("[save failed] [invalid param]")
		return
	}
	sql := fmt.Sprintf("INSERT INTO %s (long_url, short_url, unique_id) VALUES (?, ?, ?)", sss.genTableName(param.UniqueID))
	sss.logger.Infof("[sql:%s] [param:%+v]", sql, param)
	_, err = repo.SHORT_URL_RECORD_DB.Exec(sql, param.LongURL, param.ShortURL, param.UniqueID)
	if err != nil {
		sss.logger.Errorf("[save failed] [err:%v]", err)
		return config.DB_ERROR
	}
	return
}

func (sss *shortURLStorageService) Search(ctx context.Context, uniqueID uint64) (string, error) {
	if uniqueID == 0 {
		return "", config.PARAM_ERROR
	}
	sqlStr := fmt.Sprintf("SELECT long_url FROM %s WHERE unique_id = ?", sss.genTableName(uniqueID))
	sss.logger.Infof("[sql:%s] [param:%+v]", sqlStr, uniqueID)
	row := repo.SHORT_URL_RECORD_DB.QueryRowContext(ctx, sqlStr, uniqueID)

	var longURL string
	if err := row.Scan(&longURL); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		sss.logger.Errorf("[scan failed] [err:%v]", err)
		return "", config.DB_ERROR
	}
	return longURL, nil
}

func (sss *shortURLStorageService) genTableName(uniqueID uint64) string {
	return fmt.Sprintf("short_url_record_%d", uniqueID%uint64(sss.tableCnt))
}

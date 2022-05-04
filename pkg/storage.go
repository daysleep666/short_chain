package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"hash/crc32"
	"strings"

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

type ShortURLDetail struct {
	ShortURL string
	LongURL  string
	UniqueID uint64
	ViewCnt  int64
}

type ShortURLStorageService interface {
	Save(ctx context.Context, param *ShortURLStorageSaveParam) error
	QueryByUniqueID(ctx context.Context, uniqueID uint64) (detail ShortURLDetail, err error)
	QueryByLongURL(ctx context.Context, longURL string, page uint64, cnt uint64) (output []uint64, err error)
	Count(ctx context.Context, longURL string) (cnt uint64, err error)
	IncViewCnt(uniqueID uint64) (err error)
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
	tx, err := repo.SHORT_URL_RECORD_DB.Begin()
	if err != nil {
		sss.logger.Warnf("[start tx failed] [err:%v]", err)
		err = config.DB_ERROR
		return
	}
	err = func() (err error) {
		{
			sql := fmt.Sprintf("INSERT INTO %s (long_url, short_url, unique_id) VALUES (?, ?, ?)", sss.genShortURLRecordTableName(param.UniqueID))
			sss.logger.Infof("[sql:%s] [param:%+v]", sql, param)
			_, err = tx.Exec(sql, param.LongURL, param.ShortURL, param.UniqueID)
			if err != nil {
				return
			}
		}
		{
			longURLMD5 := sss.longSURL2MD5(param.LongURL)
			sql := fmt.Sprintf("INSERT INTO %s (long_url_md5, long_url, unique_id) VALUES (?, ?, ?)", sss.genLongURLRecordTableName(longURLMD5))
			sss.logger.Infof("[sql:%s] [param:%+v]", sql, param)
			_, err = tx.Exec(sql, longURLMD5, param.LongURL, param.UniqueID)
			if err != nil {
				return
			}
		}
		return nil
	}()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			sss.logger.Errorf("[rollback failed] [err:%v]", err)
			return config.DB_ERROR
		}
		if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
			sss.logger.Info("[duplicate row]")
			return config.EXIST_LONG_URL_ERROR
		}
		sss.logger.Errorf("[save failed] [err:%v]", err)
		return config.DB_ERROR
	}
	if err := tx.Commit(); err != nil {
		sss.logger.Errorf("[commit failed] [err:%v]", err)
		return config.DB_ERROR
	}
	return
}

func (sss *shortURLStorageService) QueryByUniqueID(ctx context.Context, uniqueID uint64) (detail ShortURLDetail, err error) {
	if uniqueID == 0 {
		err = config.PARAM_ERROR
		return
	}
	sqlStr := fmt.Sprintf("SELECT long_url, short_url, view_cnt FROM %s WHERE unique_id = ?", sss.genShortURLRecordTableName(uniqueID))
	sss.logger.Infof("[sql:%s] [param:%+v]", sqlStr, uniqueID)
	row := repo.SHORT_URL_RECORD_DB.QueryRowContext(ctx, sqlStr, uniqueID)

	if err := row.Scan(&detail.LongURL, &detail.ShortURL, &detail.ViewCnt); err != nil {
		if err == sql.ErrNoRows {
			return ShortURLDetail{}, nil
		}
		sss.logger.Errorf("[scan failed] [err:%v]", err)
		return ShortURLDetail{}, config.DB_ERROR
	}
	return
}

func (sss *shortURLStorageService) Count(ctx context.Context, longURL string) (cnt uint64, err error) {
	if len(longURL) == 0 {
		return 0, config.PARAM_ERROR
	}
	longURLMD5 := sss.longSURL2MD5(longURL)
	sqlStr := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE long_url_md5 = ? and long_url = ?", sss.genLongURLRecordTableName(longURLMD5))
	sss.logger.Infof("[sql:%s] [long_url:%s] [long_url_md5:%d]", sqlStr, longURL, longURLMD5)
	row := repo.SHORT_URL_RECORD_DB.QueryRowContext(ctx, sqlStr, longURLMD5, longURL)
	if err := row.Scan(&cnt); err != nil {
		sss.logger.Errorf("[query failed] [err:%v]", err)
		return 0, config.DB_ERROR
	}
	return
}

func (sss *shortURLStorageService) QueryByLongURL(ctx context.Context, longURL string, page uint64, cnt uint64) (output []uint64, err error) {
	if len(longURL) == 0 {
		return nil, config.PARAM_ERROR
	}
	longURLMD5 := sss.longSURL2MD5(longURL)
	sqlStr := fmt.Sprintf("SELECT unique_id FROM %s WHERE long_url_md5 = ? AND long_url = ? LIMIT ?, ?", sss.genLongURLRecordTableName(longURLMD5))
	sss.logger.Infof("[sql:%s] [long_url:%s] [long_url_md5:%d]", sqlStr, longURL, longURLMD5)
	rows, err := repo.SHORT_URL_RECORD_DB.QueryContext(ctx, sqlStr, longURLMD5, longURL, (page-1)*cnt, cnt)
	if err != nil {
		sss.logger.Errorf("[query failed] [err:%v]", err)
		return nil, config.DB_ERROR
	}
	defer rows.Close()
	for rows.Next() {
		var uniqueID uint64
		if err := rows.Scan(&uniqueID); err != nil {
			sss.logger.Errorf("[scan failed] [err:%v]", err)
			return nil, config.DB_ERROR
		}
		output = append(output, uniqueID)
	}

	if err := rows.Err(); err != nil {
		sss.logger.Errorf("[rows err] [err:%v]", err)
		return nil, config.DB_ERROR
	}

	return
}

func (sss *shortURLStorageService) IncViewCnt(uniqueID uint64) (err error) {
	if uniqueID == 0 {
		err = config.PARAM_ERROR
		return
	}

	sql := fmt.Sprintf("UPDATE %s SET view_cnt = view_cnt + 1 WHERE unique_id = ?", sss.genShortURLRecordTableName(uniqueID))
	sss.logger.Infof("[sql:%s] [uniqueID:%+v]", sql, uniqueID)
	_, err = repo.SHORT_URL_RECORD_DB.Exec(sql, uniqueID)
	if err != nil {
		sss.logger.Errorf("[update failed] [err:%v]", err)
		return
	}
	return
}

func (sss *shortURLStorageService) genShortURLRecordTableName(uniqueID uint64) string {
	return fmt.Sprintf("short_url_record_%d", uniqueID%uint64(sss.tableCnt))
}

func (sss *shortURLStorageService) genLongURLRecordTableName(longURLMd5 uint32) string {
	return fmt.Sprintf("long_url_record_%d", longURLMd5%uint32(sss.tableCnt))
}

func (sss *shortURLStorageService) longSURL2MD5(longURL string) uint32 {
	hash := crc32.NewIEEE()
	hash.Write([]byte(longURL))
	return hash.Sum32()
}

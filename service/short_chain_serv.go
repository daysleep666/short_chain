package service

import (
	"context"

	"github.com/daysleep666/short_chain/config"
	"github.com/daysleep666/short_chain/pkg"
)

const (
	MAX_PAGE_CNT = 100
)

type ShortChainService struct {
	log pkg.Logger

	serv struct {
		uniqueID        pkg.UniqueIDService
		converter       pkg.ConverterService
		shortURLStorage pkg.ShortURLStorageService
	}
}

type ShortChainServiceParam struct {
	UniqueIDService        pkg.UniqueIDService
	ConverterService       pkg.ConverterService
	ShortURLStorageService pkg.ShortURLStorageService
	Log                    pkg.Logger
}

func (p *ShortChainServiceParam) Check() error {
	if p.UniqueIDService == nil {
		return config.PARAM_ERROR
	}
	if p.ConverterService == nil {
		return config.PARAM_ERROR
	}
	if p.ShortURLStorageService == nil {
		return config.PARAM_ERROR
	}
	if p.Log == nil {
		return config.PARAM_ERROR
	}
	return nil
}

func NewShortChainService(param *ShortChainServiceParam) (service *ShortChainService, err error) {
	if param == nil {
		err = config.PARAM_ERROR
		return
	}
	if err = param.Check(); err != nil {
		return
	}
	service = &ShortChainService{}
	service.log = param.Log
	service.serv.uniqueID = param.UniqueIDService
	service.serv.converter = param.ConverterService
	service.serv.shortURLStorage = param.ShortURLStorageService
	return
}

func (sc *ShortChainService) Generate(ctx context.Context, longSurl string) (shortURL string, err error) {
	if len(longSurl) == 0 {
		err = config.PARAM_ERROR
		sc.log.Errorf("none long_surl")
		return
	}
	// 申请一个唯一id
	uniqueID, err := sc.serv.uniqueID.Generate(ctx)
	if err != nil {
		sc.log.Errorf("[gen unique id failed] [err:%v]", err)
		return
	}
	sc.log.Infof("[unique_id:%d]", uniqueID)
	// base62
	shortURL = sc.serv.converter.NumberToBase62(uniqueID)

	// 存db
	if err = sc.serv.shortURLStorage.Save(ctx, &pkg.ShortURLStorageSaveParam{
		UniqueID: uniqueID,
		LongURL:  longSurl,
		ShortURL: shortURL,
	}); err != nil {
		sc.log.Errorf("[short_url_storage save failed] [err:%v]", err)
		return
	}

	// 返回生成的短链
	return
}

func (sc *ShortChainService) QueryByShortURL(ctx context.Context, shortURL string) (longURL string, err error) {
	if len(shortURL) == 0 {
		err = config.PARAM_ERROR
		sc.log.Errorf("none shortURL")
		return
	}

	uniqueID := sc.serv.converter.Base62ToNumber(shortURL)

	detail, err := sc.serv.shortURLStorage.QueryByUniqueID(ctx, uniqueID)
	if err != nil {
		return
	}
	longURL = detail.LongURL
	if len(longURL) == 0 {
		err = config.NONE_LONG_URL_ERROR
		return
	}

	sc.serv.shortURLStorage.IncViewCnt(uniqueID)

	return
}

type ShortURLDetail struct {
	ShortURL string `json:"short_url,omitempty"`
	UniqueID uint64 `json:"-"`
	ViewCnt  int64  `json:"view_cnt"`
}

type QueryByLongURLParam struct {
	LongURL string
	Page    uint64
	Cnt     uint64
}

func (p *QueryByLongURLParam) Check() error {
	if len(p.LongURL) == 0 {
		return config.PARAM_ERROR
	}
	if p.Page == 0 {
		return config.PARAM_ERROR
	}
	if p.Cnt == 0 || p.Cnt > MAX_PAGE_CNT {
		return config.PARAM_ERROR
	}
	return nil
}

type QueryByLongURLRes struct {
	Group []*ShortURLDetail `json:"group"`
	Total uint64            `json:"total"`
}

func (sc *ShortChainService) QueryByLongURL(ctx context.Context, param *QueryByLongURLParam) (res QueryByLongURLRes, err error) {
	if param == nil {
		err = config.PARAM_ERROR
		sc.log.Errorf("none longURL")
		return
	}

	uniqueIDGroup, err := sc.serv.shortURLStorage.QueryByLongURL(context.TODO(), param.LongURL, param.Page, param.Cnt)
	if err != nil {
		return
	}
	if len(uniqueIDGroup) == 0 {
		return res, nil
	}

	for _, uniqueID := range uniqueIDGroup {
		detail, err := sc.serv.shortURLStorage.QueryByUniqueID(context.TODO(), uniqueID)
		if err != nil {
			return res, err
		}
		res.Group = append(res.Group, &ShortURLDetail{
			ShortURL: detail.ShortURL,
			UniqueID: detail.UniqueID,
			ViewCnt:  detail.ViewCnt,
		})
	}

	cnt, err := sc.serv.shortURLStorage.Count(context.TODO(), param.LongURL)
	if err != nil {
		return
	}
	res.Total = cnt
	return
}

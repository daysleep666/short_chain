package service

import (
	"context"

	"github.com/daysleep666/short_chain/config"
	"github.com/daysleep666/short_chain/pkg"
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

func (sc *ShortChainService) Search(ctx context.Context, shortURL string) (longURL string, err error) {
	if len(shortURL) == 0 {
		err = config.PARAM_ERROR
		sc.log.Errorf("none shortURL")
		return
	}

	longURL, err = sc.serv.shortURLStorage.Search(ctx, shortURL)
	if err != nil {
		return
	}

	if len(longURL) == 0 {
		err = config.NONE_LONG_URL_ERROR
		return
	}

	return
}

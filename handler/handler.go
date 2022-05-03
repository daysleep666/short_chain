package handler

import (
	"context"

	"github.com/daysleep666/short_chain/pkg"
	"github.com/daysleep666/short_chain/service"
	"github.com/labstack/echo"
)

type GenShortChainReq struct {
	URL string `form:"url"`
}

// GenShortChain 生成短链
func GenShortChain(c echo.Context) (err error) {
	var req GenShortChainReq
	if err = c.Bind(&req); err != nil {
		c.Logger().Errorf("[bind failed] [err:%v]", err)
		return
	}

	ser, err := service.NewShortChainService(&service.ShortChainServiceParam{
		UniqueIDService:        pkg.NewUniqueIDService(c.Logger()),
		ConverterService:       pkg.NewConverterService(),
		ShortURLStorageService: pkg.NewShortURLStorageService(c.Logger()),
		Log:                    c.Logger(),
	})
	if err != nil {
		c.Logger().Errorf("[new short_chain_service failed] [err:%v]", err)
		return MakeResponse(c, err, nil)
	}
	surl, err := ser.Generate(context.TODO(), req.URL)
	if err != nil {
		return MakeResponse(c, err, nil)
	}

	return MakeResponse(c, err, surl)
}

func ShortChainRedirect(c echo.Context) (err error) {
	shortURL := c.Param("url")

	ser, err := service.NewShortChainService(&service.ShortChainServiceParam{
		UniqueIDService:        pkg.NewUniqueIDService(c.Logger()),
		ConverterService:       pkg.NewConverterService(),
		ShortURLStorageService: pkg.NewShortURLStorageService(c.Logger()),
		Log:                    c.Logger(),
	})
	if err != nil {
		c.Logger().Errorf("[new short_chain_service failed] [err:%v]", err)
		return MakeResponse(c, err, nil)
	}
	longURL, err := ser.Search(context.TODO(), shortURL)
	if err != nil {
		return MakeResponse(c, err, nil)
	}
	c.Logger().Infof("redirect to %s", longURL)
	return c.Redirect(302, longURL)
}

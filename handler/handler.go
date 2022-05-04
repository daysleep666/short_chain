package handler

import (
	"context"
	"fmt"

	"github.com/daysleep666/short_chain/config"
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

	ser, err := newShortChainService(c.Logger())
	if err != nil {
		c.Logger().Errorf("[new short_chain_service failed] [err:%v]", err)
		return MakeResponse(c, err, nil)
	}
	surl, err := ser.Generate(context.TODO(), req.URL)
	if err != nil {
		return MakeResponse(c, err, nil)
	}

	return MakeResponse(c, err, map[string]string{
		"surl": fmt.Sprintf("%s/%s", config.CONFIG_INSTANCE.ServerConfig.Domain, surl),
	})
}

type QueryShortChainReq struct {
	URL string `form:"url"`
}

func QueryShortChain(c echo.Context) (err error) {
	var req QueryShortChainReq
	if err = c.Bind(&req); err != nil {
		c.Logger().Errorf("[bind failed] [err:%v]", err)
		return
	}

	ser, err := newShortChainService(c.Logger())
	if err != nil {
		c.Logger().Errorf("[new short_chain_service failed] [err:%v]", err)
		return MakeResponse(c, err, nil)
	}

	res, err := ser.QueryByLongURL(context.TODO(), req.URL)
	if err != nil {
		c.Logger().Errorf("[query by long_url failed] [err:%v]", err)
		return MakeResponse(c, err, nil)
	}
	return MakeResponse(c, nil, res)
}

// ShortChainRedirect 短链重定向
func ShortChainRedirect(c echo.Context) (err error) {
	shortURL := c.Param("url")

	ser, err := newShortChainService(c.Logger())
	if err != nil {
		c.Logger().Errorf("[new short_chain_service failed] [err:%v]", err)
		return MakeResponse(c, err, nil)
	}
	longURL, err := ser.QueryByShortURL(context.TODO(), shortURL)
	if err != nil {
		return MakeResponse(c, err, nil)
	}
	c.Logger().Infof("redirect to %s", longURL)
	return c.Redirect(302, longURL)
}

func newShortChainService(logger pkg.Logger) (*service.ShortChainService, error) {
	return service.NewShortChainService(&service.ShortChainServiceParam{
		UniqueIDService:        pkg.NewUniqueIDSnowflakeService(config.CONFIG_INSTANCE.ServerConfig.MachineID),
		ConverterService:       pkg.NewConverterService(),
		ShortURLStorageService: pkg.NewShortURLStorageService(logger, config.CONFIG_INSTANCE.ShortURLMysqlConfig.TableCnt),
		Log:                    logger,
	})
}

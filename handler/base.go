package handler

import (
	"github.com/daysleep666/short_chain/config"
	"github.com/labstack/echo"
)

type BaseResposne struct {
	Msg        string      `json:"msg,omitempty"`
	StatusCode int         `json:"status_code,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func MakeResponse(c echo.Context, err error, data interface{}) error {
	if err == nil {
		err = config.SUCCESS_ERROR
	}
	confErr, ok := err.(config.Error)
	if !ok {
		c.Logger().Errorf("[unknown err] [err:%v]", err)
		confErr = config.UNKNOWN_ERROR
	}
	return c.JSON(confErr.HTTPCode, BaseResposne{
		Msg:        confErr.Msg,
		StatusCode: confErr.StatusCode,
		Data:       data,
	})
}

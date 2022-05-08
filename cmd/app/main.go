package main

import (
	"github.com/daysleep666/short_chain/handler"
	"github.com/daysleep666/short_chain/mid"
	"github.com/daysleep666/short_chain/pkg"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

const (
	CONFIG_PATH = "./config/app.toml"
)

func main() {
	pkg.MustInit(CONFIG_PATH)

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
		Output: log.Output(),
	}))
	e.Use(mid.AddLogger)
	e.Use(mid.ReqStart)
	e.Use(middleware.CORS())
	e.GET("/shortchain/gen", handler.GenShortChain)
	e.POST("/shortchain/gen", handler.GenShortChain)
	e.GET("/shortchain/query", handler.QueryShortChain)
	e.POST("/shortchain/query", handler.QueryShortChain)
	e.File("/favicon.ico", "static/images/favicon.ico")
	e.GET("/:url", handler.ShortChainRedirect)
	e.Static("web", "static")

	e.Logger.Fatal(e.Start(":1234"))
}

package mid

import (
	"context"
	"fmt"

	"github.com/daysleep666/short_chain/config"
	"github.com/daysleep666/short_chain/pkg"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func AddLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Logger().SetOutput(log.Output())
		c.Logger().SetLevel(log.Lvl(config.CONFIG_INSTANCE.LoggerConfig.Level))
		logid, _ := pkg.NewUniqueIDSnowflakeService(0).Generate(context.TODO())
		c.Logger().SetPrefix(fmt.Sprintf("[logid:%d]", logid))
		return next(c)
	}
}

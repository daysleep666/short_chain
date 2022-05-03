package mid

import (
	"fmt"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func AddLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Logger().SetOutput(log.Output())
		c.Logger().SetLevel(log.DEBUG)
		c.Logger().SetPrefix(fmt.Sprintf("[logid:%d]", time.Now().Unix()))
		return next(c)
	}
}

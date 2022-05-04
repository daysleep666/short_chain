package pkg

import (
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

type Logger echo.Logger

func InitLogger(filePath string) (err error) {
	f, err := os.Create(filePath)
	if err != nil {
		return
	}
	log.SetOutput(f)
	return
}

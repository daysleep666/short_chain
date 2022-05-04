package config

import "fmt"

type Error struct {
	Msg        string
	StatusCode int
	HTTPCode   int
}

func (e Error) Error() string {
	return fmt.Sprintf("[msg:%s] [status_code:%d] [http_code:%d]", e.Msg, e.StatusCode, e.HTTPCode)
}

var (
	SUCCESS_ERROR = Error{Msg: "success", StatusCode: 0, HTTPCode: 200}

	PARAM_ERROR = Error{Msg: "param err", StatusCode: 1, HTTPCode: 403}

	NONE_LONG_URL_ERROR  = Error{Msg: "none long_url", StatusCode: 100, HTTPCode: 403}
	EXIST_LONG_URL_ERROR = Error{Msg: "exist long_url", StatusCode: 101, HTTPCode: 403}

	DB_ERROR = Error{Msg: "db err", StatusCode: 400, HTTPCode: 503}

	UNKNOWN_ERROR = Error{Msg: "unknown", StatusCode: 500, HTTPCode: 503}
)

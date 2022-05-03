package repo

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	SHORT_URL_RECORD_DB *sql.DB
)

func InitShortUrlRecordDB(addr, user, pwd, dbName string) (err error) {
	SHORT_URL_RECORD_DB, err = initMysql(addr, user, pwd, dbName)
	return
}

func initMysql(addr, user, pwd, dbName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", user, pwd, addr, dbName)) // 设置连接数据库的参数
	if err != nil {
		return nil, err
	}
	err = db.Ping() //连接数据库
	if err != nil {
		return nil, err
	}
	return db, nil
}

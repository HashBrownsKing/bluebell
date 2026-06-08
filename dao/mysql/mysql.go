package mysql

import (
	"bluebell/settings"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // 必须导入MySQL驱动
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Init(config *settings.AppConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.MySQLConfig.User,
		config.MySQLConfig.Password,
		config.MySQLConfig.Host,
		config.MySQLConfig.Port,
		config.MySQLConfig.DBName,
	)

	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(config.MySQLConfig.MaxOpenConnects)
	db.SetMaxIdleConns(config.MySQLConfig.MaxIdleConnects)
	db.SetConnMaxLifetime(time.Duration(config.MySQLConfig.MaxLifeTime) * time.Hour)
	db.SetConnMaxIdleTime(time.Duration(config.MySQLConfig.MaxIdleTime) * time.Hour)
	return
}

func Close() {
	_ = db.Close()
}

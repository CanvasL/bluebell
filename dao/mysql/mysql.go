package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go_web/gin/bluebell/setting"
)

var db *sqlx.DB

//func Init() (err error) {
//	dsn := fmt.Sprintf(
//		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
//		viper.GetString("mysql.user"),
//		viper.GetString("mysql.password"),
//		viper.GetString("mysql.host"),
//		viper.GetInt("mysql.port"),
//		viper.GetString("mysql.dbname"),
//	)
//	// 也可以使用MustConnect连接不成功就panic
//	db, err = sqlx.Connect("mysql", dsn)
//	if err != nil {
//		fmt.Printf("connect DB failed, err:%v\n", err)
//		return
//	}
//	db.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))
//	db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))
//	return
//}

// Init 初始化MySQL连接
func Init(cfg *setting.MySQLConfig) (err error) {
	// "user:password@tcp(host:port)/dbname"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return
}

func Close() {
	_ = db.Close()
}

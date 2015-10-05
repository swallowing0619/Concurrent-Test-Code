package main

import (
//"github.com/op/go-logging"
// gr "github.com/parnurzeal/gorequest"
// ws "golang.org/x/net/websocket"
//"os"
// "sync"
//"time"
//mo "yto.net.cn/kefu/models"
// o "yto.net.cn/kefu/models/objects"
)

const (
	PortalAddr = "222.73.41.159:9008"
	//mysql
	// MysqlHost        = "222.73.41.159"
	// MysqlPort        = 3306
	// MysqlDatabase    = "kefu_dev"
	// MysqlUser        = "kefu_dev_user"
	// MysqlPass        = "kefu_dev_password"
	// MysqlMaxConn     = 1024
	// MysqlMaxIdleConn = 64

	// //redis
	// RedisServer      = "222.73.41.159:6379"
	// RedisPass        = ""
	// RediMaxIdle      = 5
	// RedisMaxActive   = 0
	// RedisIdleTimeout = "20s"

	//util
	EntCode = "yto_dev"
)

// var (
// 	Mo  *mo.Models
// 	log *logging.Logger
// )

// func InitLogger() {
// 	log = logging.MustGetLogger("ws")
// 	format := logging.MustStringFormatter(
// 		"%{color}%{time:15:04:05.000000} ▶ %{level:.4s} %{module} %{shortfunc} %{id:03x}%{color:reset} %{message}",
// 	)
// 	backend := logging.NewLogBackend(os.Stderr, "", 0)
// 	backendFormatter := logging.NewBackendFormatter(backend, format)
// 	backendLeveled := logging.AddModuleLevel(backendFormatter)
// 	backendLeveled.SetLevel(logging.DEBUG, "ws")
// 	logging.SetBackend(backendLeveled)
// }

// func init() {
// 	Mo = new(mo.Models)
// 	InitLogger()

// 	if err := InitMysql(); err != nil {
// 		log.Error("Error Init Mysql:%s", err.Error())
// 		panic(err)
// 	}

// 	if err := InitRedis(); err != nil {
// 		log.Error("Error Init Redis:%s", err.Error())
// 		panic(err)
// 	}
// }

// func InitRedis() error {
// 	timeout, err := time.ParseDuration(RedisIdleTimeout)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return Mo.InitializeRedis(RedisServer, RedisPass, RediMaxIdle, RedisMaxActive, timeout)
// }

// func InitMysql() error {
// 	return Mo.InitializeMySQL(MysqlHost, MysqlPort, MysqlDatabase, MysqlUser, MysqlPass, MysqlMaxConn, MysqlMaxIdleConn)
// }

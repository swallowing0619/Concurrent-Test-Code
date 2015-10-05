package main

import (
	"github.com/op/go-logging"
	"os"
	"time"
	"yto.net.cn/kefu/models"
)

var (
	Mo  *models.Models
	log *logging.Logger
)

func init() {
	Mo = new(models.Models)
	InitLogger("SubLoadTest")
	InitModels()
}

func InitLogger(typ string) {
	log = logging.MustGetLogger(typ)
	format := logging.MustStringFormatter(
		"%{color}%{time:15:04:05.000000} ▶ %{level:.4s} %{module} %{shortfunc} %{id:03x}%{color:reset} %{message}",
	)
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(logging.DEBUG, "redis")
	logging.SetBackend(backendLeveled)
}

func InitModels() {

	err := Mo.InitializeMySQL("kefu-dev.hotpu.cn", 3306, "kefu_dev", "kefu_dev_user", "kefu_dev_password", 20, 0)

	if err != nil {
		panic(err)
	}

	idletimeout, err := time.ParseDuration("240s")
	if err != nil {
		panic(err)
	}
	err = Mo.InitializeRedis("127.0.0.1:6379", "", 5, 0, idletimeout)

	if err != nil {
		panic(err)
	}

}

package main

import (
	"net/http"
	"stabunkow/http/pkg/logging"
	"stabunkow/http/pkg/mongodb"
	"stabunkow/http/pkg/redis"
	"stabunkow/http/pkg/setting"
	"stabunkow/http/pkg/wechat"

	"github.com/gin-gonic/gin"
)

func init() {
	setting.Setup("../configs/app.ini")
	logging.Setup()
	mongodb.Setup()
	redis.Setup()
	wechat.Setup()
}

func main() {
	gin.SetMode(setting.ServerSetting.RunMode)

	server, err := InitializeServer()
	if err != nil {
		panic(err)
	}

	server.ListenAndServe()
}

func NewServer(router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:           setting.ServerSetting.Addr,
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeOut,
		WriteTimeout:   setting.ServerSetting.WriteTimeOut,
		MaxHeaderBytes: 1 << 20,
	}
}

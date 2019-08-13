package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

type App struct {
	LogSaveName string
	LogPath     string
}

var AppSetting = &App{}

type Server struct {
	Addr         string
	RunMode      string
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
}

var ServerSetting = &Server{}

type Wechat struct {
	AppId  string
	Secret string
}

var WechatSetting = &Wechat{}

type Mongodb struct {
	Addr     string
	Database string
}

var MongodbSetting = &Mongodb{}

type Redis struct {
	Addr string
}

var RedisSetting = &Redis{}

var cfg *ini.File

// Setup initialize the configuration instance
func Setup(path string) {
	var err error
	cfg, err = ini.Load(path)

	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'configs/app.ini': %v", err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("wechat", WechatSetting)
	mapTo("mongodb", MongodbSetting)
	mapTo("redis", RedisSetting)
	ServerSetting.ReadTimeOut = ServerSetting.ReadTimeOut * time.Second
	ServerSetting.WriteTimeOut = ServerSetting.WriteTimeOut * time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)

	if err != nil {
		log.Fatalf("Cfg.MapTo %v err: %v", section, err)
	}
}

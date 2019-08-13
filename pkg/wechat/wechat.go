package wechat

import (
	"stabunkow/http/pkg/setting"

	"github.com/medivhzhan/weapp"
)

var DefaultWechat *Wechat

func GetDefaultWechat() *Wechat {
	return DefaultWechat
}

func Setup() error {
	DefaultWechat = NewWechat(
		setting.WechatSetting.AppId,
		setting.WechatSetting.Secret,
	)

	return nil
}

type Wechat struct {
	AppId  string
	Secret string
}

func NewWechat(appId, secret string) *Wechat {
	return &Wechat{
		AppId:  appId,
		Secret: secret,
	}
}

func (w *Wechat) Login(code string) (weapp.LoginResponse, error) {
	return weapp.Login(w.AppId, w.Secret, code)
}

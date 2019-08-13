package api

import (
	"net/http"
	"stabunkow/http/pkg/logging"
	"stabunkow/http/pkg/wechat"
	"stabunkow/http/repository/user_repository"
	"stabunkow/http/service/auth_service"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	wechat         *wechat.Wechat
	userRepository *user_repository.UserRepository
	authService    *auth_service.AuthService
}

func NewAuthController(
	wechat *wechat.Wechat,
	userRepository *user_repository.UserRepository,
	authService *auth_service.AuthService,
) *AuthController {
	return &AuthController{
		wechat:         wechat,
		userRepository: userRepository,
		authService:    authService,
	}
}

type wxLoginRequest struct {
	Code string `valid:"Required; Length(32)"`
}

func (ctr *AuthController) WxLogin(c *gin.Context) {
	code := c.PostForm("code")
	req := &wxLoginRequest{code}

	valid := validation.Validation{}

	ok, _ := valid.Valid(req)
	if !ok {
		logging.Error(valid.Errors)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "code invaild",
		})
		return
	}

	rsp, err := ctr.wechat.Login(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	sid, err := ctr.authService.WxLogin(rsp.OpenID, rsp.SessionKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sid": sid,
	})
	return
}

package routers

import (
	"stabunkow/http/routers/api"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	authController *api.AuthController,
) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api := r.Group("/api")

	api.POST("/wxLogin", authController.WxLogin)

	return r
}

// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/google/wire"
	"net/http"
	"stabunkow/http/pkg/mongodb"
	"stabunkow/http/pkg/redis"
	"stabunkow/http/pkg/wechat"
	"stabunkow/http/repository/user_repository"
	"stabunkow/http/routers"
	"stabunkow/http/routers/api"
	"stabunkow/http/service/auth_service"
)

// Injectors from wire.go:

func InitializeUserRepository() (*user_repository.UserRepository, error) {
	db := mongodb.GetDefaultDb()
	cache := redis.GetDefaultCache()
	userRepository := user_repository.NewUserRepository(db, cache)
	return userRepository, nil
}

func InitializeAuthService() (*auth_service.AuthService, error) {
	db := mongodb.GetDefaultDb()
	cache := redis.GetDefaultCache()
	userRepository := user_repository.NewUserRepository(db, cache)
	authService := auth_service.NewAuthService(userRepository)
	return authService, nil
}

func InitializeAuthController() (*api.AuthController, error) {
	wechatWechat := wechat.GetDefaultWechat()
	userRepository, err := InitializeUserRepository()
	if err != nil {
		return nil, err
	}
	authService, err := InitializeAuthService()
	if err != nil {
		return nil, err
	}
	authController := api.NewAuthController(wechatWechat, userRepository, authService)
	return authController, nil
}

func InitializeServer() (*http.Server, error) {
	authController, err := InitializeAuthController()
	if err != nil {
		return nil, err
	}
	engine := routers.NewRouter(authController)
	server := NewServer(engine)
	return server, nil
}

// wire.go:

var UserRepositorySet = wire.NewSet(mongodb.GetDefaultDb, redis.GetDefaultCache)

var AuthControllerSet = wire.NewSet(wechat.GetDefaultWechat, InitializeUserRepository,
	InitializeAuthService,
)
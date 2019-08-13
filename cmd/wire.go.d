package main

import (
	"net/http"
	"stabunkow/http/pkg/mongodb"
	"stabunkow/http/pkg/redis"
	"stabunkow/http/pkg/wechat"
	"stabunkow/http/repository/user_repository"
	"stabunkow/http/routers"
	"stabunkow/http/routers/api"
	"stabunkow/http/service/auth_service"

	"github.com/google/wire"
)

var UserRepositorySet = wire.NewSet(
	mongodb.GetDefaultDb,
	redis.GetDefaultCache,
)

func InitializeUserRepository() (*user_repository.UserRepository, error) {
	wire.Build(
		user_repository.NewUserRepository,
		UserRepositorySet,
	)
	return nil, nil
}

func InitializeAuthService() (*auth_service.AuthService, error) {
	wire.Build(
		auth_service.NewAuthService,
		user_repository.NewUserRepository,
		UserRepositorySet,
	)

	return nil, nil
}

var AuthControllerSet = wire.NewSet(
	wechat.GetDefaultWechat,
	InitializeUserRepository,
	InitializeAuthService,
)

func InitializeAuthController() (*api.AuthController, error) {
	wire.Build(
		api.NewAuthController,
		AuthControllerSet,
	)

	return nil, nil
}

func InitializeServer() (*http.Server, error) {
	wire.Build(
		NewServer,
		routers.NewRouter,
		InitializeAuthController,
	)

	return nil, nil
}

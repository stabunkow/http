package auth_service

import (
	"stabunkow/http/pkg/util"
	"stabunkow/http/repository/user_repository"

	"gopkg.in/mgo.v2"
)

type AuthService struct {
	userRepository *user_repository.UserRepository
}

func NewAuthService(userRepo *user_repository.UserRepository) *AuthService {
	return &AuthService{userRepo}
}

func (svc *AuthService) WxLogin(openId, sessionKey string) (string, error) {
	user, err := svc.userRepository.FindUserByWechatOpenId(openId)
	if err != nil {
		if err != mgo.ErrNotFound {
			return "", err
		}
		user, err = svc.userRepository.CreateUser(openId)
		if err != nil {
			return "", err
		}
	}

	sid := util.UniqueId()
	svc.userRepository.UpdateUserSessionKeyById(user.GetId(), sessionKey)
	svc.userRepository.UpdateUserSidById(user.GetId(), sid)

	return sid, nil
}

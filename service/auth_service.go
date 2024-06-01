package service

import (
	"etalert-backend/repository"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) AuthService {
	return &authService{userRepository: userRepository}
}

func (s authService) Login(loginInput *LoginInput) (LoginResponse, error) {
	userRepo, err := s.userRepository.GetUserInfo(loginInput.GoogleId)
	if err != nil {
		return LoginResponse{}, err
	}

	claims := jwt.MapClaims{
		"googleId": userRepo.GoogleId,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return LoginResponse{}, err
	}
	return LoginResponse{Token: t}, nil
}

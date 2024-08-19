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

	acExp := time.Hour * 1
	accessToken, err := s.generateToken(userRepo.GoogleId, acExp) // 1 hour expiration for access token
	acExpTime := time.Now().Add(acExp).Unix()
	acExpTimeUnix := time.Unix(acExpTime, 0)
	acReadableTime := acExpTimeUnix.Format(time.RFC3339)
	if err != nil {
		return LoginResponse{}, err
	}

	rfExp := time.Hour * 24 * 7
	refreshToken, err := s.generateToken(userRepo.GoogleId, rfExp) // 7 days expiration for refresh token
	rfExpTime := time.Now().Add(rfExp).Unix()
	rfExpTimeUnix := time.Unix(rfExpTime, 0)
	rfReadableTime := rfExpTimeUnix.Format(time.RFC3339)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken, AccessTokenExpired: acReadableTime, RefreshTokenExpired: rfReadableTime}, nil
}

func (s *authService) generateToken(googleId string, expiresIn time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"googleId": googleId,
		"exp":      time.Now().Add(expiresIn).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (s *authService) RefreshToken(refreshToken string) (LoginResponse, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return LoginResponse{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return LoginResponse{}, err
	}

	googleId, ok := claims["googleId"].(string)
	if !ok {
		return LoginResponse{}, err
	}

	// Generate new access token
	newAccessToken, err := s.generateToken(googleId, time.Hour*1)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{AccessToken: newAccessToken}, nil
}
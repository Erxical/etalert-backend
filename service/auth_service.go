package service

import (
	"etalert-backend/repository"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) AuthService {
	return &authService{userRepository: userRepository}
}

func (s *authService) Login(loginInput *LoginInput) (LoginResponse, error) {
	userRepo, err := s.userRepository.GetUserInfo(loginInput.GoogleId)
	if err != nil {
		return LoginResponse{}, err
	}

	accessToken, refreshToken, accessTokenExpire, refreshTokenExpire, err := s.generateToken(userRepo.GoogleId)
	if err != nil {
		return LoginResponse{}, err
	}
	return LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken, AccessTokenExpired: accessTokenExpire, RefreshTokenExpired: refreshTokenExpire}, nil
}

func (s *authService) generateToken(googleId string) (string, string, string, string, error) {
	accessExpire := time.Hour * 24
	acExpTime := time.Now().Add(accessExpire).Unix()
	acExpTimeUnix := time.Unix(acExpTime, 0)
	acReadableTime := acExpTimeUnix.Format(time.RFC3339)
	aClaims := jwt.MapClaims{
		"googleId": googleId,
		"exp":      time.Now().Add(accessExpire).Unix(),
	}
	refreshExpire := time.Hour * 24 * 7
	rfExpTime := time.Now().Add(refreshExpire).Unix()
	rfExpTimeUnix := time.Unix(rfExpTime, 0)
	rfReadableTime := rfExpTimeUnix.Format(time.RFC3339)
	rClaims := jwt.MapClaims{
		"googleId": googleId,
		"exp":      time.Now().Add(refreshExpire).Unix(),
	}

	aToken := jwt.NewWithClaims(jwt.SigningMethodHS256, aClaims)
	aTokenString, err := aToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", "", "", "", err
	}

	rToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rClaims)
	rTokenString, err := rToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", "", "", "", err
	}

	return aTokenString, rTokenString, acReadableTime, rfReadableTime, nil
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
	accessToken, _, accessTokenExpire, _, err := s.generateToken(googleId)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{AccessToken: accessToken, AccessTokenExpired: accessTokenExpire}, nil
}

func (s *authService) ValidateAccessToken(tokenString string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return nil, err
    }

    return claims, nil
}

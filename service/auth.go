package service

type LoginInput struct {
	GoogleId string `bson:"googleId"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
	AccessTokenExpired      string `json:"expired"`
	RefreshTokenExpired      string `json:"refreshExpired,omitempty"`
}


type AuthService interface {
	Login(loginInput *LoginInput) (LoginResponse, error)
	RefreshToken(refreshToken string) (LoginResponse, error)
}

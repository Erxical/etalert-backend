package service

type LoginInput struct {
	GoogleId string `bson:"googleId"`
}

type LoginResponse struct {
	AccessToken  string `bson:"accessToken"`
	RefreshToken string `bson:"refreshToken"`
	AccessTokenExpired      string `bson:"expired"`
	RefreshTokenExpired      string `bson:"expired"`
}


type AuthService interface {
	Login(loginInput *LoginInput) (LoginResponse, error)
	RefreshToken(refreshToken string) (LoginResponse, error)
}

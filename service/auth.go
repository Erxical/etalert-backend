package service

type LoginInput struct {
	GoogleId string `bson:"googleId"`
	Expired  string `bson:"expired"`
}

type LoginResponse struct {
	Token string `bson:"token"`
}

type AuthService interface {
	Login(loginInput *LoginInput) (LoginResponse, error)
}

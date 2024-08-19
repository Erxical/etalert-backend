package service

type LoginInput struct {
	GoogleId string `bson:"googleId"`
	// Expired  string `bson:"expired"`
}

type LoginResponse struct {
	Token string `bson:"token"`
	Expired string `bson:"expired"`
}

type AuthService interface {
	Login(loginInput *LoginInput) (LoginResponse, error)
}

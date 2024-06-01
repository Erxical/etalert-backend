package service

type UserInput struct {
	GoogleId string `bson:"googleId"`
	Email    string `bson:"email"`
	Name     string `bson:"name"`
	Image    string `bson:"image"`
}

type UserUpdater struct {
	Name  string `bson:"name"`
	Image string `bson:"image"`
}

type UserInfoResponse struct {
	Name  string `bson:"name"`
	Image string `bson:"image"`
	Email string `bson:"email"`
}

type UserService interface {
	InsertUser(user *UserInput) error
	GetUserInfo(string) (*UserInfoResponse, error)
	UpdateUser(string, *UserUpdater) error
}

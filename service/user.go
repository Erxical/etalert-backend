package service

type UserInput struct {
	Name     string `json:"name"`
	Image    string `json:"image"`
	Email    string `json:"email"`
	GoogleId string `json:"googleId"`
}

type UserResponse struct {
	Name  string `bson:"name"`
	Image string `bson:"image"`
	Email string `bson:"email"`
}

type UserService interface {
	InsertUser(user *UserInput) error
	GetUser(string) (*UserResponse, error)
}

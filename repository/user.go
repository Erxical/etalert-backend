package repository

type User struct {
	GoogleId string `bson:"googleId"`
	Email    string `bson:"email"`
	Name     string `bson:"name"`
	Image    string `bson:"image"`
}

type UserRepository interface {
	InsertUser(user *User) error
	GetUserInfo(string) (*User, error)
}

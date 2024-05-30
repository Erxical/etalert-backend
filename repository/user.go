package repository

type User struct {
	Name     string `bson:"name"`
	Image    string `bson:"image"`
	Email    string `bson:"email"`
	GoogleId string `bson:"googleId"`
}

type UserRepository interface {
	InsertUser(user *User) error
	GetUser(string) (*User, error)
}

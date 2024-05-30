package service

import (
	"errors"
	"etalert-backend/repository"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

var ErrUserAlreadyExists = errors.New("user already exists")

func (s *userService) InsertUser(user *UserInput) error {
	existingUser, err := s.userRepo.GetUser(user.GoogleId)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return ErrUserAlreadyExists
	}
	err = s.userRepo.InsertUser(&repository.User{
		Name:     user.Name,
		Image:    user.Image,
		Email:    user.Email,
		GoogleId: user.GoogleId,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s userService) GetUser(gId string) (*UserResponse, error) {
	user, err := s.userRepo.GetUser(gId)
	if err != nil {
		return nil, err
	}

	userResponse := UserResponse{
		Name:  user.Name,
		Image: user.Image,
		Email: user.Email,
	}

	return &userResponse, nil
}

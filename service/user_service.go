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

func (s userService) InsertUser(user *UserInput) error {
	existingUser, err := s.userRepo.GetUserInfo(user.GoogleId)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return ErrUserAlreadyExists
	}

	err = s.userRepo.InsertUser(&repository.User{
		GoogleId: user.GoogleId,
		Email:    user.Email,
		Name:     user.Name,
		Image:    user.Image,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s userService) GetUserInfo(gId string) (*UserInfoResponse, error) {
	user, err := s.userRepo.GetUserInfo(gId)
	if err != nil {
		return nil, err
	}

	userResponse := UserInfoResponse{
		Name:  user.Name,
		Image: user.Image,
		Email: user.Email,
	}

	return &userResponse, nil
}

func (s userService) UpdateUser(gId string, user *UserUpdater) error {
	err := s.userRepo.UpdateUser(gId, &repository.User{
		Name:  user.Name,
		Image: user.Image,
	})
	if err != nil {
		return err
	}
	return nil
}

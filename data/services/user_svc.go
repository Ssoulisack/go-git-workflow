package services

import "go-fiber/data/repositories"

type UserService interface {
	//Methods
	GetAllUsers()
}

type userService struct {
	userRepo repositories.UserRepository
}

// GetAllUsers implements UserService.
func (u *userService) GetAllUsers() {
	panic("unimplemented")
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

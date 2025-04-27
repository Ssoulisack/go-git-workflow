package services

import (
	"go-fiber/core/utilities"
	"go-fiber/data/repositories"
	"go-fiber/domain/entities"
	"go-fiber/domain/models"
)

type UserService interface {
	//Methods
	CreateUser(data models.UserReq) error
	GetUserByID(id int) (*models.UserRes, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

// CreateUser implements UserService.
func (u *userService) CreateUser(data models.UserReq) error {
	// dont use tx here
	// because we are using transaction in repository
	//convert data to entity
	dataEntity := utilities.ConvertModelToEntity[models.UserReq, entities.UserEntity](data)
	if err := u.userRepo.CreateUser(dataEntity); err != nil {
		return err
	}
	return nil
}

// GetUserByID implements UserService.
func (u *userService) GetUserByID(id int) (*models.UserRes, error) {
	user, err := u.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	//convert entity to model
	userModel := utilities.ConvertEntityToModel[entities.UserEntity, models.UserRes](user)
	return &userModel, nil
}



func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

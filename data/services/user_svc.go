package services

import (
	"go-fiber/api/rest/middleware"
	"go-fiber/core/utilities"
	"go-fiber/data/repositories"
	"go-fiber/domain/entities"
	"go-fiber/domain/models"
)

type UserService interface {
	//Methods
	CreateUser(data models.UserReq) error
	GetUserByID(id int) (*models.UserRes, error)
	GetAllUsers(req middleware.PageQuery) (*middleware.PageQuery, error)
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

// GetAllUsers implements UserService. use pagination
func (u *userService) GetAllUsers(req middleware.PageQuery) (*middleware.PageQuery, error) {
	//convert entity to model
	page, users, err := u.userRepo.GetAllUsers(req)
	if err != nil {
		return nil,  err
	}
	//convert entity to model
	userModels := utilities.ConvertEntitiesToModels[entities.UserEntity, models.UserRes](users)

	if userModels == nil {
		page.Rows = []models.UserRes{}
	}

	page.Rows = userModels

	return page, nil
}



func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

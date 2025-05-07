package repositories

import (
	"fmt"
	"go-fiber/api/rest/middleware"
	"go-fiber/core/utilities"
	"go-fiber/domain/entities"

	"gorm.io/gorm"
)

type UserRepository interface {
	//Methods
	CreateUser(data entities.UserEntity) error
	GetUserByID(id int) (entities.UserEntity, error)
	GetAllUsers(req middleware.PageQuery) (*middleware.PageQuery, []entities.UserEntity, error)
	UpdateUser(id int, data entities.UserEntity) error
	DeleteUser(id int) error
}

type userRepository struct {
	db *gorm.DB
}

// GetAllUsers implements UserRepository.
func (u *userRepository) GetAllUsers(req middleware.PageQuery) (*middleware.PageQuery, []entities.UserEntity, error) {
	var users []entities.UserEntity
	var total int64
	var search string

	if req.Search != "" {
		search = fmt.Sprintf("name LIKE '%%%s%%' AND deleted_at IS NULL", req.Search)
	}

	tx := u.db.Model(&entities.UserEntity{})

	// Count total records
	if err := tx.Count(&total).Where(search).Error; err != nil {
		return nil, nil, err
	}

	offset := utilities.CalculateOffset(req.Page, req.Limit)
	totalPage := utilities.CalculatePageSize(total, req.Limit)

	if req.Limit < -1 {
		req.Limit = int(total)
	}

	if err := tx.Offset(offset).Limit(req.Limit).Where(search).Find(&users).Error; err != nil {
		return nil, nil, err
	}

	if req.Limit == 0 {
		totalPage = 0
	}

	query := &middleware.PageQuery{
		TotalPages: totalPage,
		TotalRows:  total,
		Page:       req.Page,
		Limit:      req.Limit,
		Rows:       users,
	}
	return query, users, nil
}
func (u *userRepository) CreateUser(data entities.UserEntity) error {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (u *userRepository) GetUserByID(id int) (entities.UserEntity, error) {
	var user entities.UserEntity
	if err := u.db.First(&user, id).Error; err != nil {
		return entities.UserEntity{}, err
	}
	return user, nil
}

func (u *userRepository) UpdateUser(id int, data entities.UserEntity) error {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Model(&entities.UserEntity{}).Where("id = ?", id).Updates(data).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (u *userRepository) DeleteUser(id int) error {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Delete(&entities.UserEntity{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}


func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

package repositories

import (
	"go-fiber/domain/entities"

	"gorm.io/gorm"
)

type UserRepository interface {
	//Methods
	GetAllUsers()
}

type userRepository struct {
	db *gorm.DB
}

// GetAllUsers implements UserRepository.
func (u *userRepository) GetAllUsers() {
	panic("unimplemented")
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


func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

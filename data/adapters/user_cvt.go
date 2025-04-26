package adapter

import (
	"go-fiber/domain/entities"
	"go-fiber/domain/models"
)

func UserEntityToModel(data entities.UserEntity) models.UserRes {
	return models.UserRes{
		ID:       data.ID,
		Name:     data.Name,
		Email:    data.Email,
	}
}

func UserModelsToEntities(data []models.UserReq) []entities.UserEntity {
	var result []entities.UserEntity
	for _, v := range data {
		result = append(result, entities.UserEntity{
			ID:       v.ID,
			Name:     v.Name,
			Email:    v.Email,
			Password: v.Password,
		})
	}
	return result
}

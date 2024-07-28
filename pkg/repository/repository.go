package repository

import "github.com/megajandrox/go-finance-api/pkg/models"

type Repository interface {
	Create(position *models.Position) error
	GetAll() ([]models.Position, error)
	GetByID(id uint) (*models.Position, error)
	Update(position *models.Position) error
	Delete(id uint) error
}

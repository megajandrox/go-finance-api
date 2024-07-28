package repository

import (
	"github.com/megajandrox/go-finance-api/pkg/models"
	"gorm.io/gorm"
)

type PositionRepository interface {
	Create(position *models.Position) error
	GetAll() ([]models.Position, error)
	GetByID(id uint) (*models.Position, error)
	Update(position *models.Position) error
	Delete(id uint) error
}

type positionRepository struct {
	db *gorm.DB
}

func NewPositionRepository(db *gorm.DB) PositionRepository {
	return &positionRepository{db}
}

func (r *positionRepository) Create(position *models.Position) error {
	return r.db.Create(position).Error
}

func (r *positionRepository) GetAll() ([]models.Position, error) {
	var positions []models.Position
	err := r.db.Find(&positions).Error
	return positions, err
}

func (r *positionRepository) GetByID(id uint) (*models.Position, error) {
	var position models.Position
	err := r.db.First(&position, id).Error
	return &position, err
}

func (r *positionRepository) Update(position *models.Position) error {
	return r.db.Save(position).Error
}

func (r *positionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Position{}, id).Error
}

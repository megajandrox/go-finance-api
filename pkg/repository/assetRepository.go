package repository

import (
	"github.com/megajandrox/go-finance-api/pkg/models"
	"gorm.io/gorm"
)

type AssetRepository interface {
	Create(asset *models.Asset) error
	GetAll() ([]models.Asset, error)
	GetByID(id uint) (*models.Asset, error)
	UpdateByID(id uint, asset *models.Asset) error
	Delete(id uint) error
}

type assetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) AssetRepository {
	return &assetRepository{db}
}

func (r *assetRepository) Create(asset *models.Asset) error {
	return r.db.Create(asset).Error
}

func (r *assetRepository) GetAll() ([]models.Asset, error) {
	var assets []models.Asset
	err := r.db.Preload("Positions").Find(&assets).Error
	return assets, err
}

func (r *assetRepository) GetByID(id uint) (*models.Asset, error) {
	var asset models.Asset
	err := r.db.First(&asset, id).Error
	return &asset, err
}

func (r *assetRepository) UpdateByID(id uint, asset *models.Asset) error {
	return r.db.Model(&models.Asset{}).Where("id = ?", id).Updates(asset).Error
}

func (r *assetRepository) Delete(id uint) error {
	return r.db.Delete(&models.Asset{}, id).Error
}

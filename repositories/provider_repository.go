package repositories

import (
	"github.com/galiherlangga/clinic-appointment/models"
	"gorm.io/gorm"
)

type ProviderRepository interface {
	FindAll(filter models.ProviderFilter) ([]models.Provider, error)
	FindByID(id uint) (*models.Provider, error)
	Create(provider *models.Provider) error
	Update(provider *models.Provider) error
	Delete(id uint) error
	WithTx(tx *gorm.DB) ProviderRepository
}

type providerRepository struct {
	db *gorm.DB
}

func NewProviderRepository(db *gorm.DB) ProviderRepository {
	return &providerRepository{db: db}
}

func (r *providerRepository) FindAll(filter models.ProviderFilter) ([]models.Provider, error) {
	var providers []models.Provider

	err := r.applyFilter(r.db, filter).Find(&providers).Error
	if err != nil {
		return nil, err
	}
	return providers, nil
}

func (r *providerRepository) applyFilter(db *gorm.DB, filter models.ProviderFilter) *gorm.DB {
	if filter.Name != "" {
		db = db.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Email != "" {
		db = db.Where("email = ?", filter.Email)
	}
	if filter.Specialty != "" {
		db = db.Where("specialty = ?", filter.Specialty)
	}
	return db
}

func (r *providerRepository) FindByID(id uint) (*models.Provider, error) {
	var provider models.Provider
	err := r.db.First(&provider, id).Error
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

func (r *providerRepository) Create(provider *models.Provider) error {
	return r.db.Create(provider).Error
}

func (r *providerRepository) Update(provider *models.Provider) error {
	return r.db.Save(provider).Error
}

func (r *providerRepository) Delete(id uint) error {
	return r.db.Delete(&models.Provider{}, id).Error
}

func (r *providerRepository) WithTx(tx *gorm.DB) ProviderRepository {
	return &providerRepository{db: tx}
}

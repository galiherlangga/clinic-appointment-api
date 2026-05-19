package repositories

import (
	"github.com/galiherlangga/clinic-appointment/models"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	FindAll(filter models.CustomerFilter) ([]models.Customer, error)
	FindByID(id uint) (*models.Customer, error)
	Create(customer *models.Customer) error
	Update(customer *models.Customer) error
	Delete(id uint) error
	WithTx(tx *gorm.DB) CustomerRepository
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) FindAll(filter models.CustomerFilter) ([]models.Customer, error) {
	var customers []models.Customer

	err := r.applyFilter(r.db, filter).Find(&customers).Error
	if err != nil {
		return nil, err
	}
	return customers, nil
}

func (r *customerRepository) applyFilter(db *gorm.DB, filter models.CustomerFilter) *gorm.DB {
	if filter.Name != "" {
		db = db.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Email != "" {
		db = db.Where("email = ?", filter.Email)
	}
	if filter.Gender != "" {
		db = db.Where("gender = ?", filter.Gender)
	}
	return db
}

func (r *customerRepository) FindByID(id uint) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.First(&customer, id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) Create(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *customerRepository) Update(customer *models.Customer) error {
	return r.db.Save(customer).Error
}

func (r *customerRepository) Delete(id uint) error {
	return r.db.Delete(&models.Customer{}, id).Error
}

func (r *customerRepository) WithTx(tx *gorm.DB) CustomerRepository {
	return &customerRepository{db: tx}
}



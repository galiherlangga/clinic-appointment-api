package services

import (
	"context"
	"fmt"
	"time"

	"github.com/galiherlangga/clinic-appointment/configs"
	"github.com/galiherlangga/clinic-appointment/models"
	"github.com/galiherlangga/clinic-appointment/repositories"
	"gorm.io/gorm"
)

type CustomerService interface {
	FindAll(filter models.CustomerFilter) ([]models.Customer, error)
	FindByID(id uint) (*models.Customer, error)
	Create(customer *models.Customer) error
	Update(id uint, customer *models.Customer) error
	Delete(id uint) error
}

type customerService struct {
	db           *gorm.DB
	customerRepo repositories.CustomerRepository
	cacheService CacheService
}

func NewCustomerService(db *gorm.DB, customerRepo repositories.CustomerRepository, cacheService CacheService) CustomerService {
	return &customerService{
		db:           db,
		customerRepo: customerRepo,
		cacheService: cacheService,
	}
}

func (s *customerService) FindAll(filter models.CustomerFilter) ([]models.Customer, error) {
	cacheKey := s.getCacheKey(filter)
	var customers []models.Customer

	// 1. Try to fetch from Cache
	err := s.cacheService.Get(context.Background(), cacheKey, &customers)
	if err == nil {
		return customers, nil
	}

	// 2. Cache Miss: Fetch from repository
	customers, err = s.customerRepo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	// 3. Save to Cache (expire in 10 minutes)
	_ = s.cacheService.Set(context.Background(), cacheKey, customers, 10*time.Minute)

	return customers, nil
}

func (s *customerService) FindByID(id uint) (*models.Customer, error) {
	return s.customerRepo.FindByID(id)
}

func (s *customerService) Create(customer *models.Customer) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.customerRepo.WithTx(tx)

		if err := txRepo.Create(customer); err != nil {
			return err
		}

		s.clearCache()
		return nil
	})
}

func (s *customerService) Update(id uint, customer *models.Customer) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.customerRepo.WithTx(tx)

		existing, err := txRepo.FindByID(id)
		if err != nil {
			return err
		}

		existing.Name = customer.Name
		existing.Email = customer.Email
		existing.Phone = customer.Phone
		existing.Gender = customer.Gender
		existing.DateOfBirth = customer.DateOfBirth
		existing.Address = customer.Address

		if err := txRepo.Update(existing); err != nil {
			return err
		}

		s.clearCache()
		return nil
	})
}

func (s *customerService) Delete(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.customerRepo.WithTx(tx)

		if _, err := txRepo.FindByID(id); err != nil {
			return err
		}

		if err := txRepo.Delete(id); err != nil {
			return err
		}

		s.clearCache()
		return nil
	})
}

func (s *customerService) getCacheKey(filter models.CustomerFilter) string {
	prefix := configs.AppConfig.RedisPrefix
	return fmt.Sprintf("%scustomers:all:name=%s:email=%s:gender=%s", prefix, filter.Name, filter.Email, filter.Gender)
}

func (s *customerService) clearCache() {
	prefix := configs.AppConfig.RedisPrefix
	pattern := prefix + "customers:all:*"
	_ = s.cacheService.DeletePattern(context.Background(), pattern)
}

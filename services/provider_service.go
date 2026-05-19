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

type ProviderService interface {
	GetAll(filter models.ProviderFilter) ([]models.Provider, error)
	GetByID(id uint) (*models.Provider, error)
	Create(provider *models.Provider) error
	Update(id uint, provider *models.Provider) error
	Delete(id uint) error
}

type providerService struct {
	db           *gorm.DB
	providerRepo repositories.ProviderRepository
	cacheService CacheService
}

func NewProviderService(
	db *gorm.DB,
	providerRepo repositories.ProviderRepository,
	cacheService CacheService,
) ProviderService {
	return &providerService{
		db:           db,
		providerRepo: providerRepo,
		cacheService: cacheService,
	}
}

func (s *providerService) GetAll(filter models.ProviderFilter) ([]models.Provider, error) {
	cacheKey := s.getCacheKey(filter)
	var providers []models.Provider

	// 1. Try to fetch from Redis Cache
	err := s.cacheService.Get(context.Background(), cacheKey, &providers)
	if err == nil {
		return providers, nil
	}

	// 2. Cache Miss: Fetch from database repository
	providers, err = s.providerRepo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	// 3. Save to Redis Cache (expire in 1 hour)
	_ = s.cacheService.Set(context.Background(), cacheKey, providers, 1*time.Hour)

	return providers, nil
}

func (s *providerService) GetByID(id uint) (*models.Provider, error) {
	cacheKey := fmt.Sprintf("%sprovider:%d", configs.AppConfig.RedisPrefix, id)
	var provider models.Provider

	// 1. Try to fetch from Redis Cache
	err := s.cacheService.Get(context.Background(), cacheKey, &provider)
	if err == nil {
		return &provider, nil
	}

	// 2. Cache Miss: Fetch from database repository
	res, err := s.providerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 3. Save to Redis Cache (expire in 1 hour)
	_ = s.cacheService.Set(context.Background(), cacheKey, res, 1*time.Hour)

	return res, nil
}

func (s *providerService) Create(provider *models.Provider) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.providerRepo.WithTx(tx)

		if err := txRepo.Create(provider); err != nil {
			return err
		}

		s.clearCache(0)
		return nil
	})
}

func (s *providerService) Update(id uint, provider *models.Provider) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.providerRepo.WithTx(tx)

		existing, err := txRepo.FindByID(id)
		if err != nil {
			return err
		}

		existing.Name = provider.Name
		existing.Email = provider.Email
		existing.Phone = provider.Phone
		existing.Specialty = provider.Specialty

		if err := txRepo.Update(existing); err != nil {
			return err
		}

		s.clearCache(id)
		return nil
	})
}

func (s *providerService) Delete(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.providerRepo.WithTx(tx)

		if _, err := txRepo.FindByID(id); err != nil {
			return err
		}

		if err := txRepo.Delete(id); err != nil {
			return err
		}

		s.clearCache(id)
		return nil
	})
}

// Helper: Generate dynamic cache key based on query filters
func (s *providerService) getCacheKey(filter models.ProviderFilter) string {
	prefix := configs.AppConfig.RedisPrefix
	return fmt.Sprintf("%sproviders:all:name=%s:email=%s:specialty=%s", prefix, filter.Name, filter.Email, filter.Specialty)
}

// Helper: Scan and delete all provider caches from Redis safely
func (s *providerService) clearCache(id uint) {
	prefix := configs.AppConfig.RedisPrefix
	
	// Clear the list cache
	_ = s.cacheService.DeletePattern(context.Background(), prefix+"providers:all:*")
	
	// Clear specific item cache
	if id > 0 {
		_ = s.cacheService.Delete(context.Background(), fmt.Sprintf("%sprovider:%d", prefix, id))
	}
}

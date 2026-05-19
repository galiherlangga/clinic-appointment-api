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

type AppointmentService interface {
	GetAll(filter models.AppointmentFilter) ([]models.Appointment, error)
	GetByID(id uint) (*models.Appointment, error)
	Create(appointment *models.Appointment) error
	Update(id uint, appointment *models.Appointment) error
	Delete(id uint) error
}

type appointmentService struct {
	db              *gorm.DB
	appointmentRepo repositories.AppointmentRepository
	cacheService    CacheService
}

func NewAppointmentService(
	db *gorm.DB,
	appointmentRepo repositories.AppointmentRepository,
	cacheService CacheService,
) AppointmentService {
	return &appointmentService{
		db:              db,
		appointmentRepo: appointmentRepo,
		cacheService:    cacheService,
	}
}

func (s *appointmentService) GetAll(filter models.AppointmentFilter) ([]models.Appointment, error) {
	cacheKey := s.getCacheKey(filter)
	var appointments []models.Appointment

	// 1. Try to fetch from Cache
	err := s.cacheService.Get(context.Background(), cacheKey, &appointments)
	if err == nil {
		return appointments, nil
	}

	// 2. Cache Miss: Fetch from repository
	appointments, err = s.appointmentRepo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	// 3. Save to Cache (expire in 10 minutes)
	_ = s.cacheService.Set(context.Background(), cacheKey, appointments, 10*time.Minute)

	return appointments, nil
}

func (s *appointmentService) GetByID(id uint) (*models.Appointment, error) {
	cacheKey := fmt.Sprintf("%sappointment:%d", configs.AppConfig.RedisPrefix, id)
	var appointment models.Appointment

	// 1. Try to fetch from Cache
	err := s.cacheService.Get(context.Background(), cacheKey, &appointment)
	if err == nil {
		return &appointment, nil
	}

	// 2. Cache Miss: Fetch from repository
	res, err := s.appointmentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 3. Save to Cache (expire in 10 minutes)
	_ = s.cacheService.Set(context.Background(), cacheKey, res, 10*time.Minute)

	return res, nil
}

func (s *appointmentService) Create(appointment *models.Appointment) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.appointmentRepo.WithTx(tx)

		if err := txRepo.Create(appointment); err != nil {
			return err
		}

		s.clearCache(0)
		return nil
	})
}

func (s *appointmentService) Update(id uint, appointment *models.Appointment) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.appointmentRepo.WithTx(tx)

		existing, err := txRepo.FindByID(id)
		if err != nil {
			return err
		}

		existing.CustomerID = appointment.CustomerID
		existing.ProviderID = appointment.ProviderID
		existing.ServiceID = appointment.ServiceID
		existing.AppointmentTime = appointment.AppointmentTime
		existing.Status = appointment.Status
		existing.Notes = appointment.Notes

		if err := txRepo.Update(existing); err != nil {
			return err
		}

		s.clearCache(id)
		return nil
	})
}

func (s *appointmentService) Delete(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.appointmentRepo.WithTx(tx)

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

func (s *appointmentService) getCacheKey(filter models.AppointmentFilter) string {
	prefix := configs.AppConfig.RedisPrefix
	return fmt.Sprintf("%sappointments:all:customer=%d:provider=%d:service=%d:status=%s",
		prefix, filter.CustomerID, filter.ProviderID, filter.ServiceID, filter.Status)
}

func (s *appointmentService) clearCache(id uint) {
	prefix := configs.AppConfig.RedisPrefix
	
	// Clear all appointment list caches
	_ = s.cacheService.DeletePattern(context.Background(), prefix+"appointments:all:*")

	// Clear specific item cache
	if id > 0 {
		_ = s.cacheService.Delete(context.Background(), fmt.Sprintf("%sappointment:%d", prefix, id))
	}
}

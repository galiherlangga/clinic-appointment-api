package repositories

import (
	"github.com/galiherlangga/clinic-appointment/models"
	"gorm.io/gorm"
)

type AppointmentRepository interface {
	FindAll(filter models.AppointmentFilter) ([]models.Appointment, error)
	FindByID(id uint) (*models.Appointment, error)
	Create(appointment *models.Appointment) error
	Update(appointment *models.Appointment) error
	Delete(id uint) error
	WithTx(tx *gorm.DB) AppointmentRepository
}

type appointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &appointmentRepository{db: db}
}

func (r *appointmentRepository) FindAll(filter models.AppointmentFilter) ([]models.Appointment, error) {
	var appointments []models.Appointment

	err := r.applyFilter(r.db, filter).
		Preload("Customer").
		Preload("Provider").
		Preload("Service").
		Find(&appointments).Error
	if err != nil {
		return nil, err
	}
	return appointments, nil
}

func (r *appointmentRepository) applyFilter(db *gorm.DB, filter models.AppointmentFilter) *gorm.DB {
	if filter.CustomerID != 0 {
		db = db.Where("customer_id = ?", filter.CustomerID)
	}
	if filter.ProviderID != 0 {
		db = db.Where("provider_id = ?", filter.ProviderID)
	}
	if filter.ServiceID != 0 {
		db = db.Where("service_id = ?", filter.ServiceID)
	}
	if filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}
	return db
}

func (r *appointmentRepository) FindByID(id uint) (*models.Appointment, error) {
	var appointment models.Appointment
	err := r.db.
		Preload("Customer").
		Preload("Provider").
		Preload("Service").
		First(&appointment, id).Error
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (r *appointmentRepository) Create(appointment *models.Appointment) error {
	return r.db.Create(appointment).Error
}

func (r *appointmentRepository) Update(appointment *models.Appointment) error {
	return r.db.Save(appointment).Error
}

func (r *appointmentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Appointment{}, id).Error
}

func (r *appointmentRepository) WithTx(tx *gorm.DB) AppointmentRepository {
	return &appointmentRepository{db: tx}
}

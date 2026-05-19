package models

import (
	"time"

	"gorm.io/gorm"
)

type Appointment struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CustomerID      uint           `gorm:"not null;index" json:"customer_id"`
	Customer        *Customer      `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	ProviderID      uint           `gorm:"not null;index" json:"provider_id"`
	Provider        *Provider      `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
	ServiceID       uint           `gorm:"not null;index" json:"service_id"`
	Service         *Service       `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
	AppointmentTime time.Time      `gorm:"not null" json:"appointment_time"`
	Status          string         `gorm:"type:varchar(20);default:'PENDING'" json:"status"`
	Notes           string         `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type AppointmentFilter struct {
	CustomerID uint
	ProviderID uint
	ServiceID  uint
	Status     string
}


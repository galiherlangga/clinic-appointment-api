package models

import (
	"time"

	"gorm.io/gorm"
)

type Provider struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	Specialty    string         `gorm:"type:varchar(255);not null" json:"specialty"`
	Email        string         `gorm:"type:varchar(255);unique" json:"email"`
	Phone        string         `gorm:"type:varchar(50)" json:"phone"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Appointments []Appointment  `gorm:"foreignKey:ProviderID" json:"appointments,omitempty"`
}

type ProviderFilter struct {
	Name      string
	Specialty string
	Email     string
}

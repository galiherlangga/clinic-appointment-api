package models

import (
	"time"

	"gorm.io/gorm"
)

type Service struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	Duration     int            `gorm:"not null" json:"duration"` // in minutes
	Price        float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Appointments []Appointment  `gorm:"foreignKey:ServiceID" json:"appointments,omitempty"`
}


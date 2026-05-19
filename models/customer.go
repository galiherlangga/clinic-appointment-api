package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Phone        string         `gorm:"type:varchar(50)" json:"phone"`
	Gender       string         `gorm:"type:varchar(20)" json:"gender"`
	DateOfBirth  *time.Time     `gorm:"type:date" json:"date_of_birth"`
	Address      string         `gorm:"type:text" json:"address"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Appointments []Appointment  `gorm:"foreignKey:CustomerID" json:"appointments,omitempty"`
}

type CustomerFilter struct {
	Name   string
	Email  string
	Gender string
}



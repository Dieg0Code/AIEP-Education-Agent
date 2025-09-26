package models

import "gorm.io/gorm"

// Module (módulo académico que cursa el estudiante)
type Module struct {
	gorm.Model
	Code        string `json:"code" gorm:"type:varchar(50);uniqueIndex:ux_modules_code;not null"`
	Name        string `json:"name" gorm:"type:varchar(150);not null"`
	Description string `json:"description" gorm:"type:varchar(300)"`

	// Relaciones
	Topics      []Topic      `json:"topics,omitempty"`
	Enrollments []Enrollment `json:"enrollments,omitempty"`
}

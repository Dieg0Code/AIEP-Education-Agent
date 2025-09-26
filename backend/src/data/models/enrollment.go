package models

import "gorm.io/gorm"

// Enrollment (usuario inscrito en módulo)
type Enrollment struct {
	gorm.Model
	// Un usuario no debe estar inscrito dos veces al mismo módulo.
	// Para eso usamos la MISMA etiqueta uniqueIndex con el mismo nombre en ambos campos.
	UserID   uint   `json:"user_id" gorm:"index;uniqueIndex:ux_enrollment_user_module"`
	ModuleID uint   `json:"module_id" gorm:"index;uniqueIndex:ux_enrollment_user_module"`
	Status   string `json:"status" gorm:"type:varchar(20);default:'active'"` // active | dropped | completed

	// Relaciones
	User   User   `json:"user,omitzero"`
	Module Module `json:"module,omitzero"`
}

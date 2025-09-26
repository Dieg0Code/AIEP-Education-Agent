package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Topic (tema dentro del módulo)
type Topic struct {
	gorm.Model
	ModuleID      uint           `json:"module_id" gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ScheduledDate datatypes.Date `json:"scheduled_date" gorm:"type:date;not null;index"` // Fecha programada del tema

	UnitTitle         string `json:"unit_title" gorm:"type:varchar(200)"` // Título de la unidad
	OfficialContent   string `json:"official_content" gorm:"type:text"`   // Contenido oficial del tema
	ModernizedContent string `json:"modernized_content" gorm:"type:text"` // Contenido modernizado del tema

	// Relación
	Module Module `json:"module,omitzero"`
}

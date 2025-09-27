package models

import (
	"github.com/pgvector/pgvector-go"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Topic (tema dentro del módulo)
type Topic struct {
	gorm.Model
	ModuleID      uint           `json:"module_id" gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ScheduledDate datatypes.Date `json:"scheduled_date" gorm:"type:date;not null;index"` // Fecha programada del tema

	UnitTitle string          `json:"unit_title" gorm:"type:varchar(200)"` // Título de la unidad
	Content   string          `json:"content" gorm:"type:text"`            // Contenido del tema
	Embedding pgvector.Vector `json:"embedding" gorm:"type:vector(1536)"`  // Embedding para búsquedas vectoriales

	// Relación
	Module Module `json:"module,omitzero"`
}

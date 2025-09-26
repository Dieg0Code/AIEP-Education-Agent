package models

import (
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type Insight struct {
	gorm.Model
	UserID      uint            `json:"user_id" gorm:"index"`                        // Estudiante dueño
	InsightType string          `json:"insight_type" gorm:"type:varchar(100);index"` // Tipo de insight (e.g., "estilo_de_aprendizaje", "sesgo_cognitivo, "interes_academico", "habilidad_blanda", "problema_de_aprendizaje", "motivacion", etc.)
	Content     string          `json:"content" gorm:"type:text"`                    // Descripción o contenido del insight
	Embedding   pgvector.Vector `json:"embedding" gorm:"type:vector(1536)"`          // Embedding para búsquedas vectoriales
	// Relaciones
	User User `json:"user,omitzero"`
}

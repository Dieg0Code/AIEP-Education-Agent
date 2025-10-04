package topicdto

import (
	"time"

	"github.com/Dieg0Code/aiep-agent/src/data/models"
	"github.com/Dieg0Code/aiep-agent/src/pkg/date"
	"gorm.io/datatypes"
)

type TopicDTO struct {
	ID            uint    `json:"id" example:"1"`
	ModuleID      uint    `json:"module_id,omitempty" example:"2"`
	ModuleName    string  `json:"module_name,omitempty" example:"Go Basics"`
	UnitTitle     string  `json:"unit_title" example:"Introduction to Go"`
	Content       string  `json:"content" example:"This topic covers the basics of Go programming."`
	ScheduledDate string  `json:"scheduled_date,omitempty" example:"2023-10-01"` // Formato YYYY-MM-DD
	CreatedAt     string  `json:"created_at" example:"2023-09-01T12:00:00Z"`     // Formato RFC3339
	SimilarityPct float32 `json:"similarity_pct,omitempty" example:"0.85"`       // 0.0 a 1.0
}

// FromTopicModel convierte un models.Topic a TopicDTO (nil-safe), formateando fechas con pkg/date.
func FromTopicModel(t *models.Topic) TopicDTO {
	if t == nil {
		return TopicDTO{}
	}

	dto := TopicDTO{
		ID:         t.ID,
		UnitTitle:  t.UnitTitle,
		Content:    t.Content,
		ModuleID:   t.ModuleID,
		ModuleName: t.Module.Name,
	}

	if !t.CreatedAt.IsZero() {
		dto.CreatedAt = date.FormatDateTime(t.CreatedAt)
	}

	// datatypes.Date is a value type; convert to time.Time before checking
	if !time.Time(t.ScheduledDate).IsZero() {
		dto.ScheduledDate = date.FormatDate(time.Time(t.ScheduledDate))
	}

	return dto
}

// ParseScheduledDateToDatatypes parsea la cadena ScheduledDate del DTO a datatypes.Date.
// Si la cadena está vacía devuelve la fecha cero y nil error.
func (d *TopicDTO) ParseScheduledDateToDatatypes() (datatypes.Date, error) {
	if d == nil || d.ScheduledDate == "" {
		return datatypes.Date{}, nil
	}
	tm, err := date.ParseDate(d.ScheduledDate)
	if err != nil {
		return datatypes.Date{}, err
	}
	return datatypes.Date(tm), nil
}

package topicdto

import (
	"strings"
	"time"

	"github.com/Dieg0Code/aiep-agent/src/data/models"
	"github.com/Dieg0Code/aiep-agent/src/pkg/date"
)

// TopicFilter representa los filtros para consultar topics en el repositorio.
type TopicFilter struct {
	ModuleID uint
	Search   string
	Limit    int
	Offset   int
}

// ListTopicsRequestDTO representa los parámetros de consulta (query params).
type ListTopicsRequestDTO struct {
	ModuleID uint   `form:"module_id" json:"module_id" example:"1"`
	Search   string `form:"search" json:"search" example:"historia"`
	Limit    int    `form:"limit" json:"limit" binding:"omitempty,min=1,max=100" example:"20"`
	Offset   int    `form:"offset" json:"offset" example:"0"`
}

// Getters nil-safe / normalización
func (d *ListTopicsRequestDTO) GetModuleID() uint {
	if d == nil {
		return 0
	}
	return d.ModuleID
}

func (d *ListTopicsRequestDTO) GetSearch() string {
	if d == nil {
		return ""
	}
	return strings.TrimSpace(d.Search)
}

func (d *ListTopicsRequestDTO) GetLimit() int {
	if d == nil || d.Limit <= 0 {
		return 20 // default
	}
	if d.Limit > 100 {
		return 100 // cap
	}
	return d.Limit
}

func (d *ListTopicsRequestDTO) GetOffset() int {
	if d == nil || d.Offset < 0 {
		return 0
	}
	return d.Offset
}

// ToRepoFilter convierte el DTO al filtro esperado por el repo.
func (d *ListTopicsRequestDTO) ToRepoFilter() TopicFilter {
	return TopicFilter{
		ModuleID: d.GetModuleID(),
		Search:   d.GetSearch(),
		Limit:    d.GetLimit(),
		Offset:   d.GetOffset(),
	}
}

// TopicListItemDTO representa un resumen de Topic para la lista.
type TopicListItemDTO struct {
	ID            uint   `json:"id"`
	ModuleID      uint   `json:"module_id,omitempty"`
	UnitTitle     string `json:"unit_title"`
	Content       string `json:"content,omitempty"`
	ScheduledDate string `json:"scheduled_date,omitempty"` // ISO date "2006-01-02"
	CreatedAt     string `json:"created_at,omitempty"`     // RFC3339
}

// FromModel convierte models.Topic a TopicListItemDTO (nil-safe).
func FromModel(t *models.Topic) TopicListItemDTO {
	if t == nil {
		return TopicListItemDTO{}
	}

	topicList := TopicListItemDTO{
		ID:        t.ID,
		UnitTitle: t.UnitTitle,
		Content:   t.Content,
		ModuleID:  t.ModuleID,
	}

	// campos temporales/fechas: serializar de forma segura si existen
	if !t.CreatedAt.IsZero() {
		topicList.CreatedAt = date.FormatDateTime(t.CreatedAt)
	}

	if !time.Time(t.ScheduledDate).IsZero() {
		topicList.ScheduledDate = date.FormatDate(time.Time(t.ScheduledDate))
	}

	return topicList
}

// ListTopicsResponseDTO envuelve la respuesta paginada.
type ListTopicsResponseDTO struct {
	Items  []TopicListItemDTO `json:"items"`
	Limit  int                `json:"limit,omitempty"`
	Offset int                `json:"offset,omitempty"`
	// Total int64 `json:"total,omitempty"` // activa si el repo devuelve total
}

// MakeListResponse mapea la lista de modelos a la respuesta DTO.
func MakeListResponse(topics []models.Topic, req *ListTopicsRequestDTO) ListTopicsResponseDTO {
	items := make([]TopicListItemDTO, 0, len(topics))
	for i := range topics {
		t := topics[i]
		items = append(items, FromModel(&t))
	}
	return ListTopicsResponseDTO{
		Items:  items,
		Limit:  req.GetLimit(),
		Offset: req.GetOffset(),
	}
}

package topicrepo

import (
	"context"

	"github.com/Dieg0Code/aiep-agent/src/data/models"
	"gorm.io/datatypes"
)

// Lectura de temas
type TopicReader interface {
	TopicByID(ctx context.Context, id uint) (*models.Topic, error)
	ListTopics(ctx context.Context, filter TopicFilter) ([]models.Topic, error)
	TopicsByModule(ctx context.Context, moduleID uint) ([]models.Topic, error)
	TopicWithModule(ctx context.Context, topicID uint) (*models.Topic, error)                         // Con módulo incluido
	TopicsByDateRange(ctx context.Context, startDate, endDate datatypes.Date) ([]models.Topic, error) // Por rango de fechas
}

// Escritura de temas
type TopicWriter interface {
	CreateTopic(ctx context.Context, topic *models.Topic) (*models.Topic, error)
	UpdateTopic(ctx context.Context, id uint, updates TopicUpdate) error
	DeleteTopic(ctx context.Context, id uint) error
}

// Interfaz principal
type TopicRepo interface {
	TopicReader
	TopicWriter
}

// Filtro para temas
type TopicFilter struct {
	ModuleID uint   // Filtrar por módulo específico
	Search   string // Buscar en unit_title/official_content/modernized_content
	Limit    int
	Offset   int
}

// Actualización de tema
type TopicUpdate struct {
	UnitTitle         string
	OfficialContent   string
	ModernizedContent string
	ScheduledDate     *datatypes.Date // Pointer para permitir nil (no actualizar)
}

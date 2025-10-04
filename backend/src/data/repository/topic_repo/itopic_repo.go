package topicrepo

import (
	"context"

	topicdto "github.com/Dieg0Code/aiep-agent/src/data/dtos/topic_dto"
	"github.com/Dieg0Code/aiep-agent/src/data/models"
	"github.com/pgvector/pgvector-go"
	"gorm.io/datatypes"
)

// Lectura de temas
type TopicReader interface {
	TopicByID(ctx context.Context, id uint) (*models.Topic, error)
	ListTopics(ctx context.Context, filter TopicFilter) ([]models.Topic, error)
	TopicsByModule(ctx context.Context, moduleID uint) ([]models.Topic, error)
	TopicWithModule(ctx context.Context, topicID uint) (*models.Topic, error)                         // Con módulo incluido
	TopicsByDateRange(ctx context.Context, startDate, endDate datatypes.Date) ([]models.Topic, error) // Por rango de fechas

	// Búsquedas semánticas con embeddings
	SearchTopicsByEmbedding(ctx context.Context, embedding pgvector.Vector, limit int) ([]topicdto.VectorSearchResultDTO, error)
	SearchTopicsByEmbeddingWithFilter(ctx context.Context, embedding pgvector.Vector, filter SemanticFilter) ([]topicdto.VectorSearchResultDTO, error)
	FindSimilarTopics(ctx context.Context, topicID uint, limit int) ([]topicdto.VectorSearchResultDTO, error)
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

// Filtro para búsquedas semánticas
type SemanticFilter struct {
	ModuleID      uint    // Filtrar por módulo específico
	Limit         int     // Límite de resultados
	MinSimilarity float32 // Umbral mínimo de similitud (0.0 a 1.0)
}

// Actualización de tema
type TopicUpdate struct {
	UnitTitle     string
	Content       string
	ScheduledDate *datatypes.Date // Pointer para permitir nil (no actualizar)
}

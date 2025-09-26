package insightrepo

import (
	"context"

	"github.com/Dieg0Code/aiep-agent/src/data/models"
	"github.com/pgvector/pgvector-go"
)

// Lectura de insights
type InsightReader interface {
	InsightByID(ctx context.Context, id uint) (*models.Insight, error)
	ListInsights(ctx context.Context, filter InsightFilter) ([]models.Insight, error)
	InsightsByUser(ctx context.Context, userID uint) ([]models.Insight, error)
	InsightsByType(ctx context.Context, insightType string) ([]models.Insight, error)
	InsightsByUserAndType(ctx context.Context, userID uint, insightType string) ([]models.Insight, error)
	InsightWithUser(ctx context.Context, insightID uint) (*models.Insight, error) // Con User incluido

	// Búsquedas semánticas con embeddings
	SearchInsightsByEmbedding(ctx context.Context, embedding pgvector.Vector, limit int) ([]models.Insight, error)
	SearchInsightsByEmbeddingWithFilter(ctx context.Context, embedding pgvector.Vector, filter SemanticFilter) ([]models.Insight, error)
	FindSimilarInsights(ctx context.Context, insightID uint, limit int) ([]models.Insight, error)
}

// Escritura de insights
type InsightWriter interface {
	CreateInsight(ctx context.Context, insight *models.Insight) (*models.Insight, error)
	UpdateInsight(ctx context.Context, id uint, updates InsightUpdates) error
	DeleteInsight(ctx context.Context, id uint) error

	// Actualización de embeddings
	UpdateInsightEmbedding(ctx context.Context, id uint, embedding pgvector.Vector) error
	BatchUpdateEmbeddings(ctx context.Context, updates []EmbeddingUpdate) error
}

// Interfaz principal
type InsightRepo interface {
	InsightReader
	InsightWriter
}

// Filtro para insights tradicional
type InsightFilter struct {
	UserID      uint   // Filtrar por usuario específico
	InsightType string // Filtrar por tipo de insight
	Limit       int
	Offset      int
}

// Filtro para búsquedas semánticas
type SemanticFilter struct {
	UserID        uint    // Filtrar por usuario específico
	InsightType   string  // Filtrar por tipo de insight
	Limit         int     // Límite de resultados
	MinSimilarity float32 // Umbral mínimo de similitud (0.0 a 1.0)
}

// Estructura para actualizaciones parciales
type InsightUpdates struct {
	InsightType *string          // Tipo de insight
	Content     *string          // Contenido del insight
	Embedding   *pgvector.Vector // Embedding actualizado
}

// Estructura para actualizaciones batch de embeddings
type EmbeddingUpdate struct {
	ID        uint
	Embedding pgvector.Vector
}

// Constantes de tipos de insight
const (
	InsightTypeEstiloAprendizaje   = "estilo_de_aprendizaje"
	InsightTypeSesgoConitivo       = "sesgo_cognitivo"
	InsightTypeInteresAcademico    = "interes_academico"
	InsightTypeHabilidadBlanda     = "habilidad_blanda"
	InsightTypeProblemaAprendizaje = "problema_de_aprendizaje"
	InsightTypeMotivacion          = "motivacion"
	InsightTypePreferenciaHorario  = "preferencia_horario"
	InsightTypeMetodoEstudio       = "metodo_estudio"
	InsightTypeFortalezaAcademica  = "fortaleza_academica"
	InsightTypeAreaMejora          = "area_mejora"
)

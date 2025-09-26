package chatrepo

import (
	"context"

	"github.com/Dieg0Code/aiep-agent/src/data/models"
	"github.com/pgvector/pgvector-go"
)

// ============================================================================
// ChatSession Interfaces
// ============================================================================

// Lectura de sesiones de chat
type ChatSessionReader interface {
	ChatSessionByID(ctx context.Context, id uint) (*models.ChatSession, error)
	ChatSessionByUserID(ctx context.Context, userID uint) (*models.ChatSession, error)
	ListChatSessions(ctx context.Context, filter ChatSessionFilter) ([]models.ChatSession, error)
	ChatSessionWithMessages(ctx context.Context, sessionID uint) (*models.ChatSession, error) // Con Messages incluidos
	ChatSessionExists(ctx context.Context, userID uint) (bool, error)
}

// Escritura de sesiones de chat
type ChatSessionWriter interface {
	CreateChatSession(ctx context.Context, session *models.ChatSession) (*models.ChatSession, error)
	UpdateChatSession(ctx context.Context, id uint, updates ChatSessionUpdates) error
	DeleteChatSession(ctx context.Context, id uint) error
	DeleteChatSessionByUserID(ctx context.Context, userID uint) error
}

// ============================================================================
// ChatMessage Interfaces
// ============================================================================

// Lectura de mensajes de chat
type ChatMessageReader interface {
	ChatMessageByID(ctx context.Context, id uint) (*models.ChatMessage, error)
	ListChatMessages(ctx context.Context, filter ChatMessageFilter) ([]models.ChatMessage, error)
	ChatMessagesByConversationID(ctx context.Context, conversationID uint) ([]models.ChatMessage, error)
	ChatMessagesByRole(ctx context.Context, conversationID uint, role string) ([]models.ChatMessage, error)
	ChatMessageByToolCallID(ctx context.Context, toolCallID string) (*models.ChatMessage, error)
	GetConversationHistory(ctx context.Context, conversationID uint, limit int) ([]models.ChatMessage, error)

	// Búsquedas semánticas con embeddings
	SearchMessagesByEmbedding(ctx context.Context, embedding pgvector.Vector, limit int) ([]models.ChatMessage, error)
	SearchMessagesByEmbeddingWithFilter(ctx context.Context, embedding pgvector.Vector, filter SemanticMessageFilter) ([]models.ChatMessage, error)
	FindSimilarMessages(ctx context.Context, messageID uint, limit int) ([]models.ChatMessage, error)
	SearchMessagesByContent(ctx context.Context, query string, conversationID uint, limit int) ([]models.ChatMessage, error)
}

// Escritura de mensajes de chat
type ChatMessageWriter interface {
	CreateChatMessage(ctx context.Context, message *models.ChatMessage) (*models.ChatMessage, error)
	UpdateChatMessage(ctx context.Context, id uint, updates ChatMessageUpdates) error
	DeleteChatMessage(ctx context.Context, id uint) error
	DeleteMessagesByConversationID(ctx context.Context, conversationID uint) error

	// Operaciones en lote para mensajes
	BatchCreateMessages(ctx context.Context, messages []models.ChatMessage) ([]models.ChatMessage, error)
	BatchUpdateEmbeddings(ctx context.Context, updates []MessageEmbeddingUpdate) error

	// Actualización de embeddings
	UpdateMessageEmbedding(ctx context.Context, id uint, embedding pgvector.Vector) error
}

// ============================================================================
// Interfaces principales
// ============================================================================

// Interfaz principal para ChatSession
type ChatSessionRepo interface {
	ChatSessionReader
	ChatSessionWriter
}

// Interfaz principal para ChatMessage
type ChatMessageRepo interface {
	ChatMessageReader
	ChatMessageWriter
}

// Interfaz combinada para todo el chat
type ChatRepo interface {
	ChatSessionRepo
	ChatMessageRepo
}

// ============================================================================
// Filtros y Estructuras
// ============================================================================

// Filtro para sesiones de chat
type ChatSessionFilter struct {
	UserID    uint   // Filtrar por usuario específico
	AgentName string // Filtrar por nombre del agente
	Limit     int
	Offset    int
}

// Filtro para mensajes de chat
type ChatMessageFilter struct {
	ConversationID uint   // Filtrar por conversación específica
	Role           string // Filtrar por rol del mensaje
	ToolCallID     string // Filtrar por tool call ID
	Limit          int
}

// Filtro para búsquedas semánticas de mensajes
type SemanticMessageFilter struct {
	ConversationID uint    // Filtrar por conversación específica
	Role           string  // Filtrar por rol del mensaje
	Limit          int     // Límite de resultados
	MinSimilarity  float32 // Umbral mínimo de similitud (0.0 a 1.0)
}

// Estructura para actualizaciones parciales de ChatSession
type ChatSessionUpdates struct {
	UserName  *string // Nombre del usuario actualizado
	AgentName *string // Nombre del agente actualizado
}

// Estructura para actualizaciones parciales de ChatMessage
type ChatMessageUpdates struct {
	Role       *string          // Rol del mensaje
	Name       *string          // Nombre del emisor
	Content    *string          // Contenido del mensaje
	ToolCallID *string          // ID del tool call
	ToolCalls  any              // Tool calls (datatypes.JSON)
	Embedding  *pgvector.Vector // Embedding actualizado
}

// Estructura para actualizaciones batch de embeddings de mensajes
type MessageEmbeddingUpdate struct {
	ID        uint
	Embedding pgvector.Vector
}

// ============================================================================
// Constantes de roles de mensaje
// ============================================================================

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
	RoleTool      = "tool"
)

// ============================================================================
// Constantes de ordenamiento
// ============================================================================

const (
	OrderByCreatedAtASC  = "created_at ASC"
	OrderByCreatedAtDESC = "created_at DESC"
	OrderByUpdatedAtASC  = "updated_at ASC"
	OrderByUpdatedAtDESC = "updated_at DESC"
)

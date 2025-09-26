package chatrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Dieg0Code/aiep-agent/src/data/models"
	pgvector "github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type chatRepo struct {
	db *gorm.DB
}

func NewChatRepo(db *gorm.DB) (ChatRepo, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}

	return &chatRepo{
		db: db,
	}, nil
}

// BatchCreateMessages implements ChatRepo.
func (c *chatRepo) BatchCreateMessages(ctx context.Context, messages []models.ChatMessage) ([]models.ChatMessage, error) {
	if len(messages) == 0 {
		return []models.ChatMessage{}, nil
	}

	// Validar roles permitidos una sola vez
	validRoles := map[string]bool{
		RoleUser: true, RoleAssistant: true, RoleSystem: true, RoleTool: true,
	}

	// Validación de cada mensaje
	for _, message := range messages {
		if message.ConversationID == 0 {
			return nil, ErrInvalidConversationID
		}
		if message.Role == "" {
			return nil, ErrInvalidMessageRole
		}
		if !validRoles[message.Role] {
			return nil, ErrInvalidMessageRole
		}
		if message.Content == "" {
			return nil, ErrInvalidMessageContent
		}
		if len(message.Embedding.Slice()) > 0 && len(message.Embedding.Slice()) != 1536 {
			return nil, ErrInvalidEmbedding
		}
	}

	err := c.db.WithContext(ctx).Create(&messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// BatchUpdateEmbeddings implements ChatRepo.
func (c *chatRepo) BatchUpdateEmbeddings(ctx context.Context, updates []MessageEmbeddingUpdate) error {
	if len(updates) == 0 {
		return nil
	}

	// Validación de cada update
	for _, update := range updates {
		if update.ID == 0 {
			return ErrInvalidChatMessageID
		}
		if len(update.Embedding.Slice()) != 1536 {
			return ErrInvalidEmbedding
		}
	}

	// Actualizar cada embedding individualmente
	for _, update := range updates {
		result := c.db.WithContext(ctx).Model(&models.ChatMessage{}).
			Where("id = ?", update.ID).
			Update("embedding", update.Embedding)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return ErrChatMessageNotFound
		}
	}

	return nil
}

// ChatMessageByID implements ChatRepo.
func (c *chatRepo) ChatMessageByID(ctx context.Context, id uint) (*models.ChatMessage, error) {
	if id == 0 {
		return nil, ErrInvalidChatMessageID
	}

	var message models.ChatMessage
	err := c.db.WithContext(ctx).Where("id = ?", id).First(&message).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrChatMessageNotFound
		}
		return nil, err
	}

	return &message, nil
}

// ChatMessageByToolCallID implements ChatRepo.
func (c *chatRepo) ChatMessageByToolCallID(ctx context.Context, toolCallID string) (*models.ChatMessage, error) {
	if toolCallID == "" {
		return nil, ErrInvalidToolCallID
	}

	var message models.ChatMessage
	err := c.db.WithContext(ctx).Where("tool_call_id = ?", toolCallID).First(&message).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrChatMessageNotFound
		}
		return nil, err
	}

	return &message, nil
}

// ChatMessagesByConversationID implements ChatRepo.
func (c *chatRepo) ChatMessagesByConversationID(ctx context.Context, conversationID uint) ([]models.ChatMessage, error) {
	if conversationID == 0 {
		return nil, ErrInvalidConversationID
	}

	var messages []models.ChatMessage
	err := c.db.WithContext(ctx).Where("conversation_id = ?", conversationID).Order("created_at ASC").Find(&messages).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrChatMessageNotFound
		}
		return nil, err
	}
	return messages, nil
}

// ChatMessagesByRole implements ChatRepo.
func (c *chatRepo) ChatMessagesByRole(ctx context.Context, conversationID uint, role string) ([]models.ChatMessage, error) {
	if conversationID == 0 {
		return nil, ErrInvalidConversationID
	}
	if role == "" {
		return nil, ErrInvalidMessageRole
	}

	var messages []models.ChatMessage
	err := c.db.WithContext(ctx).
		Where("conversation_id = ? AND role = ?", conversationID, role).
		Order("created_at ASC").
		Find(&messages).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrChatMessageNotFound
		}
		return nil, err
	}
	return messages, nil
}

// ChatSessionByID implements ChatRepo.
func (c *chatRepo) ChatSessionByID(ctx context.Context, id uint) (*models.ChatSession, error) {
	if id == 0 {
		return nil, ErrInvalidChatSessionID
	}

	var session models.ChatSession
	err := c.db.WithContext(ctx).
		Where("id = ?", id).
		First(&session).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrChatSessionNotFound
		}
		return nil, err
	}
	return &session, nil
}

// ChatSessionByUserID implements ChatRepo.
func (c *chatRepo) ChatSessionByUserID(ctx context.Context, userID uint) (*models.ChatSession, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}

	var session models.ChatSession
	err := c.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&session).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrChatSessionNotFound
		}
	}

	return &session, nil
}

// ChatSessionExists implements ChatRepo.
func (c *chatRepo) ChatSessionExists(ctx context.Context, userID uint) (bool, error) {
	if userID == 0 {
		return false, ErrInvalidUserID
	}

	var count int64
	err := c.db.WithContext(ctx).
		Model(&models.ChatSession{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ChatSessionWithMessages implements ChatRepo.
func (c *chatRepo) ChatSessionWithMessages(ctx context.Context, sessionID uint) (*models.ChatSession, error) {
	if sessionID == 0 {
		return nil, ErrInvalidChatSessionID
	}

	var session models.ChatSession
	err := c.db.WithContext(ctx).
		Preload("Messages").
		First(&session, sessionID).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrChatSessionNotFound
		}
		return nil, err
	}

	return &session, nil
}

// CreateChatMessage implements ChatRepo.
func (c *chatRepo) CreateChatMessage(ctx context.Context, message *models.ChatMessage) (*models.ChatMessage, error) {
	if message == nil {
		return nil, ErrChatMessageNil
	}

	if message.ConversationID == 0 || message.Role == "" || message.Content == "" {
		return nil, ErrMissingRequiredFields
	}

	if len(message.Embedding.Slice()) > 0 && len(message.Embedding.Slice()) != 1536 {
		return nil, ErrInvalidEmbedding
	}

	// verificar la existencia del chat session
	var count int64
	err := c.db.WithContext(ctx).
		Model(&models.ChatSession{}).
		Where("id = ?", message.ConversationID).
		Count(&count).Error
	if err != nil {
		return nil, fmt.Errorf("error inesperado verificando la existencia del chat session: %w", err)
	}

	if count == 0 {
		return nil, ErrChatSessionNotFound
	}

	err = c.db.WithContext(ctx).Create(message).Error
	if err != nil {
		return nil, fmt.Errorf("error inesperado creando el chat message: %w", err)
	}

	return message, nil
}

// CreateChatSession implements ChatRepo.
func (c *chatRepo) CreateChatSession(ctx context.Context, session *models.ChatSession) (*models.ChatSession, error) {
	if session == nil {
		return nil, ErrChatSessionNil
	}

	if session.UserID == 0 || session.UserName == "" || session.AgentName == "" {
		return nil, ErrMissingRequiredFields
	}

	// verificar que no exista ya un chat session para el usuario
	var count int64
	err := c.db.WithContext(ctx).
		Model(&models.ChatSession{}).
		Where("user_id = ?", session.UserID).
		Count(&count).Error
	if err != nil {
		return nil, fmt.Errorf("error inesperado verificando la existencia del chat session: %w", err)
	}
	if count > 0 {
		return nil, ErrUserAlreadyHasSession
	}

	err = c.db.WithContext(ctx).Create(session).Error
	if err != nil {
		return nil, fmt.Errorf("error inesperado creando el chat session: %w", err)
	}

	return session, nil
}

// DeleteChatMessage implements ChatRepo.
func (c *chatRepo) DeleteChatMessage(ctx context.Context, id uint) error {
	if id == 0 {
		return ErrInvalidChatMessageID
	}

	result := c.db.WithContext(ctx).Delete(&models.ChatMessage{}, id)
	if result.Error != nil {
		return fmt.Errorf("error inesperado eliminando el chat message: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrChatMessageNotFound
	}

	return nil
}

// DeleteChatSession implements ChatRepo.
func (c *chatRepo) DeleteChatSession(ctx context.Context, id uint) error {
	if id == 0 {
		return ErrInvalidChatSessionID
	}

	result := c.db.WithContext(ctx).Delete(&models.ChatSession{}, id)
	if result.Error != nil {
		return fmt.Errorf("error inesperado eliminando el chat session: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrChatSessionNotFound
	}

	return nil
}

// DeleteChatSessionByUserID implements ChatRepo.
func (c *chatRepo) DeleteChatSessionByUserID(ctx context.Context, userID uint) error {
	if userID == 0 {
		return ErrInvalidUserID
	}

	// Verificar si existe el chat session
	var count int64
	err := c.db.WithContext(ctx).
		Model(&models.ChatSession{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("error inesperado verificando la existencia del chat session: %w", err)
	}
	if count == 0 {
		return ErrChatSessionNotFound
	}

	result := c.db.WithContext(ctx).Delete(&models.ChatSession{}, userID)
	if result.Error != nil {
		return fmt.Errorf("error inesperado eliminando el chat session: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrChatSessionNotFound
	}

	return nil
}

// DeleteMessagesByConversationID implements ChatRepo.
func (c *chatRepo) DeleteMessagesByConversationID(ctx context.Context, conversationID uint) error {
	if conversationID == 0 {
		return ErrInvalidConversationID
	}

	// Verificar si existen mensajes en la conversación
	var count int64
	err := c.db.WithContext(ctx).
		Model(&models.ChatMessage{}).
		Where("conversation_id = ?", conversationID).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("error inesperado verificando la existencia de mensajes en la conversación: %w", err)
	}

	if count == 0 {
		return ErrChatMessageNotFound
	}

	result := c.db.WithContext(ctx).Where("conversation_id = ?", conversationID).Delete(&models.ChatMessage{})
	if result.Error != nil {
		return fmt.Errorf("error inesperado eliminando los mensajes de la conversación: %w", result.Error)
	}

	return nil
}

// FindSimilarMessages implements ChatRepo.
func (c *chatRepo) FindSimilarMessages(ctx context.Context, messageID uint, limit int) ([]models.ChatMessage, error) {
	if messageID == 0 {
		return nil, ErrInvalidChatMessageID
	}

	if limit <= 0 {
		return nil, ErrInvalidLimit
	}

	var referenceMessage models.ChatMessage
	err := c.db.WithContext(ctx).First(&referenceMessage, messageID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrChatMessageNotFound
		}
		return nil, fmt.Errorf("error inesperado obteniendo el mensaje de referencia: %w", err)
	}

	if len(referenceMessage.Embedding.Slice()) == 0 {
		return nil, ErrChatMessageInvalid
	}

	if len(referenceMessage.Embedding.Slice()) != 1536 {
		return nil, ErrInvalidEmbedding
	}

	var similarMessages []models.ChatMessage
	result := c.db.WithContext(ctx).
		Model(&models.ChatMessage{}).
		Select("*, embedding <=> ? AS distance", referenceMessage.Embedding).
		Where("id != ?", messageID). // Excluir el mensaje de referencia
		Order("distance").
		Limit(limit).
		Find(&similarMessages)
	if result.Error != nil {
		return nil, fmt.Errorf("error inesperado buscando mensajes similares: %w", result.Error)
	}

	if len(similarMessages) == 0 {
		return nil, ErrNoSimilarMessagesFound
	}

	return similarMessages, nil
}

// GetConversationHistory implements ChatRepo.
func (c *chatRepo) GetConversationHistory(ctx context.Context, conversationID uint, limit int) ([]models.ChatMessage, error) {
	if conversationID == 0 {
		return nil, ErrInvalidConversationID
	}

	if limit <= 0 {
		return nil, ErrInvalidLimit
	}

	// Verificar si existen mensajes en la conversación
	var count int64
	err := c.db.WithContext(ctx).
		Model(&models.ChatMessage{}).
		Where("conversation_id = ?", conversationID).
		Count(&count).Error
	if err != nil {
		return nil, fmt.Errorf("error inesperado verificando la existencia de mensajes en la conversación: %w", err)
	}

	if count == 0 {
		return nil, ErrChatMessageNotFound
	}

	var messages []models.ChatMessage
	err = c.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, fmt.Errorf("error inesperado obteniendo el historial de la conversación: %w", err)
	}

	return messages, nil
}

// ListChatMessages implements ChatRepo.
func (c *chatRepo) ListChatMessages(ctx context.Context, filter ChatMessageFilter) ([]models.ChatMessage, error) {
	query := c.db.WithContext(ctx).Model(&models.ChatMessage{})

	// Aplicar filtros dinámicos
	if filter.ConversationID != 0 {
		query = query.Where("conversation_id = ?", filter.ConversationID)
	}

	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}

	if filter.ToolCallID != "" {
		query = query.Where("tool_call_id = ?", filter.ToolCallID)
	}

	query = query.Order("created_at desc")

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	var messages []models.ChatMessage
	err := query.Find(&messages).Error
	if err != nil {
		return nil, fmt.Errorf("error inesperado listando los mensajes de chat: %w", err)
	}

	// Invertir el slice para devolver en orden cronológico ascendente para el llm
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil

}

// ListChatSessions implements ChatRepo.
func (c *chatRepo) ListChatSessions(ctx context.Context, filter ChatSessionFilter) ([]models.ChatSession, error) {
	query := c.db.WithContext(ctx).Model(&models.ChatSession{})

	// Filtros dinámicos
	if filter.UserID != 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}

	if filter.AgentName != "" {
		query = query.Where("agent_name = ?", filter.AgentName)
	}

	// Paginación
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	// Ordenar por creación más reciente primero
	query = query.Order("created_at desc")

	var sessions []models.ChatSession
	if err := query.Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("error inesperado listando sesiones de chat: %w", err)
	}

	if len(sessions) == 0 {
		return nil, ErrChatSessionNotFound
	}

	return sessions, nil
}

// SearchMessagesByContent implements ChatRepo.
func (c *chatRepo) SearchMessagesByContent(
	ctx context.Context,
	text string,
	conversationID uint,
	limit int,
) ([]models.ChatMessage, error) {
	if text == "" {
		return nil, ErrInvalidSearchQuery
	}

	query := c.db.WithContext(ctx).Model(&models.ChatMessage{})

	// Filtrar por conversación si se especifica
	if conversationID != 0 {
		query = query.Where("conversation_id = ?", conversationID)
	}

	// Búsqueda textual (ILIKE para case-insensitive)
	query = query.Where("content ILIKE ?", "%"+text+"%")

	if limit > 0 {
		query = query.Limit(limit)
	}

	// Ordenar por más reciente
	query = query.Order("created_at desc")

	var messages []models.ChatMessage
	if err := query.Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("error buscando mensajes de chat: %w", err)
	}

	if len(messages) == 0 {
		return nil, ErrNoSimilarMessagesFound
	}

	return messages, nil
}

// SearchMessagesByEmbedding implements ChatRepo.
func (c *chatRepo) SearchMessagesByEmbedding(ctx context.Context, embedding pgvector.Vector, limit int) ([]models.ChatMessage, error) {
	if len(embedding.Slice()) != 1536 {
		return nil, ErrInvalidEmbedding
	}

	if limit <= 0 {
		return nil, ErrInvalidLimit
	}

	var messages []models.ChatMessage
	result := c.db.WithContext(ctx).
		Model(&models.ChatMessage{}).
		Select("*, embedding <=> ? AS distance", embedding).
		Order("distance").
		Limit(limit).
		Find(&messages)
	if result.Error != nil {
		return nil, fmt.Errorf("error inesperado buscando mensajes por embedding: %w", result.Error)
	}

	if len(messages) == 0 {
		return nil, ErrNoSimilarMessagesFound
	}

	return messages, nil
}

// SearchMessagesByEmbeddingWithFilter implements ChatRepo.
func (c *chatRepo) SearchMessagesByEmbeddingWithFilter(
	ctx context.Context,
	embedding pgvector.Vector,
	filter SemanticMessageFilter,
) ([]models.ChatMessage, error) {
	// Validación embedding
	if len(embedding.Slice()) != 1536 {
		return nil, ErrInvalidEmbedding
	}

	// Validación limit
	if filter.Limit <= 0 {
		return nil, ErrInvalidLimit
	}

	// Validación threshold
	if filter.MinSimilarity < 0.0 || filter.MinSimilarity > 1.0 {
		return nil, ErrSimilarityThreshold
	}

	// Construir query
	query := c.db.WithContext(ctx).
		Model(&models.ChatMessage{}).
		Select("*, embedding <=> ? AS distance", embedding)

	if filter.ConversationID != 0 {
		query = query.Where("conversation_id = ?", filter.ConversationID)
	}

	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}

	// El operador <=> retorna una distancia (0 = idéntico, mayor = menos parecido).
	// Si tienes un threshold en similitud (0.0–1.0), se transforma a distancia:
	// similitud = 1 - distancia  => distancia <= 1 - similitud
	if filter.MinSimilarity > 0 {
		maxDistance := 1.0 - filter.MinSimilarity
		query = query.Where("embedding <=> ? <= ?", embedding, maxDistance)
	}

	query = query.Order("distance").Limit(filter.Limit)

	var messages []models.ChatMessage
	if err := query.Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("error inesperado buscando mensajes por embedding con filtro: %w", err)
	}

	if len(messages) == 0 {
		return nil, ErrNoSimilarMessagesFound
	}

	return messages, nil
}

// UpdateChatMessage implements ChatRepo.
func (c *chatRepo) UpdateChatMessage(ctx context.Context, id uint, updates ChatMessageUpdates) error {
	panic("unimplemented")
}

// UpdateChatSession implements ChatRepo.
func (c *chatRepo) UpdateChatSession(ctx context.Context, id uint, updates ChatSessionUpdates) error {
	panic("unimplemented")
}

// UpdateMessageEmbedding implements ChatRepo.
func (c *chatRepo) UpdateMessageEmbedding(ctx context.Context, id uint, embedding pgvector.Vector) error {
	panic("unimplemented")
}

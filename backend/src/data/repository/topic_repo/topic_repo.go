package topicrepo

import (
	"context"
	"errors"

	topicdto "github.com/Dieg0Code/aiep-agent/src/data/dtos/topic_dto"
	"github.com/Dieg0Code/aiep-agent/src/data/models"
	pgvector "github.com/pgvector/pgvector-go"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type topicRepo struct {
	db *gorm.DB
}

func NewTopicRepo(db *gorm.DB) (TopicRepo, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}
	return &topicRepo{
		db: db,
	}, nil
}

// FindSimilarTopics implements TopicRepo.
func (t *topicRepo) FindSimilarTopics(ctx context.Context, topicID uint, limit int) ([]topicdto.VectorSearchResultDTO, error) {
	if topicID == 0 {
		return nil, ErrInvalidTopicID
	}

	if limit <= 0 {
		return nil, ErrInvalidLimit
	}

	// Obtener el topic de referencia con su embedding
	var referenceTopic models.Topic
	err := t.db.WithContext(ctx).First(&referenceTopic, topicID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTopicNotFound
		}
		return nil, err
	}

	// Verificar que el topic tiene embedding
	if len(referenceTopic.Embedding.Slice()) == 0 {
		return nil, ErrEmbeddingRequired
	}

	// Validar dimensiones del embedding de referencia
	if len(referenceTopic.Embedding.Slice()) != 1536 {
		return nil, ErrEmbeddingDimensions
	}

	// Buscar topics similares usando distancia coseno

	// Buscar topics similares directamente en el DTO
	var results []topicdto.VectorSearchResultDTO
	result := t.db.WithContext(ctx).
		Model(&models.Topic{}).
		Preload("Module").
		Select(`
		topics.id, 
		topics.scheduled_date, 
		topics.unit_title, 
		topics.content, 
		topics.module_id, 
		modules.name as module_name, 
		embedding <=> ? AS distance
		`, referenceTopic.Embedding).
		Joins("LEFT JOIN modules ON topics.module_id = modules.id").
		Where("topics.id != ?", topicID). // Excluir el topic de referencia
		Order("distance").
		Limit(limit).
		Scan(&results)

	if result.Error != nil {
		return nil, ErrSemanticSearchFailed
	}

	// Si no se encontraron topics similares
	if len(results) == 0 {
		return nil, ErrTopicNotFound
	}

	return results, nil
}

// SearchTopicsByEmbedding implements TopicRepo.
func (t *topicRepo) SearchTopicsByEmbedding(ctx context.Context, embedding pgvector.Vector, limit int) ([]topicdto.VectorSearchResultDTO, error) {
	if limit <= 0 {
		return nil, ErrInvalidLimit
	}

	if len(embedding.Slice()) != 1536 {
		return nil, ErrEmbeddingDimensions
	}

	var topics []topicdto.VectorSearchResultDTO
	result := t.db.WithContext(ctx).
		Model(&models.Topic{}).
		Select(`
		topics.id, 
		topics.scheduled_date, 
		topics.unit_title,
		topics.content,
		topics.module_id,
		modules.name as module_name,
		embedding <=> ? AS distance
		`, embedding).
		Joins("LEFT JOIN modules ON topics.module_id = modules.id").
		Order("distance").
		Limit(limit).
		Scan(&topics)

	if result.Error != nil {
		return nil, ErrSemanticSearchFailed
	}

	return topics, nil
}

// SearchTopicsByEmbeddingWithFilter implements TopicRepo.
func (t *topicRepo) SearchTopicsByEmbeddingWithFilter(ctx context.Context, embedding pgvector.Vector, filter SemanticFilter) ([]topicdto.VectorSearchResultDTO, error) {
	if filter.Limit <= 0 {
		return nil, ErrInvalidLimit
	}

	if len(embedding.Slice()) != 1536 {
		return nil, ErrEmbeddingDimensions
	}

	if filter.MinSimilarity < 0 || filter.MinSimilarity > 1.0 {
		return nil, ErrInvalidLimit
	}

	// Construir la consulta base con JOIN
	query := t.db.WithContext(ctx).
		Model(&models.Topic{}).
		Select(`
			topics.id, 
			topics.scheduled_date, 
			topics.unit_title,
			topics.content,
			topics.module_id,
			modules.name as module_name,
			topics.embedding <=> ? AS distance
		`, embedding).
		Joins("LEFT JOIN modules ON topics.module_id = modules.id")

	// Aplicar filtros
	if filter.ModuleID != 0 {
		query = query.Where("topics.module_id = ?", filter.ModuleID)
	}

	if filter.MinSimilarity > 0.0 {
		maxDistance := 1.0 - filter.MinSimilarity
		query = query.Where("topics.embedding <=> ? <= ?", embedding, maxDistance)
	}

	// Ejecutar consulta y escanear directamente al DTO
	var topics []topicdto.VectorSearchResultDTO
	result := query.Order("distance").Limit(filter.Limit).Scan(&topics)
	if result.Error != nil {
		return nil, ErrSemanticSearchFailed
	}

	return topics, nil
}

// CreateTopic implements TopicRepo.
func (t *topicRepo) CreateTopic(ctx context.Context, topic *models.Topic) (*models.Topic, error) {
	if topic == nil {
		return nil, ErrTopicNil
	}

	// Validaciones mínimas
	if topic.UnitTitle == "" || topic.ModuleID == 0 {
		return nil, ErrMissingRequiredFields
	}

	// Verificar que el módulo existe
	var count int64
	err := t.db.WithContext(ctx).Model(&models.Module{}).Where("id = ?", topic.ModuleID).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, ErrModuleNotExists
	}

	err = t.db.WithContext(ctx).Create(topic).Error
	if err != nil {
		return nil, err
	}

	return topic, nil
}

// DeleteTopic implements TopicRepo.
func (t *topicRepo) DeleteTopic(ctx context.Context, id uint) error {
	if id == 0 {
		return ErrInvalidTopicID
	}

	result := t.db.WithContext(ctx).Delete(&models.Topic{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrTopicNotFound
	}

	return nil
}

// ListTopics implements TopicRepo.
func (t *topicRepo) ListTopics(ctx context.Context, filter TopicFilter) ([]models.Topic, error) {
	query := t.db.WithContext(ctx).Model(&models.Topic{})

	// Filtrar por módulo específico
	if filter.ModuleID != 0 {
		query = query.Where("module_id = ?", filter.ModuleID)
	}

	// Búsqueda por texto en unit_title/official_content/modernized_content
	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("unit_title ILIKE ? OR content ILIKE ?", searchTerm, searchTerm)
	}

	// Aplicar limit y offset para paginación
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var topics []models.Topic
	err := query.Find(&topics).Error
	return topics, err
}

// TopicByID implements TopicRepo.
func (t *topicRepo) TopicByID(ctx context.Context, id uint) (*models.Topic, error) {
	if id == 0 {
		return nil, ErrInvalidTopicID
	}

	var topic models.Topic
	err := t.db.WithContext(ctx).Where("id = ?", id).First(&topic).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTopicNotFound
		}
		return nil, err
	}

	return &topic, nil
}

// TopicWithModule implements TopicRepo.
func (t *topicRepo) TopicWithModule(ctx context.Context, topicID uint) (*models.Topic, error) {
	if topicID == 0 {
		return nil, ErrInvalidTopicID
	}

	var topic models.Topic
	err := t.db.WithContext(ctx).Preload("Module").First(&topic, topicID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTopicNotFound
		}
		return nil, err
	}

	return &topic, nil
}

// TopicsByDateRange implements TopicRepo.
func (t *topicRepo) TopicsByDateRange(ctx context.Context, startDate datatypes.Date, endDate datatypes.Date) ([]models.Topic, error) {
	var topics []models.Topic
	err := t.db.WithContext(ctx).
		Where("scheduled_date >= ? AND scheduled_date <= ?", startDate, endDate).
		Order("scheduled_date ASC").
		Find(&topics).Error

	return topics, err
}

// TopicsByModule implements TopicRepo.
func (t *topicRepo) TopicsByModule(ctx context.Context, moduleID uint) ([]models.Topic, error) {
	if moduleID == 0 {
		return nil, ErrInvalidModuleID
	}

	var topics []models.Topic
	err := t.db.WithContext(ctx).
		Where("module_id = ?", moduleID).
		Order("scheduled_date ASC").
		Find(&topics).Error

	return topics, err
}

// UpdateTopic implements TopicRepo.
func (t *topicRepo) UpdateTopic(ctx context.Context, id uint, updates TopicUpdate) error {
	if id == 0 {
		return ErrInvalidTopicID
	}

	// Crear mapa de campos a actualizar (solo campos no vacíos)
	updateFields := make(map[string]any)

	if updates.UnitTitle != "" {
		updateFields["unit_title"] = updates.UnitTitle
	}
	if updates.Content != "" {
		updateFields["content"] = updates.Content
	}
	if updates.ScheduledDate != nil {
		updateFields["scheduled_date"] = *updates.ScheduledDate
	}

	// Si no hay campos para actualizar, no hacer nada
	if len(updateFields) == 0 {
		return nil
	}

	result := t.db.WithContext(ctx).Model(&models.Topic{}).Where("id = ?", id).Updates(updateFields)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrTopicNotFound
	}

	return nil
}

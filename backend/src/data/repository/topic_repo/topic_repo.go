package topicrepo

import (
	"context"
	"errors"

	"github.com/Dieg0Code/aiep-agent/src/data/models"
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
		query = query.Where("unit_title ILIKE ? OR official_content ILIKE ? OR modernized_content ILIKE ?", searchTerm, searchTerm, searchTerm)
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
	if updates.OfficialContent != "" {
		updateFields["official_content"] = updates.OfficialContent
	}
	if updates.ModernizedContent != "" {
		updateFields["modernized_content"] = updates.ModernizedContent
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

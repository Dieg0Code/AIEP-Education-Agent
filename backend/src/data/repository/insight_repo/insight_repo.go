package insightrepo

import (
	"context"

	"github.com/Dieg0Code/aiep-agent/src/data/models"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type insightRepo struct {
	db *gorm.DB
}

func NewInsightRepo(db *gorm.DB) (InsightRepo, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}
	return &insightRepo{
		db: db,
	}, nil
}

// CreateInsight implements InsightRepo.
func (i *insightRepo) CreateInsight(ctx context.Context, insight *models.Insight) (*models.Insight, error) {
	if insight == nil {
		return nil, ErrInsightNil
	}

	// Validaciones mínimas
	if insight.UserID == 0 || insight.InsightType == "" || insight.Content == "" {
		return nil, ErrMissingRequiredFields
	}

	// Validar que el tipo de insight sea válido (opcional: implementar lista permitida)
	if len(insight.InsightType) > 100 {
		return nil, ErrInvalidInsightType
	}

	// Verificar que el usuario existe
	var userCount int64
	err := i.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", insight.UserID).Count(&userCount).Error
	if err != nil {
		return nil, err
	}
	if userCount == 0 {
		return nil, ErrUserNotExists
	}

	err = i.db.WithContext(ctx).Create(insight).Error
	if err != nil {
		return nil, err
	}

	return insight, nil
}

// DeleteInsight implements InsightRepo.
func (i *insightRepo) DeleteInsight(ctx context.Context, id uint) error {
	if id == 0 {
		return ErrInvalidInsightID
	}

	result := i.db.WithContext(ctx).Delete(&models.Insight{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrInsightNotFound
	}

	return nil
}

// InsightByID implements InsightRepo.
func (i *insightRepo) InsightByID(ctx context.Context, id uint) (*models.Insight, error) {
	if id == 0 {
		return nil, ErrInvalidInsightID
	}

	var insight models.Insight
	err := i.db.WithContext(ctx).First(&insight, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrInsightNotFound
		}
		return nil, err
	}

	return &insight, nil
}

// InsightWithUser implements InsightRepo.
func (i *insightRepo) InsightWithUser(ctx context.Context, insightID uint) (*models.Insight, error) {
	if insightID == 0 {
		return nil, ErrInvalidInsightID
	}

	var insight models.Insight
	err := i.db.WithContext(ctx).
		Preload("User").
		First(&insight, insightID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrInsightNotFound
		}
		return nil, err
	}

	return &insight, nil
}

// InsightsByType implements InsightRepo.
func (i *insightRepo) InsightsByType(ctx context.Context, insightType string) ([]models.Insight, error) {
	if insightType == "" {
		return nil, ErrInvalidInsightType
	}

	var insights []models.Insight
	err := i.db.WithContext(ctx).
		Where("insight_type = ?", insightType).
		Order("created_at DESC").
		Find(&insights).Error
	if err != nil {
		return nil, err
	}

	return insights, nil
}

// InsightsByUser implements InsightRepo.
func (i *insightRepo) InsightsByUser(ctx context.Context, userID uint) ([]models.Insight, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}

	var insights []models.Insight
	err := i.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&insights).Error
	if err != nil {
		return nil, err
	}

	return insights, nil
}

// InsightsByUserAndType implements InsightRepo.
func (i *insightRepo) InsightsByUserAndType(ctx context.Context, userID uint, insightType string) ([]models.Insight, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}
	if insightType == "" {
		return nil, ErrInvalidInsightType
	}

	var insights []models.Insight
	err := i.db.WithContext(ctx).
		Where("user_id = ? AND insight_type = ?", userID, insightType).
		Order("created_at DESC").
		Find(&insights).Error
	if err != nil {
		return nil, err
	}

	return insights, nil
}

// ListInsights implements InsightRepo.
func (i *insightRepo) ListInsights(ctx context.Context, filter InsightFilter) ([]models.Insight, error) {
	query := i.db.WithContext(ctx).Model(&models.Insight{})

	// Aplicar filtros dinámicos
	if filter.UserID != 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.InsightType != "" {
		query = query.Where("insight_type = ?", filter.InsightType)
	}

	// Ordenamiento
	query = query.Order("created_at DESC")

	// Paginación
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var insights []models.Insight
	err := query.Find(&insights).Error
	if err != nil {
		return nil, err
	}

	return insights, nil
}

// UpdateInsight implements InsightRepo.
func (i *insightRepo) UpdateInsight(ctx context.Context, id uint, updates InsightUpdates) error {
	if id == 0 {
		return ErrInvalidInsightID
	}

	// Construir map de actualizaciones dinámicamente
	updateMap := make(map[string]interface{})

	// Solo actualizar campos que fueron especificados (no nil)
	if updates.InsightType != nil {
		if *updates.InsightType == "" {
			return ErrInvalidInsightType
		}
		if len(*updates.InsightType) > 100 {
			return ErrInvalidInsightType
		}
		updateMap["insight_type"] = *updates.InsightType
	}

	if updates.Content != nil {
		if *updates.Content == "" {
			return ErrEmptyContent
		}
		if len(*updates.Content) > 10000 {
			return ErrContentTooLong
		}
		updateMap["content"] = *updates.Content
	}

	// Si no hay nada que actualizar, no hacer nada
	if len(updateMap) == 0 {
		return nil
	}

	// Realizar la actualización
	result := i.db.WithContext(ctx).
		Model(&models.Insight{}).
		Where("id = ?", id).
		Updates(updateMap)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrInsightNotFound
	}

	return nil
}

// BatchUpdateEmbeddings implements InsightRepo.
func (i *insightRepo) BatchUpdateEmbeddings(ctx context.Context, updates []EmbeddingUpdate) error {
	if len(updates) == 0 {
		return ErrBatchUpdateEmpty
	}

	// Límite para evitar timeouts (ajustable según necesidades)
	if len(updates) > 1000 {
		return ErrBatchUpdateTooLarge
	}

	// Validar dimensiones de embeddings y IDs antes de procesar
	ids := make([]uint, 0, len(updates))
	for _, update := range updates {
		if update.ID == 0 {
			return ErrInvalidInsightID
		}
		// Validar que el embedding tenga las dimensiones correctas (1536 para OpenAI)
		if len(update.Embedding.Slice()) != 1536 {
			return ErrEmbeddingDimensions
		}
		ids = append(ids, update.ID)
	}

	// Ejecutar actualizaciones en una transacción
	err := i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Pre-verificar que todos los insights existen
		var existingCount int64
		err := tx.Model(&models.Insight{}).Where("id IN ?", ids).Count(&existingCount).Error
		if err != nil {
			return err
		}
		if int(existingCount) != len(updates) {
			return ErrInsightNotFound // Al menos uno de los IDs no existe
		}

		// Realizar actualizaciones usando CASE WHEN para mejor rendimiento
		for _, update := range updates {
			result := tx.Model(&models.Insight{}).
				Where("id = ?", update.ID).
				Update("embedding", update.Embedding)

			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})

	return err
}

// FindSimilarInsights implements InsightRepo.
func (i *insightRepo) FindSimilarInsights(ctx context.Context, insightID uint, limit int) ([]models.Insight, error) {
	if insightID == 0 {
		return nil, ErrInvalidInsightID
	}
	if limit <= 0 {
		return nil, ErrInvalidSearchLimit
	}

	// Primero obtener el insight de referencia con su embedding
	var referenceInsight models.Insight
	err := i.db.WithContext(ctx).First(&referenceInsight, insightID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrInsightNotFound
		}
		return nil, err
	}

	// Verificar que tenga embedding
	if len(referenceInsight.Embedding.Slice()) == 0 {
		return nil, ErrEmbeddingRequired
	}

	// Validar dimensiones del embedding de referencia
	if len(referenceInsight.Embedding.Slice()) != 1536 {
		return nil, ErrEmbeddingDimensions
	}

	// Buscar insights similares usando distancia coseno (más estándar para embeddings)
	var insights []models.Insight

	result := i.db.WithContext(ctx).
		Model(&models.Insight{}).
		Select("*, embedding <=> ? AS distance", referenceInsight.Embedding).
		Where("id != ?", insightID). // Excluir el insight de referencia
		Order("distance").           // Menor distancia = más similar
		Limit(limit).
		Find(&insights)

	if result.Error != nil {
		return nil, ErrSemanticSearchFailed
	}

	// Si no se encontraron resultados similares
	if len(insights) == 0 {
		return nil, ErrNoSimilarInsights
	}

	return insights, nil
}

// SearchInsightsByEmbedding implements InsightRepo.
func (i *insightRepo) SearchInsightsByEmbedding(ctx context.Context, embedding pgvector.Vector, limit int) ([]models.Insight, error) {
	if limit <= 0 {
		return nil, ErrInvalidSearchLimit
	}

	// Validar que el embedding tenga las dimensiones correctas
	if len(embedding.Slice()) != 1536 {
		return nil, ErrEmbeddingDimensions
	}

	// Buscar insights ordenados por similitud usando distancia coseno (estándar para embeddings)
	var insights []models.Insight

	result := i.db.WithContext(ctx).
		Model(&models.Insight{}).
		Select("*, embedding <=> ? AS distance", embedding).
		Order("distance"). // Menor distancia = más similar
		Limit(limit).
		Find(&insights)

	if result.Error != nil {
		return nil, ErrSemanticSearchFailed
	}

	return insights, nil
}

// SearchInsightsByEmbeddingWithFilter implements InsightRepo.
func (i *insightRepo) SearchInsightsByEmbeddingWithFilter(ctx context.Context, embedding pgvector.Vector, filter SemanticFilter) ([]models.Insight, error) {
	if filter.Limit <= 0 {
		return nil, ErrInvalidSearchLimit
	}

	// Validar que el embedding tenga las dimensiones correctas
	if len(embedding.Slice()) != 1536 {
		return nil, ErrEmbeddingDimensions
	}

	// Validar umbral de similitud si se especifica (debe ser válido cuando > 0.0)
	if filter.MinSimilarity < 0.0 || filter.MinSimilarity > 1.0 {
		return nil, ErrInvalidSimilarity
	}

	// Construir query base con búsqueda semántica usando distancia coseno (consistente con otras funciones)
	query := i.db.WithContext(ctx).
		Model(&models.Insight{}).
		Select("*, embedding <=> ? AS distance", embedding)

	// Aplicar filtros tradicionales
	if filter.UserID != 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.InsightType != "" {
		query = query.Where("insight_type = ?", filter.InsightType)
	}

	// Aplicar filtro de similitud mínima si se especifica
	// Nota: <=> retorna distancia coseno (0 = idénticos, 1 = diferentes)
	// Para MinSimilarity necesitamos convertir: distancia_maxima = 1.0 - MinSimilarity
	// Solo aplicar si MinSimilarity > 0.0 (0.0 significa "sin filtro de similitud")
	if filter.MinSimilarity > 0.0 {
		maxDistance := 1.0 - filter.MinSimilarity
		query = query.Where("embedding <=> ? <= ?", embedding, maxDistance)
	}

	// Ordenar por distancia (menor distancia = más similar)
	query = query.Order("distance").Limit(filter.Limit)

	var insights []models.Insight
	result := query.Find(&insights)

	if result.Error != nil {
		return nil, ErrSemanticSearchFailed
	}

	return insights, nil
}

// UpdateInsightEmbedding implements InsightRepo.
func (i *insightRepo) UpdateInsightEmbedding(ctx context.Context, id uint, embedding pgvector.Vector) error {
	if id == 0 {
		return ErrInvalidInsightID
	}

	// Validar que el embedding tenga las dimensiones correctas
	if len(embedding.Slice()) != 1536 {
		return ErrEmbeddingDimensions
	}

	// Actualizar el embedding del insight específico
	result := i.db.WithContext(ctx).
		Model(&models.Insight{}).
		Where("id = ?", id).
		Update("embedding", embedding)

	if result.Error != nil {
		return result.Error
	}

	// Verificar que el insight existe
	if result.RowsAffected == 0 {
		return ErrInsightNotFound
	}

	return nil
}

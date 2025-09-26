package insightrepo

import "errors"

var (
	// Errores de búsqueda
	ErrInsightNotFound = errors.New("insight no encontrado")

	// Errores de configuración
	ErrDatabaseRequired = errors.New("insight error: la conexión a la base de datos es requerida")

	// Errores de validación básica
	ErrInsightNil            = errors.New("insight error: el insight no puede ser nil")
	ErrInvalidInsightID      = errors.New("insight error: id de insight inválido")
	ErrInvalidUserID         = errors.New("insight error: id de usuario inválido")
	ErrMissingRequiredFields = errors.New("insight error: faltan campos requeridos: user_id/insight_type/content")

	// Errores de relación
	ErrUserNotExists = errors.New("insight error: el usuario especificado no existe")

	// Errores de negocio/validación de contenido
	ErrInvalidInsightType = errors.New("insight error: tipo de insight inválido")
	ErrEmptyContent       = errors.New("insight error: el contenido del insight no puede estar vacío")
	ErrContentTooLong     = errors.New("insight error: el contenido excede la longitud máxima permitida")

	// Errores de embeddings y búsquedas semánticas
	ErrInvalidEmbedding    = errors.New("insight error: embedding inválido o corrupto")
	ErrEmbeddingDimensions = errors.New("insight error: dimensiones del embedding incorrectas (debe ser 1536)")
	ErrEmbeddingRequired   = errors.New("insight error: embedding requerido para búsquedas semánticas")
	ErrInvalidSimilarity   = errors.New("insight error: umbral de similitud inválido (debe estar entre 0.0 y 1.0)")
	ErrBatchUpdateEmpty    = errors.New("insight error: lista de actualizaciones batch no puede estar vacía")
	ErrBatchUpdateTooLarge = errors.New("insight error: batch de actualizaciones excede el límite máximo")

	// Errores de búsquedas semánticas
	ErrSemanticSearchFailed = errors.New("insight error: fallo en búsqueda semántica")
	ErrNoSimilarInsights    = errors.New("insight error: no se encontraron insights similares")
	ErrInvalidSearchLimit   = errors.New("insight error: límite de búsqueda inválido (debe ser mayor a 0)")
)

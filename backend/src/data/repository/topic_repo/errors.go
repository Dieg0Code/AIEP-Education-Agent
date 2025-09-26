package topicrepo

import "errors"

var (
	// Errores de búsqueda
	ErrTopicNotFound = errors.New("tema no encontrado")

	// Errores de configuración
	ErrDatabaseRequired = errors.New("topic error: la conexión a la base de datos es requerida")

	// Errores de validación
	ErrTopicNil              = errors.New("topic error: el tema no puede ser nil")
	ErrInvalidTopicID        = errors.New("topic error: id de tema inválido")
	ErrInvalidModuleID       = errors.New("topic error: id de módulo inválido")
	ErrTitleEmpty            = errors.New("topic error: el título no puede estar vacío")
	ErrContentEmpty          = errors.New("topic error: el contenido no puede estar vacío")
	ErrMissingRequiredFields = errors.New("topic error: faltan campos requeridos: title/content/module_id")

	// Errores de relación
	ErrModuleNotExists = errors.New("topic error: el módulo especificado no existe")

	// Errores de negocio
	ErrInvalidScheduledDate  = errors.New("topic error: fecha programada inválida")
	ErrTopicAlreadyCompleted = errors.New("topic error: el tema ya está completado")
)

package modulerepo

import "errors"

var (
	// Errores de búsqueda
	ErrModuleNotFound = errors.New("módulo no encontrado")

	// Errores de configuración
	ErrDatabaseRequired = errors.New("module error: la conexión a la base de datos es requerida")

	// Errores de validación
	ErrModuleNil             = errors.New("module error: el módulo no puede ser nil")
	ErrInvalidModuleID       = errors.New("module error: id de módulo inválido")
	ErrCodeEmpty             = errors.New("module error: el código no puede estar vacío")
	ErrNameEmpty             = errors.New("module error: el nombre no puede estar vacío")
	ErrMissingRequiredFields = errors.New("module error: faltan campos requeridos: code/name")

	// Errores de conflicto/unicidad
	ErrModuleCodeConflict = errors.New("module error: el código del módulo ya está en uso")

	// Errores de negocio
	ErrInvalidModuleCode = errors.New("module error: código de módulo inválido: debe seguir el formato estándar")
)

package enrollementrepo

import "errors"

var (
	// Errores de búsqueda
	ErrEnrollmentNotFound = errors.New("inscripción no encontrada")

	// Errores de configuración
	ErrDatabaseRequired = errors.New("enrollment error: la conexión a la base de datos es requerida")

	// Errores de validación
	ErrEnrollmentNil         = errors.New("enrollment error: la inscripción no puede ser nil")
	ErrInvalidEnrollmentID   = errors.New("enrollment error: id de inscripción inválido")
	ErrInvalidUserID         = errors.New("enrollment error: id de usuario inválido")
	ErrInvalidModuleID       = errors.New("enrollment error: id de módulo inválido")
	ErrMissingRequiredFields = errors.New("enrollment error: faltan campos requeridos: user_id/module_id")

	// Errores de relación
	ErrUserNotExists   = errors.New("enrollment error: el usuario especificado no existe")
	ErrModuleNotExists = errors.New("enrollment error: el módulo especificado no existe")

	// Errores de unicidad/conflicto
	ErrUserAlreadyEnrolled = errors.New("enrollment error: el usuario ya está inscrito en este módulo")

	// Errores de negocio/estado
	ErrInvalidStatus            = errors.New("enrollment error: estado inválido (debe ser: active, dropped, completed)")
	ErrEnrollmentAlreadyDropped = errors.New("enrollment error: la inscripción ya está dada de baja")
	ErrEnrollmentAlreadyActive  = errors.New("enrollment error: la inscripción ya está activa")
	ErrCannotDropCompleted      = errors.New("enrollment error: no se puede dar de baja una inscripción completada")
)

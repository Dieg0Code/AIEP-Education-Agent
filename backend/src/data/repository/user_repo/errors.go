package userrepo

import "errors"

var (
	// Errores de búsqueda
	ErrUserNotFound = errors.New("usuario no encontrado")

	// Errores de configuración
	ErrDatabaseRequired = errors.New("user error: la conexión a la base de datos es requerida")

	// Errores de validación
	ErrUserNil               = errors.New("user error: el usuario no puede ser nil")
	ErrInvalidUserID         = errors.New("user error: id de usuario inválido")
	ErrEmailEmpty            = errors.New("user error: el email no puede estar vacío")
	ErrUsernameEmpty         = errors.New("user error: el nombre de usuario no puede estar vacío")
	ErrPasswordHashEmpty     = errors.New("user error: el hash de contraseña no puede estar vacío")
	ErrRoleEmpty             = errors.New("user error: el rol no puede estar vacío")
	ErrMissingRequiredFields = errors.New("user error: faltan campos requeridos: username/email/password")

	// Errores de conflicto/unicidad
	ErrUserNameConflict  = errors.New("user error: el nombre de usuario ya está en uso")
	ErrUserEmailConflict = errors.New("user error: el email ya está en uso")

	// Errores de negocio
	ErrInvalidRole = errors.New("user error: rol inválido: debe ser student, teacher o admin")
)

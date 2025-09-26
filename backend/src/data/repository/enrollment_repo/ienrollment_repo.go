package enrollementrepo

import (
	"context"

	"github.com/Dieg0Code/aiep-agent/src/data/models"
)

// Lectura de inscripciones
type EnrollmentReader interface {
	EnrollmentByID(ctx context.Context, id uint) (*models.Enrollment, error)
	ListEnrollments(ctx context.Context, filter EnrollmentFilter) ([]models.Enrollment, error)
	EnrollmentsByUser(ctx context.Context, userID uint) ([]models.Enrollment, error)
	EnrollmentsByModule(ctx context.Context, moduleID uint) ([]models.Enrollment, error)
	EnrollmentWithDetails(ctx context.Context, enrollmentID uint) (*models.Enrollment, error) // Con User y Module incluidos
	GetUserEnrollment(ctx context.Context, userID, moduleID uint) (*models.Enrollment, error) // Inscripción específica user-module
}

// Escritura de inscripciones
type EnrollmentWriter interface {
	CreateEnrollment(ctx context.Context, enrollment *models.Enrollment) (*models.Enrollment, error)
	UpdateEnrollmentStatus(ctx context.Context, id uint, status string) error
	DeleteEnrollment(ctx context.Context, id uint) error
}

// Interfaz principal
type EnrollmentRepo interface {
	EnrollmentReader
	EnrollmentWriter
}

// Filtro para inscripciones
type EnrollmentFilter struct {
	UserID   uint   // Filtrar por usuario específico
	ModuleID uint   // Filtrar por módulo específico
	Status   string // Filtrar por estado: active, dropped, completed
	Limit    int
	Offset   int
}

// Constantes de estado
const (
	StatusActive    = "active"
	StatusDropped   = "dropped"
	StatusCompleted = "completed"
)

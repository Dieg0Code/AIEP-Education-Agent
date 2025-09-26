package enrollementrepo

import (
	"context"
	"strings"

	"github.com/Dieg0Code/aiep-agent/src/data/models"
	"gorm.io/gorm"
)

type enrollmentRepo struct {
	db *gorm.DB
}

func NewEnrollmentRepo(db *gorm.DB) (EnrollmentRepo, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}
	return &enrollmentRepo{
		db: db,
	}, nil
}

// CreateEnrollment implements EnrollmentRepo.
func (e *enrollmentRepo) CreateEnrollment(ctx context.Context, enrollment *models.Enrollment) (*models.Enrollment, error) {
	if enrollment == nil {
		return nil, ErrEnrollmentNil
	}

	// Validaciones mínimas
	if enrollment.UserID == 0 || enrollment.ModuleID == 0 {
		return nil, ErrMissingRequiredFields
	}

	// Validar estado si está especificado
	if enrollment.Status != "" {
		if enrollment.Status != StatusActive && enrollment.Status != StatusDropped && enrollment.Status != StatusCompleted {
			return nil, ErrInvalidStatus
		}
	}

	// Verificar que el usuario existe
	var userCount int64
	err := e.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", enrollment.UserID).Count(&userCount).Error
	if err != nil {
		return nil, err
	}
	if userCount == 0 {
		return nil, ErrUserNotExists
	}

	// Verificar que el módulo existe
	var moduleCount int64
	err = e.db.WithContext(ctx).Model(&models.Module{}).Where("id = ?", enrollment.ModuleID).Count(&moduleCount).Error
	if err != nil {
		return nil, err
	}
	if moduleCount == 0 {
		return nil, ErrModuleNotExists
	}

	err = e.db.WithContext(ctx).Create(enrollment).Error
	if err != nil {
		// Detectar error de constraint único
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "ux_enrollment_user_module") {
			return nil, ErrUserAlreadyEnrolled
		}
		return nil, err
	}

	return enrollment, nil
}

// DeleteEnrollment implements EnrollmentRepo.
func (e *enrollmentRepo) DeleteEnrollment(ctx context.Context, id uint) error {
	if id == 0 {
		return ErrInvalidEnrollmentID
	}

	result := e.db.WithContext(ctx).Delete(&models.Enrollment{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrEnrollmentNotFound
	}

	return nil
}

// EnrollmentByID implements EnrollmentRepo.
func (e *enrollmentRepo) EnrollmentByID(ctx context.Context, id uint) (*models.Enrollment, error) {
	if id == 0 {
		return nil, ErrInvalidEnrollmentID
	}

	var enrollment models.Enrollment
	err := e.db.WithContext(ctx).First(&enrollment, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrEnrollmentNotFound
		}
		return nil, err
	}

	return &enrollment, nil
}

// EnrollmentWithDetails implements EnrollmentRepo.
func (e *enrollmentRepo) EnrollmentWithDetails(ctx context.Context, enrollmentID uint) (*models.Enrollment, error) {
	if enrollmentID == 0 {
		return nil, ErrInvalidEnrollmentID
	}

	var enrollment models.Enrollment
	err := e.db.WithContext(ctx).
		Preload("User").
		Preload("Module").
		First(&enrollment, enrollmentID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrEnrollmentNotFound
		}
		return nil, err
	}

	return &enrollment, nil
}

// EnrollmentsByModule implements EnrollmentRepo.
func (e *enrollmentRepo) EnrollmentsByModule(ctx context.Context, moduleID uint) ([]models.Enrollment, error) {
	if moduleID == 0 {
		return nil, ErrInvalidModuleID
	}

	var enrollments []models.Enrollment
	err := e.db.WithContext(ctx).
		Where("module_id = ?", moduleID).
		Order("created_at DESC").
		Find(&enrollments).Error
	if err != nil {
		return nil, err
	}

	return enrollments, nil
}

// EnrollmentsByUser implements EnrollmentRepo.
func (e *enrollmentRepo) EnrollmentsByUser(ctx context.Context, userID uint) ([]models.Enrollment, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}

	var enrollments []models.Enrollment
	err := e.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&enrollments).Error
	if err != nil {
		return nil, err
	}

	return enrollments, nil
}

// GetUserEnrollment implements EnrollmentRepo.
func (e *enrollmentRepo) GetUserEnrollment(ctx context.Context, userID uint, moduleID uint) (*models.Enrollment, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}
	if moduleID == 0 {
		return nil, ErrInvalidModuleID
	}

	var enrollment models.Enrollment
	err := e.db.WithContext(ctx).
		Where("user_id = ? AND module_id = ?", userID, moduleID).
		First(&enrollment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrEnrollmentNotFound
		}
		return nil, err
	}

	return &enrollment, nil
}

// ListEnrollments implements EnrollmentRepo.
func (e *enrollmentRepo) ListEnrollments(ctx context.Context, filter EnrollmentFilter) ([]models.Enrollment, error) {
	query := e.db.WithContext(ctx).Model(&models.Enrollment{})

	// Aplicar filtros dinámicos
	if filter.UserID != 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.ModuleID != 0 {
		query = query.Where("module_id = ?", filter.ModuleID)
	}
	if filter.Status != "" {
		// Validar estado antes de aplicar filtro
		if filter.Status != StatusActive && filter.Status != StatusDropped && filter.Status != StatusCompleted {
			return nil, ErrInvalidStatus
		}
		query = query.Where("status = ?", filter.Status)
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

	var enrollments []models.Enrollment
	err := query.Find(&enrollments).Error
	if err != nil {
		return nil, err
	}

	return enrollments, nil
}

// UpdateEnrollmentStatus implements EnrollmentRepo.
func (e *enrollmentRepo) UpdateEnrollmentStatus(ctx context.Context, id uint, status string) error {
	if id == 0 {
		return ErrInvalidEnrollmentID
	}

	// Validar estado
	if status != StatusActive && status != StatusDropped && status != StatusCompleted {
		return ErrInvalidStatus
	}

	// Obtener la inscripción actual para validar transiciones
	var currentEnrollment models.Enrollment
	err := e.db.WithContext(ctx).First(&currentEnrollment, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrEnrollmentNotFound
		}
		return err
	}

	// Validaciones de negocio para transiciones de estado
	if currentEnrollment.Status == status {
		// Evitar actualizar al mismo estado
		switch status {
		case StatusActive:
			return ErrEnrollmentAlreadyActive
		case StatusDropped:
			return ErrEnrollmentAlreadyDropped
		}
		// Para completed no hay error específico, simplemente no hacemos nada
		return nil
	}

	// No permitir cambiar de completed a dropped
	if currentEnrollment.Status == StatusCompleted && status == StatusDropped {
		return ErrCannotDropCompleted
	}

	// Actualizar el estado
	result := e.db.WithContext(ctx).
		Model(&models.Enrollment{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrEnrollmentNotFound
	}

	return nil
}

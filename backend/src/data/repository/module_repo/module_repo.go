package modulerepo

import (
	"context"
	"errors"
	"strings"

	"github.com/Dieg0Code/aiep-agent/src/data/models"
	"gorm.io/gorm"
)

type moduleRepo struct {
	db *gorm.DB
}

func NewModuleRepo(db *gorm.DB) (ModuleRepo, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}
	return &moduleRepo{
		db: db,
	}, nil
}

// CreateModule implements ModuleRepo.
func (m *moduleRepo) CreateModule(ctx context.Context, module *models.Module) (*models.Module, error) {
	if module == nil {
		return nil, ErrModuleNil
	}
	// Validaciones mínimas
	if module.Code == "" || module.Name == "" {
		return nil, ErrMissingRequiredFields
	}

	err := m.db.WithContext(ctx).Create(module).Error
	if err != nil {
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "ux_modules_code") {
			return nil, ErrModuleCodeConflict
		}
		return nil, err
	}
	return module, nil
}

// DeleteModule implements ModuleRepo.
func (m *moduleRepo) DeleteModule(ctx context.Context, id uint) error {
	if id == 0 {
		return ErrInvalidModuleID
	}

	result := m.db.WithContext(ctx).Delete(&models.Module{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrModuleNotFound
	}

	return nil
}

// ListModules implements ModuleRepo.
func (m *moduleRepo) ListModules(ctx context.Context, filter ModuleFilter) ([]models.Module, error) {
	query := m.db.WithContext(ctx).Model(&models.Module{})

	// Búsqueda por texto en code/name/description
	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("code ILIKE ? OR name ILIKE ? OR description ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	// Aplicar limit y offset para paginación
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var modules []models.Module
	err := query.Find(&modules).Error
	return modules, err
}

// ModuleWithEnrollments implements ModuleRepo.
func (m *moduleRepo) ModuleWithEnrollments(ctx context.Context, moduleID uint) (*models.Module, error) {
	if moduleID == 0 {
		return nil, ErrInvalidModuleID
	}

	var module models.Module
	err := m.db.WithContext(ctx).Preload("Enrollments").Preload("Enrollments.User").First(&module, moduleID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrModuleNotFound
		}
		return nil, err
	}

	return &module, nil
}

// ModuleByCode implements ModuleRepo.
func (m *moduleRepo) ModuleByCode(ctx context.Context, code string) (*models.Module, error) {
	if code == "" {
		return nil, ErrCodeEmpty
	}

	var module models.Module
	err := m.db.WithContext(ctx).Where("code = ?", code).First(&module).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrModuleNotFound
		}
		return nil, err
	}

	return &module, nil
}

// ModuleByID implements ModuleRepo.
func (m *moduleRepo) ModuleByID(ctx context.Context, id uint) (*models.Module, error) {
	if id == 0 {
		return nil, ErrInvalidModuleID
	}

	var module models.Module
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&module).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrModuleNotFound
		}
		return nil, err
	}

	return &module, nil
}

// ModuleWithTopics implements ModuleRepo.
func (m *moduleRepo) ModuleWithTopics(ctx context.Context, moduleID uint) (*models.Module, error) {
	if moduleID == 0 {
		return nil, ErrInvalidModuleID
	}

	var module models.Module
	err := m.db.WithContext(ctx).Preload("Topics").First(&module, moduleID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrModuleNotFound
		}
		return nil, err
	}

	return &module, nil
}

// UpdateModule implements ModuleRepo.
func (m *moduleRepo) UpdateModule(ctx context.Context, id uint, updates ModuleUpdate) error {
	if id == 0 {
		return ErrInvalidModuleID
	}

	// Crear mapa de campos a actualizar (solo campos no vacíos)
	updateFields := make(map[string]any)

	if updates.Name != "" {
		updateFields["name"] = updates.Name
	}
	if updates.Description != "" {
		updateFields["description"] = updates.Description
	}

	// Si no hay campos para actualizar, no hacer nada
	if len(updateFields) == 0 {
		return nil
	}

	result := m.db.WithContext(ctx).Model(&models.Module{}).Where("id = ?", id).Updates(updateFields)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrModuleNotFound
	}

	return nil
}

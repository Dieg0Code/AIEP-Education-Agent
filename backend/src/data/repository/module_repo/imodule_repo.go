package modulerepo

import (
	"context"

	"github.com/Dieg0Code/aiep-agent/src/data/models"
)

// Lectura de módulos
type ModuleReader interface {
	ModuleByID(ctx context.Context, id uint) (*models.Module, error)
	ModuleByCode(ctx context.Context, code string) (*models.Module, error)
	ListModules(ctx context.Context, filter ModuleFilter) ([]models.Module, error)
	ModuleWithTopics(ctx context.Context, moduleID uint) (*models.Module, error)      // Con temas incluidos
	ModuleWithEnrollments(ctx context.Context, moduleID uint) (*models.Module, error) // Con inscripciones incluidas
}

// Escritura de módulos
type ModuleWriter interface {
	CreateModule(ctx context.Context, module *models.Module) (*models.Module, error)
	UpdateModule(ctx context.Context, id uint, updates ModuleUpdate) error
	DeleteModule(ctx context.Context, id uint) error
}

// Interfaz principal
type ModuleRepo interface {
	ModuleReader
	ModuleWriter
}

// Filtro para módulos
type ModuleFilter struct {
	Search string // Buscar en code/name/description
	Limit  int
	Offset int
}

// Actualización de módulo
type ModuleUpdate struct {
	Name        string
	Description string
}

package userrepo

import (
	"context"

	"github.com/Dieg0Code/aiep-agent/src/data/models"
)

// Lectura
type UserReader interface {
	UserByID(ctx context.Context, id uint) (*models.User, error)
	UserByUsername(ctx context.Context, username string) (*models.User, error)
	UserByEmail(ctx context.Context, email string) (*models.User, error)
	ListUsers(ctx context.Context, filter UserFilter) ([]models.User, error)
}

// Escritura
type UserWriter interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdatePassword(ctx context.Context, id uint, newHash string) error
	UpdateRole(ctx context.Context, id uint, newRole string) error
	DeleteUser(ctx context.Context, id uint) error // soft delete (gorm)
}

type UserRepo interface {
	UserReader
	UserWriter
}

// Filtros simples
type UserFilter struct {
	Role   string // "" = todos, "student" = solo estudiantes
	Search string // "" = todos, "juan" = buscar "juan" en username/email
	Limit  int
	Offset int
}

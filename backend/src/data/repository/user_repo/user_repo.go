package userrepo

import (
	"context"
	"errors"
	"strings"

	"github.com/Dieg0Code/aiep-agent/src/data/models"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) (UserRepo, error) {
	if db == nil {
		return nil, ErrDatabaseRequired
	}
	return &userRepo{
		db: db,
	}, nil
}

// CreateUser implements UserRepo.
func (u *userRepo) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if user == nil {
		return nil, ErrUserNil
	}
	// Validaciones mínimas (el hashing debe venir hecho antes de llamar aquí)
	if user.UserName == "" || user.Email == "" || user.PasswordHash == "" {
		return nil, ErrMissingRequiredFields
	}

	err := u.db.WithContext(ctx).Create(user).Error
	if err != nil {
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "ux_users_username") {
			return nil, ErrUserNameConflict
		}
		if strings.Contains(errStr, "ux_users_email") {
			return nil, ErrUserEmailConflict
		}
		return nil, err
	}
	return user, nil
}

// DeleteUser implements UserRepo.
func (u *userRepo) DeleteUser(ctx context.Context, id uint) error {
	if id == 0 {
		return ErrInvalidUserID
	}

	result := u.db.WithContext(ctx).Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// ListUsers implements UserRepo.
func (u *userRepo) ListUsers(ctx context.Context, filter UserFilter) ([]models.User, error) {
	query := u.db.WithContext(ctx).Model(&models.User{})

	// Filtrar por rol si se especifica
	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}

	// Búsqueda por texto en username o email
	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("user_name ILIKE ? OR email ILIKE ?", searchTerm, searchTerm)
	}

	// Aplicar limit y offset para paginación
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var users []models.User
	err := query.Find(&users).Error
	return users, err
}

// UpdatePassword implements UserRepo.
func (u *userRepo) UpdatePassword(ctx context.Context, id uint, newHash string) error {
	if id == 0 {
		return ErrInvalidUserID
	}
	if newHash == "" {
		return ErrPasswordHashEmpty
	}

	result := u.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Update("password_hash", newHash)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// UpdateRole implements UserRepo.
func (u *userRepo) UpdateRole(ctx context.Context, id uint, newRole string) error {
	if id == 0 {
		return ErrInvalidUserID
	}
	if newRole == "" {
		return ErrRoleEmpty
	}

	// Validar roles permitidos
	validRoles := map[string]bool{
		"student": true,
		"teacher": true,
		"admin":   true,
	}
	if !validRoles[newRole] {
		return ErrInvalidRole
	}

	result := u.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Update("role", newRole)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// UserByEmail implements UserRepo.
func (u *userRepo) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, ErrEmailEmpty
	}

	var user models.User
	err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// UserByID implements UserRepo.
func (u *userRepo) UserByID(ctx context.Context, id uint) (*models.User, error) {
	if id == 0 {
		return nil, ErrInvalidUserID
	}

	var user models.User
	err := u.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// UserByUsername implements UserRepo.
func (u *userRepo) UserByUsername(ctx context.Context, username string) (*models.User, error) {
	if username == "" {
		return nil, ErrUsernameEmpty
	}

	var user models.User
	err := u.db.WithContext(ctx).Where("user_name = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

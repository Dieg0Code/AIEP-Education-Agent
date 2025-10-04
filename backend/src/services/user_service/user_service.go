package userservice

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Dieg0Code/aiep-agent/src/auth/bcrypt"
	userdto "github.com/Dieg0Code/aiep-agent/src/data/dtos/user_dto"
	userrepo "github.com/Dieg0Code/aiep-agent/src/data/repository/user_repo"
)

type userService struct {
	userRepo userrepo.UserRepo
	bcrypt   bcrypt.Bcrypt
	logger   *slog.Logger
}

// NewUserService crea una instancia de IUserService con el repositorio inyectado.
func NewUserService(userRepo userrepo.UserRepo, bcrypt bcrypt.Bcrypt, logger *slog.Logger) IUserService {
	return &userService{
		userRepo: userRepo,
		bcrypt:   bcrypt,
		logger:   logger,
	}
}

// Authenticate implements IUserService.
func (u *userService) Authenticate(ctx context.Context, req userdto.LoginRequestDTO) (userdto.UserDetailDTO, error) {

	if req.Email == "" && req.Password == "" {
		u.logger.ErrorContext(ctx, "Email and password must be provided",
			"email", req.Email,
			"password", req.Password,
		)
		return userdto.UserDetailDTO{}, fmt.Errorf("email and password must be provided")
	}

	user, err := u.userRepo.UserByEmail(ctx, req.Email)
	if err != nil {
		u.logger.ErrorContext(ctx, "Failed to get user by email during authentication",
			"error", err,
			"email", req.Email,
		)
		return userdto.UserDetailDTO{}, fmt.Errorf("failed to get user by email: %w", err)
	}

	if err := u.bcrypt.CompareHashAndPassword(user.PasswordHash, req.Password); err != nil {
		u.logger.ErrorContext(ctx, "Invalid password during authentication",
			"error", err,
			"email", req.Email,
		)
		return userdto.UserDetailDTO{}, fmt.Errorf("invalid password: %w", err)
	}

	u.logger.InfoContext(ctx, "User authenticated successfully",
		"user_id", user.ID,
		"email", user.Email,
	)

	return userdto.FromModelToDetail(user), nil
}

// CreateUser implements IUserService.
func (u *userService) CreateUser(ctx context.Context, req userdto.CreateUserDTO) (userdto.UserDetailDTO, error) {
	hashedPassword, err := u.bcrypt.HashPassword(req.Password)
	if err != nil {
		u.logger.ErrorContext(ctx, "Failed to hash password",
			"error", err,
			"username", req.UserName,
		)
		return userdto.UserDetailDTO{}, fmt.Errorf("failed to hash password: %w", err)
	}

	userModel := req.ToModelWithHash(hashedPassword)

	createdUser, err := u.userRepo.CreateUser(ctx, userModel)
	if err != nil {
		u.logger.ErrorContext(ctx, "Failed to create user in repository",
			"error", err,
			"username", req.UserName,
		)
		return userdto.UserDetailDTO{}, fmt.Errorf("failed to create user: %w", err)
	}

	u.logger.InfoContext(ctx, "User created successfully",
		"user_id", createdUser.ID,
		"username", createdUser.UserName,
		"role", createdUser.Role,
	)

	return userdto.FromModelToDetail(createdUser), nil
}

// DeleteUser implements IUserService.
func (u *userService) DeleteUser(ctx context.Context, id uint) error {

	if id == 0 {
		return fmt.Errorf("invalid user ID")
	}

	if err := u.userRepo.DeleteUser(ctx, id); err != nil {
		u.logger.ErrorContext(ctx, "Failed to delete user",
			"error", err,
			"user_id", id,
		)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	u.logger.InfoContext(ctx, "User deleted successfully", "user_id", id)
	return nil
}

// GetByEmail implements IUserService.
func (u *userService) GetByEmail(ctx context.Context, email string) (userdto.UserDetailDTO, error) {

	user, err := u.userRepo.UserByEmail(ctx, email)
	if err != nil {
		u.logger.ErrorContext(ctx, "Failed to get user by email",
			"error", err,
			"email", email,
		)
		return userdto.UserDetailDTO{}, fmt.Errorf("failed to get user by email: %w", err)
	}

	return userdto.FromModelToDetail(user), nil
}

// GetByID implements IUserService.
func (u *userService) GetByID(ctx context.Context, id uint) (userdto.UserDetailDTO, error) {
	if id == 0 {
		return userdto.UserDetailDTO{}, fmt.Errorf("invalid user ID")
	}

	user, err := u.userRepo.UserByID(ctx, id)
	if err != nil {
		u.logger.ErrorContext(ctx, "Failed to get user by ID",
			"error", err,
			"user_id", id,
		)
		return userdto.UserDetailDTO{}, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return userdto.FromModelToDetail(user), nil
}

// GetByUsername implements IUserService.
func (u *userService) GetByUsername(ctx context.Context, username string) (userdto.UserDetailDTO, error) {
	user, err := u.userRepo.UserByUsername(ctx, username)
	if err != nil {
		u.logger.ErrorContext(ctx, "Failed to get user by username",
			"error", err,
			"username", username,
		)
		return userdto.UserDetailDTO{}, fmt.Errorf("failed to get user by username: %w", err)
	}

	return userdto.FromModelToDetail(user), nil
}

// ListUsers implements IUserService.
func (s *userService) ListUsers(ctx context.Context, req userdto.ListUsersRequestDTO) (userdto.ListUsersResponseDTO, error) {
	// Validar y establecer valores por defecto
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	users, err := s.userRepo.ListUsers(ctx, req.ToRepoFilter())
	if err != nil {

		s.logger.ErrorContext(ctx, "Failed to list users",
			"error", err,
			"page_size", req.Limit,
			"search", req.Search,
			"role", req.Role,
		)
		return userdto.ListUsersResponseDTO{}, fmt.Errorf("failed to list users: %w", err)
	}

	// Usar la función helper del DTO
	response := userdto.MakeListUsersResponse(users, &req)

	return response, nil
}

// UpdatePassword implements IUserService.
func (u *userService) UpdatePassword(ctx context.Context, req userdto.UpdatePasswordRequestDTO) error {

	if req.UserID == 0 {
		u.logger.ErrorContext(ctx, "Invalid user ID",
			"user_id", req.UserID,
		)
		return fmt.Errorf("invalid user ID")
	}
	if req.OldPassword == "" {
		u.logger.ErrorContext(ctx, "Old password cannot be empty",
			"user_id", req.UserID,
		)
		return fmt.Errorf("old password cannot be empty")
	}
	if req.NewPassword == "" {
		u.logger.ErrorContext(ctx, "New password cannot be empty",
			"user_id", req.UserID,
		)
		return fmt.Errorf("new password cannot be empty")
	}

	// Validar contraseña actual
	user, err := u.userRepo.UserByID(ctx, req.UserID)
	if err != nil {
		u.logger.ErrorContext(ctx, "Failed to get user by ID",
			"error", err,
			"user_id", req.UserID,
		)
		return fmt.Errorf("failed to get user by ID: %w", err)
	}

	if err := u.bcrypt.CompareHashAndPassword(user.PasswordHash, req.OldPassword); err != nil {
		u.logger.ErrorContext(ctx, "Old password is incorrect",
			"user_id", req.UserID,
		)
		return fmt.Errorf("old password is incorrect: %w", err)
	}

	hashedPassword, err := u.bcrypt.HashPassword(req.NewPassword)
	if err != nil {
		u.logger.ErrorContext(ctx, "Failed to hash new password",
			"error", err,
			"user_id", req.UserID,
		)
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	if err := u.userRepo.UpdatePassword(ctx, req.UserID, hashedPassword); err != nil {
		u.logger.ErrorContext(ctx, "Failed to update user password",
			"error", err,
			"user_id", req.UserID,
		)
		return fmt.Errorf("failed to update user password: %w", err)
	}

	u.logger.InfoContext(ctx, "User password updated successfully", "user_id", req.UserID)
	return nil
}

// UpdateRole implements IUserService.
func (u *userService) UpdateRole(ctx context.Context, req userdto.UpdateRoleRequestDTO) error {

	if req.UserID == 0 {
		u.logger.ErrorContext(ctx, "Invalid user ID",
			"user_id", req.UserID,
		)
		return fmt.Errorf("invalid user ID")
	}
	if req.NewRole == "" {
		u.logger.ErrorContext(ctx, "New role cannot be empty",
			"user_id", req.UserID,
		)
		return fmt.Errorf("new role cannot be empty")
	}

	if err := u.userRepo.UpdateRole(ctx, req.UserID, req.NewRole); err != nil {
		u.logger.ErrorContext(ctx, "Failed to update user role",
			"error", err,
			"user_id", req.UserID,
			"new_role", req.NewRole,
		)
		return fmt.Errorf("failed to update user role: %w", err)
	}

	u.logger.InfoContext(ctx, "User role updated successfully", "user_id", req.UserID, "new_role", req.NewRole)
	return nil
}

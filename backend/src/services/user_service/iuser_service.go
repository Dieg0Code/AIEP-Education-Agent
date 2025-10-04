package userservice

import (
	"context"

	userdto "github.com/Dieg0Code/aiep-agent/src/data/dtos/user_dto"
)

// UserReader agrupa operaciones de lectura sobre usuarios.
type UserReader interface {
	GetByID(ctx context.Context, id uint) (userdto.UserDetailDTO, error)
	GetByUsername(ctx context.Context, username string) (userdto.UserDetailDTO, error)
	GetByEmail(ctx context.Context, email string) (userdto.UserDetailDTO, error)
	ListUsers(ctx context.Context, req userdto.ListUsersRequestDTO) (userdto.ListUsersResponseDTO, error)
}

// UserWriter agrupa operaciones de escritura/acciones sobre usuarios.
type UserWriter interface {
	CreateUser(ctx context.Context, req userdto.CreateUserDTO) (userdto.UserDetailDTO, error)
	UpdatePassword(ctx context.Context, req userdto.UpdatePasswordRequestDTO) error
	UpdateRole(ctx context.Context, req userdto.UpdateRoleRequestDTO) error
	DeleteUser(ctx context.Context, id uint) error
	Authenticate(ctx context.Context, req userdto.LoginRequestDTO) (userdto.UserDetailDTO, error)
}

// IUserService es la composición de lectura y escritura.
type IUserService interface {
	UserReader
	UserWriter
}

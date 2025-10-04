package userdto

import (
	"github.com/Dieg0Code/aiep-agent/src/data/models"
)

// CreateUserDTO represents the data required to create a new user.
// @Description CreateUserDTO is used for creating a new user in the system.
type CreateUserDTO struct {
	UserName string `json:"user_name" binding:"required,min=3,max=50" example:"juan123"`
	Password string `json:"password" binding:"required,min=6,max=100" example:"securePassword!"`
	Role     string `json:"role" binding:"required,oneof=student teacher admin" example:"student"`
	Email    string `json:"email" binding:"required,email" example:"juan@example.com"`
}

// ToModel convierte el DTO a un models.User sin incluir el PasswordHash.
// El service debe encargarse de hashear la contraseña antes de llamar a
// repo.CreateUser si se necesita guardar el hash.
func (d *CreateUserDTO) ToModel() *models.User {
	return &models.User{
		UserName: d.UserName,
		Role:     d.Role,
		Email:    d.Email,
	}
}

// ToModelWithHash convierte el DTO a models.User incluyendo el PasswordHash
// (útil cuando el service ya ha hasheado la contraseña).
func (d *CreateUserDTO) ToModelWithHash(hash string) *models.User {
	u := d.ToModel()
	u.PasswordHash = hash
	return u
}

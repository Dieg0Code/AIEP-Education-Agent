package userdto

import "strings"

// LoginRequestDTO represents the data required for user login.
// @Description LoginRequestDTO is used for authenticating a user.
type LoginRequestDTO struct {
	Email    string `json:"email" binding:"required,email" example:"juan@example.com"`
	Password string `json:"password" binding:"required,min=6,max=100" example:"securePassword!"`
}

// GetEmail devuelve el email contenido en el DTO (helper nil-safe).
func (d *LoginRequestDTO) GetEmail() string {
	if d == nil {
		return ""
	}
	return strings.TrimSpace(strings.ToLower(d.Email))
}

// GetPassword devuelve la contraseña contenida en el DTO (helper nil-safe).
func (d *LoginRequestDTO) GetPassword() string {
	if d == nil {
		return ""
	}
	return d.Password
}

package userdto

// UpdatePasswordRequestDTO represents the data required to update a user's password.
// @Description UpdatePasswordRequestDTO is used for changing a user's password.
type UpdatePasswordRequestDTO struct {
	UserID      uint   `json:"user_id" binding:"required" example:"1"`
	OldPassword string `json:"old_password" binding:"required,min=6,max=100" example:"oldPassword!"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=100" example:"newSecurePassword!"`
}

// GetUserID devuelve el id del usuario (helper nil-safe).
func (d *UpdatePasswordRequestDTO) GetUserID() uint {
	if d == nil {
		return 0
	}
	return d.UserID
}

// GetOldPassword devuelve la contraseña antigua (helper nil-safe).
func (d *UpdatePasswordRequestDTO) GetOldPassword() string {
	if d == nil {
		return ""
	}
	return d.OldPassword
}

// GetNewPassword devuelve la contraseña nueva (helper nil-safe).
func (d *UpdatePasswordRequestDTO) GetNewPassword() string {
	if d == nil {
		return ""
	}
	return d.NewPassword
}

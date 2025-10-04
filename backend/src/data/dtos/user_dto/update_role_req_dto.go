package userdto

import "strings"

// UpdateRoleRequestDTO represents the data required to update a user's role.
// @Description UpdateRoleRequestDTO is used for changing a user's role in the system.
type UpdateRoleRequestDTO struct {
	UserID  uint   `json:"user_id" binding:"required" example:"1"`
	NewRole string `json:"new_role" binding:"required,oneof=student teacher admin" example:"teacher"`
}

// GetUserID devuelve el ID del usuario (helper nil-safe).
func (d *UpdateRoleRequestDTO) GetUserID() uint {
	if d == nil {
		return 0
	}
	return d.UserID
}

// GetNewRole devuelve el nuevo rol solicitado (helper nil-safe).
func (d *UpdateRoleRequestDTO) GetNewRole() string {
	if d == nil {
		return ""
	}
	return strings.TrimSpace(strings.ToLower(d.NewRole))
}

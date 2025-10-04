package userdto

import (
	"github.com/Dieg0Code/aiep-agent/src/data/models"
	"github.com/Dieg0Code/aiep-agent/src/pkg/date"
)

// UserDetailDTO representa la vista pública de un usuario individual
// sin exponer campos sensibles como PasswordHash.
type UserDetailDTO struct {
	ID              uint   `json:"id"`
	UserName        string `json:"user_name"`
	Email           string `json:"email"`
	Role            string `json:"role"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at,omitempty"`
	Deleted         bool   `json:"deleted,omitempty"`
	ConversationID  *uint  `json:"conversation_id,omitempty"`
	EnrollmentCount int    `json:"enrollment_count,omitempty"`
	InsightsCount   int    `json:"insights_count,omitempty"`
}

// FromModel convierte models.User a UserDetailDTO (nil-safe).
func FromModelToDetail(u *models.User) UserDetailDTO {
	if u == nil {
		return UserDetailDTO{}
	}

	var convID *uint
	if u.Conversation != nil {
		id := u.Conversation.ID
		convID = &id
	}

	created := date.FormatDateTime(u.CreatedAt)
	var updated string
	if !u.UpdatedAt.IsZero() {
		updated = date.FormatDateTime(u.UpdatedAt)
	}

	return UserDetailDTO{
		ID:              u.ID,
		UserName:        u.UserName,
		Email:           u.Email,
		Role:            u.Role,
		CreatedAt:       created,
		UpdatedAt:       updated,
		Deleted:         u.DeletedAt.Valid,
		ConversationID:  convID,
		EnrollmentCount: len(u.Enrollments),
		InsightsCount:   len(u.Insights),
	}
}

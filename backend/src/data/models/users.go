package models

import "gorm.io/gorm"

// User representa a un estudiante, profesor o admin del sistema.
// Mantiene un único hilo de conversación asociado (Conversation) vía user_id en esa tabla.
type User struct {
	gorm.Model
	UserName     string `json:"user_name" gorm:"type:varchar(255);not null;uniqueIndex:ux_users_username"`
	PasswordHash string `json:"password" gorm:"type:varchar(255);not null"`
	Role         string `json:"role" gorm:"type:varchar(50);not null;default:'student'"` // student | teacher | admin
	Email        string `json:"email" gorm:"type:varchar(255);not null;uniqueIndex:ux_users_email"`

	// Relaciones
	Conversation *ChatSession `json:"conversation,omitempty"` // Un único hilo por usuario
	Enrollments  []Enrollment `json:"enrollments,omitempty"`  // Módulos en los que está inscrito
	Insights     []Insight    `json:"insights,omitempty"`
}

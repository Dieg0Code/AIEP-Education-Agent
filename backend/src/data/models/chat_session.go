package models

import (
	"gorm.io/gorm"
)

// Conversation: hilo entre un estudiante y un agente
type ChatSession struct {
	gorm.Model
	UserID    uint   `json:"user_id" gorm:"uniqueIndex"`                // Un único hilo por usuario
	UserName  string `json:"user_name" gorm:"type:varchar(255);index"`  // Cache opcional del nombre de usuario
	AgentName string `json:"agent_name" gorm:"type:varchar(255);index"` // Nombre del agente

	Messages []ChatMessage `json:"messages"`
}

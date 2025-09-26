package models

import (
	"github.com/pgvector/pgvector-go"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ChatMessage
type ChatMessage struct {
	gorm.Model
	ConversationID uint            `json:"conversation_id" gorm:"index"`
	Role           string          `json:"role" gorm:"type:varchar(20);index"` // user | assistant | system | tool
	Name           string          `json:"name" gorm:"type:varchar(255)"`
	Content        string          `json:"content" gorm:"type:text"`
	ToolCallID     string          `json:"tool_call_id" gorm:"type:varchar(255);index"`
	ToolCalls      datatypes.JSON  `json:"tool_calls" gorm:"type:jsonb"`
	Embedding      pgvector.Vector `json:"embedding" gorm:"type:vector(1536)"`
}

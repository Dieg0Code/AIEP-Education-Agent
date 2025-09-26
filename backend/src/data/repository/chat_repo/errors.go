package chatrepo

import "errors"

var (
	// Errores de búsqueda - ChatSession
	ErrChatSessionNotFound = errors.New("sesión de chat no encontrada")

	// Errores de configuración
	ErrDatabaseRequired = errors.New("chat error: la conexión a la base de datos es requerida")

	// Errores de validación - ChatSession
	ErrChatSessionNil        = errors.New("chat error: la sesión de chat no puede ser nil")
	ErrInvalidChatSessionID  = errors.New("chat error: id de sesión de chat inválido")
	ErrInvalidUserID         = errors.New("chat error: id de usuario inválido")
	ErrInvalidAgentName      = errors.New("chat error: nombre del agente inválido")
	ErrMissingRequiredFields = errors.New("chat error: faltan campos requeridos: user_id/agent_name")

	// Errores de relación - ChatSession
	ErrUserNotExists = errors.New("chat error: el usuario especificado no existe")

	// Errores de unicidad/conflicto - ChatSession
	ErrUserAlreadyHasSession = errors.New("chat error: el usuario ya tiene una sesión de chat activa")

	// Errores de búsqueda - ChatMessage
	ErrChatMessageNotFound = errors.New("mensaje de chat no encontrado")

	// Errores de validación - ChatMessage
	ErrChatMessageNil        = errors.New("chat error: el mensaje de chat no puede ser nil")
	ErrChatMessageInvalid    = errors.New("chat error: mensaje invalido")
	ErrInvalidChatMessageID  = errors.New("chat error: id de mensaje de chat inválido")
	ErrInvalidConversationID = errors.New("chat error: id de conversación inválido")
	ErrInvalidMessageRole    = errors.New("chat error: rol de mensaje inválido (debe ser: user, assistant, system, tool)")
	ErrInvalidMessageContent = errors.New("chat error: contenido del mensaje inválido")
	ErrInvalidToolCallID     = errors.New("chat error: id de tool call inválido")
	ErrInvalidToolCalls      = errors.New("chat error: tool calls inválidos")
	ErrInvalidEmbedding      = errors.New("chat error: embedding inválido (debe tener 1536 dimensiones)")

	// Errores de relación - ChatMessage
	ErrConversationNotExists = errors.New("chat error: la conversación especificada no existe")

	// Errores de negocio/estado
	ErrConversationHasMessages     = errors.New("chat error: la conversación tiene mensajes y no puede ser eliminada")
	ErrCannotDeleteSystemMessage   = errors.New("chat error: no se puede eliminar un mensaje del sistema")
	ErrCannotModifyArchivedMessage = errors.New("chat error: no se puede modificar un mensaje archivado")

	// Errores de búsqueda semántica
	ErrInvalidSearchQuery     = errors.New("chat error: consulta de búsqueda inválida")
	ErrInvalidLimit           = errors.New("chat error: límite inválido (debe ser > 0 y <= 100)")
	ErrInvalidOffset          = errors.New("chat error: offset inválido (debe ser >= 0)")
	ErrSimilarityThreshold    = errors.New("chat error: umbral de similitud inválido (debe ser entre 0.0 y 1.0)")
	ErrNoSimilarMessagesFound = errors.New("chat error: no se encontraron mensajes similares")
)

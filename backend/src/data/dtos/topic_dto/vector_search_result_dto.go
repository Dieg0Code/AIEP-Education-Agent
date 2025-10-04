package topicdto

// DTO para resultados de búsqueda vectorial

type VectorSearchResultDTO struct {
	ID            uint    `json:"id"`
	ScheduledDate string  `json:"scheduled_date"`
	UnitTitle     string  `json:"unit_title"`
	Content       string  `json:"content"`
	ModuleID      uint    `json:"module_id"`
	ModuleName    string  `json:"module_name"`
	Distance      float32 `json:"distance"` // Distancia del embedding
}

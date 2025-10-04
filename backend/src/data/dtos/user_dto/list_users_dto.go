package userdto

import (
	"github.com/Dieg0Code/aiep-agent/src/data/models"
	userrepo "github.com/Dieg0Code/aiep-agent/src/data/repository/user_repo"
	"github.com/Dieg0Code/aiep-agent/src/pkg/date"
)

// ListUsersRequestDTO representa los parámetros de consulta (query params).
type ListUsersRequestDTO struct {
	Role   string `form:"role" json:"role" example:"student"`
	Search string `form:"search" json:"search" example:"juan"`
	Limit  int    `form:"limit" json:"limit" binding:"omitempty,min=1,max=100" example:"20"`
	Offset int    `form:"offset" json:"offset" example:"0"`
}

// ToRepoFilter convierte el DTO al filtro esperado por el repo.
func (q *ListUsersRequestDTO) ToRepoFilter() userrepo.UserFilter {
	return userrepo.UserFilter{
		Role:   q.Role,
		Search: q.Search,
		Limit:  q.Limit,
		Offset: q.Offset,
		// Nota: si userrepo.UserFilter se expande (Order/IncludeDeleted), mapea aquí.
	}
}

// UserListItemDTO representa un usuario en la lista pública (resumen).
type UserListItemDTO struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	Deleted   bool   `json:"deleted,omitempty"`
}

// FromModel convierte models.User a UserListItemDTO.
func FromModel(u *models.User) UserListItemDTO {
	return UserListItemDTO{
		ID:        u.ID,
		Username:  u.UserName,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: date.FormatDateTime(u.CreatedAt),
		Deleted:   u.DeletedAt.Valid, // asumiendo gorm.DeletedAt en models.User
	}
}

// ListUsersResponseDTO envuelve la lista devuelta al cliente.
// Nota: el repo actual no devuelve total; si quieres paginación total, agrega total al repo.
type ListUsersResponseDTO struct {
	Items  []UserListItemDTO `json:"items"`
	Limit  int               `json:"limit,omitempty"`
	Offset int               `json:"offset,omitempty"`
	// Total int64 `json:"total,omitempty"` // activar si repo devuelve el total
}

// MakeListUsersResponse crea la respuesta a partir de modelos.
func MakeListUsersResponse(users []models.User, req *ListUsersRequestDTO) ListUsersResponseDTO {
	items := make([]UserListItemDTO, 0, len(users))
	for i := range users {
		u := users[i] // evitar &users[i] al iterar sobre el slice
		items = append(items, FromModel(&u))
	}
	resp := ListUsersResponseDTO{
		Items:  items,
		Limit:  req.Limit,
		Offset: req.Offset,
	}
	return resp
}

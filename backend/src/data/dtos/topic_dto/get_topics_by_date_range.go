package topicdto

import (
	"fmt"

	"github.com/Dieg0Code/aiep-agent/src/pkg/date"
	"gorm.io/datatypes"
)

// GetTopicByIDDTO represents the data required to retrieve a topic by ID.
// @Description GetTopicByIDDTO is used for fetching a topic by its unique identifier.
type GetByDateRangeDTO struct {
	StartDate string `json:"start_date" binding:"required,datetime=2006-01-02" example:"2023-01-01"`
	EndDate   string `json:"end_date" binding:"required,datetime=2006-01-02" example:"2023-12-31"`
}

// GetStartDate devuelve la fecha de inicio contenida en el DTO (helper nil-safe para servicios/repos).
func (d *GetByDateRangeDTO) GetStartDate() string {
	if d == nil {
		return ""
	}
	return d.StartDate
}

// GetEndDate devuelve la fecha de fin contenida en el DTO (helper nil-safe para servicios/repos).
func (d *GetByDateRangeDTO) GetEndDate() string {
	if d == nil {
		return ""
	}
	return d.EndDate
}

func (d *GetByDateRangeDTO) ParseToDatatypes() (datatypes.Date, datatypes.Date, error) {
	if d == nil {
		return datatypes.Date{}, datatypes.Date{}, fmt.Errorf("empty dto")
	}
	startStr := d.GetStartDate()
	endStr := d.GetEndDate()

	sd, err := date.ParseDate(startStr)
	if err != nil {
		return datatypes.Date{}, datatypes.Date{}, fmt.Errorf("invalid start_date: %w", err)
	}
	ed, err := date.ParseDate(endStr)
	if err != nil {
		return datatypes.Date{}, datatypes.Date{}, fmt.Errorf("invalid end_date: %w", err)
	}

	if sd.After(ed) {
		return datatypes.Date{}, datatypes.Date{}, fmt.Errorf("start_date must be <= end_date")
	}

	return datatypes.Date(sd), datatypes.Date(ed), nil
}

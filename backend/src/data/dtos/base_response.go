package dtos

// BaseResponse represents a standardized response structure for API responses.
// @Description BaseResponse is the standard response format for all API endpoints.
type BaseResponse struct {
	Code    int    `json:"code" example:"200" extension:"x-order=0"`
	Status  string `json:"status" example:"success" extension:"x-order=1"`
	Message string `json:"message" example:"Operation completed successfully" extension:"x-order=2"`
	Data    any    `json:"data,omitempty" extension:"x-order=3"`
}

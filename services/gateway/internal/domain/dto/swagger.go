package dto

// Used for better documentation in swagger
type ErrorResponse struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"Invalid input"`
	Error   string `json:"error" example:"Email is not valid"`
}

package dto

// StandardResponse represents the standard API response format
type StandardResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Timestamp string      `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// ErrorInfo represents error details in API responses
type ErrorInfo struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// Generic error responses grouped by HTTP status code
// These show the structure, not specific error messages

type BadRequestError struct {
	Success   bool      `json:"success" example:"false"`
	Message   string    `json:"message" example:"Please check the request parameters"`
	Error     ErrorInfo `json:"error"`
	Timestamp string    `json:"timestamp" example:"2025-01-27T10:30:00Z"`
	RequestID string    `json:"request_id,omitempty" example:"req_123"`
}

type UnauthorizedError struct {
	Success   bool      `json:"success" example:"false"`
	Message   string    `json:"message" example:"You are not authorized to access this resource"`
	Error     ErrorInfo `json:"error"`
	Timestamp string    `json:"timestamp" example:"2025-01-27T10:30:00Z"`
	RequestID string    `json:"request_id,omitempty" example:"req_123"`
}

type ForbiddenError struct {
	Success   bool      `json:"success" example:"false"`
	Message   string    `json:"message" example:"You are not allowed to access this resource"`
	Error     ErrorInfo `json:"error"`
	Timestamp string    `json:"timestamp" example:"2025-01-27T10:30:00Z"`
	RequestID string    `json:"request_id,omitempty" example:"req_123"`
}

type NotFoundError struct {
	Success   bool      `json:"success" example:"false"`
	Message   string    `json:"message" example:"The resource you are looking for might have been removed, had its name changed, or is temporarily unavailable."`
	Error     ErrorInfo `json:"error"`
	Timestamp string    `json:"timestamp" example:"2025-01-27T10:30:00Z"`
	RequestID string    `json:"request_id,omitempty" example:"req_123"`
}

type ConflictError struct {
	Success   bool      `json:"success" example:"false"`
	Message   string    `json:"message" example:"The resource you are looking for might have been removed, had its name changed, or is temporarily unavailable."`
	Error     ErrorInfo `json:"error"`
	Timestamp string    `json:"timestamp" example:"2025-01-27T10:30:00Z"`
	RequestID string    `json:"request_id,omitempty" example:"req_123"`
}

type InternalServerError struct {
	Success   bool      `json:"success" example:"false"`
	Message   string    `json:"message" example:"Something went wrong. Please try again later."`
	Error     ErrorInfo `json:"error"`
	Timestamp string    `json:"timestamp" example:"2025-01-27T10:30:00Z"`
	RequestID string    `json:"request_id,omitempty" example:"req_123"`
}

package dto

// Generic error examples showing structure by HTTP status
// The actual error codes/messages come from apperrors package at runtime

var (
	// 400 - Client input errors
	BadRequestExample = ErrorInfo{
		Code:    "VALIDATION_ERROR_CODE",
		Message: "Validation error message",
		Details: map[string]string{"field": "error details"},
	}

	// 401 - Authentication errors  
	UnauthorizedExample = ErrorInfo{
		Code:    "AUTH_ERROR_CODE",
		Message: "Authentication error message",
		Details: nil,
	}

	// 403 - Permission errors
	ForbiddenExample = ErrorInfo{
		Code:    "PERMISSION_ERROR_CODE", 
		Message: "Permission error message",
		Details: nil,
	}

	// 404 - Resource not found errors
	NotFoundExample = ErrorInfo{
		Code:    "NOT_FOUND_ERROR_CODE",
		Message: "Resource not found message",
		Details: nil,
	}

	// 409 - Conflict errors (already exists, etc.)
	ConflictExample = ErrorInfo{
		Code:    "CONFLICT_ERROR_CODE",
		Message: "Resource conflict message", 
		Details: nil,
	}

	// 500 - Server errors
	InternalServerExample = ErrorInfo{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "Something went wrong. Please try again later.",
		Details: nil,
	}
)
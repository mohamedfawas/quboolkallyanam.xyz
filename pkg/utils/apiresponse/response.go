package apiresponse

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type StandardResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Timestamp string      `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

type ErrorInfo struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

func Success(c *gin.Context, message string, data interface{}) {
	response := StandardResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		RequestID: getRequestID(c),
	}
	c.JSON(http.StatusOK, response)
}

func Fail(c *gin.Context, err error) {
	httpCode, errorInfo := mapErrorToResponse(err)

	response := StandardResponse{
		Success:   false,
		Message:   "Request failed",
		Error:     errorInfo,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		RequestID: getRequestID(c),
	}

	c.JSON(httpCode, response)
	c.Abort()
}

func FailWithDetails(c *gin.Context, err error, details map[string]string) {
	httpCode, errorInfo := mapErrorToResponse(err)
	if errorInfo.Details == nil {
		errorInfo.Details = make(map[string]string)
	}
	for k, v := range details {
		errorInfo.Details[k] = v
	}

	response := StandardResponse{
		Success:   false,
		Message:   "Request failed",
		Error:     errorInfo,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		RequestID: getRequestID(c),
	}

	c.JSON(httpCode, response)
	c.Abort()
}

func ValidationError(c *gin.Context, message string, details map[string]string) {
	response := StandardResponse{
		Success: false,
		Message: "Validation failed",
		Error: &ErrorInfo{
			Code:    "VALIDATION_ERROR",
			Message: message,
			Details: details,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		RequestID: getRequestID(c),
	}

	c.JSON(http.StatusBadRequest, response)
	c.Abort()
}

func getRequestID(c *gin.Context) string {
	if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
		return requestID
	}
	if requestID := c.GetString("request_id"); requestID != "" {
		return requestID
	}
	return ""
}

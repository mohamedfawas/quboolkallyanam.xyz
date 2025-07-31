package apiresponse

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
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

func Success(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, StandardResponse{
		Success:   true,
		Message:   msg,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		RequestID: c.GetString(constants.ContextKeyRequestID),
	})
}

func Error(c *gin.Context, err error, details map[string]string) {
	// Check if the error is a gRPC error
	if st, ok := status.FromError(err); ok {
		if errorInfo := extractErrorInfo(st); errorInfo != nil {
			httpStatusStr := errorInfo.Metadata[constants.HTTPStatusCode]
			userMessage := errorInfo.Metadata[constants.UserFriendlyMessage]

			httpStatus := http.StatusInternalServerError
			if httpStatusStr != "" {
				if parsed, err := strconv.Atoi(httpStatusStr); err == nil {
					httpStatus = parsed
				}
			}

			if userMessage == "" {
				userMessage = st.Message()
			}

			respond(c, httpStatus, &ErrorInfo{
				Code:    errorInfo.Reason,
				Message: userMessage,
				Details: details,
			})
			return
		}
	}

	// check if the error is an AppError
	var ae *apperrors.AppError
	if errors.As(err, &ae) {
		respond(c, ae.HTTPStatusCode, &ErrorInfo{
			Code:    ae.Code,
			Message: ae.PublicMsg,
			Details: details,
		})
		return
	}

	// Fallback for unexpected or unknown errors.
	respond(c, http.StatusInternalServerError, &ErrorInfo{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: constants.InteralServerErrorMessage,
		Details: nil,
	})
}

func extractErrorInfo(st *status.Status) *errdetails.ErrorInfo {
	for _, detail := range st.Details() {
		if errorInfo, ok := detail.(*errdetails.ErrorInfo); ok {
			return errorInfo
		}
	}
	return nil
}

func respond(c *gin.Context, httpCode int, ei *ErrorInfo) {
	c.JSON(httpCode, StandardResponse{
		Success:   false,
		Message:   "Request failed",
		Error:     ei,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		RequestID: c.GetString(constants.ContextKeyRequestID),
	})
	c.Abort()
}

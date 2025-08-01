package apperrors

import (
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
)

var (
	ErrConversationNotFound = &AppError{
		Err:            errors.New("conversation not found"),
		Code:           "CONVERSATION_NOT_FOUND",
		HTTPStatusCode: http.StatusNotFound,
		GRPCStatusCode: codes.NotFound,
		PublicMsg:      "Conversation not found. Please create the conversation again and try again.",
	}
	ErrInvalidConversationID = &AppError{
		Err:            errors.New("invalid conversation ID"),
		Code:           "INVALID_CONVERSATION_ID",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Invalid conversation ID. Please provide a valid conversation ID.",
	}
	ErrUserNotParticipant = &AppError{
		Err:            errors.New("user not a participant in the conversation"),
		Code:           "USER_NOT_PARTICIPANT",
		HTTPStatusCode: http.StatusForbidden,
		GRPCStatusCode: codes.PermissionDenied,
		PublicMsg:      "You are not a participant in this conversation. Please create the conversation again and try again.",
	}
)
package apperrors

import (
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
)

// media file errors
var (
	ErrFileTooLarge = &AppError{		
		Err:            errors.New("file too large"),
		Code:           "FILE_TOO_LARGE",
		HTTPStatusCode: http.StatusRequestEntityTooLarge,
		GRPCStatusCode: codes.ResourceExhausted,
		PublicMsg:      "The file exceeds the maximum size allowed."}
	ErrInvalidFileType = &AppError{
		Err:            errors.New("invalid file type"),
		Code:           "INVALID_FILE_TYPE",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "The file type is not supported."}
	ErrInvalidImageType = &AppError{
		Err: errors.New("invalid image type"),
		Code: "INVALID_IMAGE_TYPE",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg: "The image type is not supported. The supported formats are jpeg, jpg, and png",
	}
)
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
	ErrInvalidImageFileSize = &AppError{
		Err: errors.New("invalid image file size"),
		Code: "INVALID_IMAGE_FILE_SIZE",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg: "The image file size is not supported. The maximum size is 5MB",
	}
	ErrProfilePhotoNotFound = &AppError{
		Err: errors.New("profile photo not found"),
		Code: "PROFILE_PHOTO_NOT_FOUND",
		HTTPStatusCode: http.StatusNotFound,
		GRPCStatusCode: codes.NotFound,
		PublicMsg: "Profile photo not found for the user. Please upload the profile photo.",
	}
	ErrImageDisplayOrderOccupied = &AppError{
		Err: errors.New("image display order occupied"),
		Code: "IMAGE_DISPLAY_ORDER_OCCUPIED",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg: "The requested image display order is already occupied. Please try a different display order or retry after deleting existing image.",
	}
	ErrDisplayOrderNotOccupied = &AppError{
		Err: errors.New("display order not occupied"),
		Code: "DISPLAY_ORDER_NOT_OCCUPIED",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg: "There are no images at the requested display order.",
	}
	ErrInvalidDisplayOrder = &AppError{
		Err: errors.New("invalid display order"),
		Code: "INVALID_DISPLAY_ORDER",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg: "The requested display order is not valid. The display order must be between 1 and 3.",
	}
)
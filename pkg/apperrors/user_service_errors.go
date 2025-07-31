package apperrors

import (
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
)

// User Profile errors
var (
	ErrUserProfileNotFound = &AppError{
		Err:            errors.New("user profile not found"),
		Code:           "USER_PROFILE_NOT_FOUND",
		HTTPStatusCode: http.StatusNotFound,
		GRPCStatusCode: codes.NotFound,
		PublicMsg:      "User profile not found."}
	ErrInvalidFullName = &AppError{
		Err:            errors.New("invalid full name"),
		Code:           "INVALID_FULL_NAME",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Full name must be 2 to 100 English letters. Please check your input."}
	ErrInvalidCommunity = &AppError{
		Err:            errors.New("invalid community"),
		Code:           "INVALID_COMMUNITY",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "The selected community is not recognized."}
	ErrInvalidMaritalStatus = &AppError{
		Err:            errors.New("invalid marital status"),
		Code:           "INVALID_MARITAL_STATUS",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "The selected marital status is not recognized."}
	ErrInvalidProfession = &AppError{
		Err:            errors.New("invalid profession"),
		Code:           "INVALID_PROFESSION",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "The selected profession is not recognized."}
	ErrInvalidProfessionType = &AppError{
		Err:            errors.New("invalid profession type"),
		Code:           "INVALID_PROFESSION_TYPE",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "The selected profession type is not recognized."}
	ErrInvalidEducationLevel = &AppError{
		Err:            errors.New("invalid highest education level"),
		Code:           "INVALID_HIGHEST_EDUCATION_LEVEL",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "The selected highest education level is not recognized."}
	ErrInvalidHomeDistrict = &AppError{
		Err:            errors.New("invalid home district"),
		Code:           "INVALID_HOME_DISTRICT",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Given home district is not recognized"}
	ErrInvalidDateOfBirth = &AppError{
		Err:            errors.New("invalid date of birth"),
		Code:           "INVALID_DATE_OF_BIRTH",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Give a valid date of birth in YYYY-MM-DD format"}
	ErrInvalidHeight = &AppError{
		Err:            errors.New("invalid height"),
		Code:           "INVALID_HEIGHT",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Height should be between 130 and 220 cm"}
)

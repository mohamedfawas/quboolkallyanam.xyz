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
		PublicMsg:      "User profile not found. Please try to create your profile first."}
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
	ErrInvalidAge = &AppError{
		Err:            errors.New("invalid age"),
		Code:           "INVALID_AGE",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Age should be between 18 and 100 years"}
	ErrInvalidHeightRange = &AppError{
		Err:            errors.New("invalid height range"),
		Code:           "INVALID_HEIGHT_RANGE",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Height should be between 130 and 220 cm and max height should be greater than min height"}
	ErrInvalidAgeRange = &AppError{
		Err:            errors.New("invalid age range"),
		Code:           "INVALID_AGE_RANGE",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Age should be between 18 and 100 years and max age should be greater than min age"}
	ErrUserProfileNotCompleted = &AppError{
		Err:            errors.New("user profile not completed"),
		Code:           "USER_PROFILE_NOT_COMPLETED",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "User profile is not completed. Please complete your profile first."}
	ErrSameGenderProfileAccessDenied = &AppError{
		Err:            errors.New("same gender profile access denied"),
		Code:           "SAME_GENDER_PROFILE_ACCESS_DENIED",
		HTTPStatusCode: http.StatusForbidden,
		GRPCStatusCode: codes.PermissionDenied,
		PublicMsg:      "You cannot access the profile."}
)

// Partner Preferences errors
var (
	ErrPartnerPreferencesNotFound = &AppError{
		Err:            errors.New("partner preferences not found"),
		Code:           "PARTNER_PREFERENCES_NOT_FOUND",
		HTTPStatusCode: http.StatusNotFound,
		GRPCStatusCode: codes.NotFound,
		PublicMsg:      "Partner preferences not found. Please try to create your partner preferences first."}
	ErrPartnerPreferencesAlreadyExists = &AppError{
		Err:            errors.New("partner preferences already exists"),
		Code:           "PARTNER_PREFERENCES_ALREADY_EXISTS",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.AlreadyExists,
		PublicMsg:      "Partner preferences already exists. Please try to update your partner preferences instead."}
)


// Match Making errors
var (
	ErrInvalidMatchAction = &AppError{
		Err:            errors.New("invalid match action"),
		Code:           "INVALID_MATCH_ACTION",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "The selected match action is not recognized. Please try again."}
	ErrInvalidTargetProfileID = &AppError{
		Err:            errors.New("invalid target profile ID"),
		Code:           "INVALID_TARGET_PROFILE_ID",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "The selected target profile ID does not exist. Please try again."}
)
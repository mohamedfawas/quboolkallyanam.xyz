package apperrors

import (
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
)

// subscription plan errors
var (
	ErrSubscriptionPlanNotFound = &AppError{
		Err:            errors.New("subscription plan not found"),
		Code:           "SUBSCRIPTION_PLAN_NOT_FOUND",
		HTTPStatusCode: http.StatusNotFound,
		GRPCStatusCode: codes.NotFound,
		PublicMsg:      "The selected subscription plan does not exist. Please try again."}
	ErrSubscriptionPlanNotActive = &AppError{
		Err:            errors.New("subscription plan not active"),
		Code:           "SUBSCRIPTION_PLAN_NOT_ACTIVE",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "The selected subscription plan is not active. Please select another plan."}
	ErrInvalidSubscriptionPlanID = &AppError{
		Err:            errors.New("invalid subscription plan ID"),
		Code:           "INVALID_SUBSCRIPTION_PLAN_ID",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "The provided subscription plan ID is invalid. Please try again."}
	ErrInvalidSubscriptionPlanDurationDays = &AppError{
		Err:            errors.New("invalid subscription plan duration days"),
		Code:           "INVALID_SUBSCRIPTION_PLAN_DURATION_DAYS",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Invalid subscription plan duration. Please choose a valid duration."}
	ErrInvalidSubscriptionPlanAmount = &AppError{
		Err:            errors.New("invalid subscription plan amount"),
		Code:           "INVALID_SUBSCRIPTION_PLAN_AMOUNT",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Invalid subscription plan amount. Please enter a valid amount."}
	ErrInvalidCurrency = &AppError{
		Err:            errors.New("invalid subscription plan currency"),
		Code:           "INVALID_SUBSCRIPTION_PLAN_CURRENCY",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "The selected currency is not supported. Please choose a valid currency."}
	ErrActiveSubscriptionNotFound = &AppError{
		Err:            errors.New("active subscription not found"),
		Code:           "ACTIVE_SUBSCRIPTION_NOT_FOUND",
		HTTPStatusCode: http.StatusNotFound,
		GRPCStatusCode: codes.NotFound,
		PublicMsg:      "No active subscription found for the given user. Please try again."}
	)

// payment errors
var (
	ErrPaymentNotFound = &AppError{
		Err:            errors.New("payment not found"),
		Code:           "PAYMENT_NOT_FOUND",
		HTTPStatusCode: http.StatusNotFound,
		GRPCStatusCode: codes.NotFound,
		PublicMsg:      "The given payment details does not match with our records. Please try again."}
	ErrPaymentAlreadyCompleted = &AppError{
		Err:            errors.New("payment already completed"),
		Code:           "PAYMENT_ALREADY_COMPLETED",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "The payment has already been completed. "}
	ErrPaymentExpired = &AppError{
		Err:            errors.New("payment expired"),
		Code:           "PAYMENT_EXPIRED",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "The time alloted for payment has expired. Please try again."}
	ErrPaymentSignatureInvalid = &AppError{	
		Err:            errors.New("payment signature invalid"),
		Code:           "PAYMENT_SIGNATURE_INVALID",
		HTTPStatusCode: http.StatusBadRequest,
		GRPCStatusCode: codes.InvalidArgument,
		PublicMsg:      "Payment verification failed due to invalid signature. Please try again."}
	ErrPaymentProcessingFailed = &AppError{
		Err:            errors.New("payment processing failed"),
		Code:           "PAYMENT_PROCESSING_FAILED",
		HTTPStatusCode: http.StatusInternalServerError,
		GRPCStatusCode: codes.Internal,
		PublicMsg:      "Unable to process the payment at the moment. Please try again later."}
)
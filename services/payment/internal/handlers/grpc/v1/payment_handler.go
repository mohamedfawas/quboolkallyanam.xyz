package v1

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	paymentpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/payment/v1"
	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/usecase"
)

type PaymentHandler struct {
	paymentpbv1.UnimplementedPaymentServiceServer
	paymentUsecase      usecase.PaymentUsecase
	subscriptionUsecase usecase.SubscriptionUsecase
}

func NewPaymentHandler(
	paymentUsecase usecase.PaymentUsecase,
	subscriptionUsecase usecase.SubscriptionUsecase,
) *PaymentHandler {
	return &PaymentHandler{
		paymentUsecase:      paymentUsecase,
		subscriptionUsecase: subscriptionUsecase,
	}
}

func (h *PaymentHandler) CreatePaymentOrder(ctx context.Context, req *paymentpbv1.CreatePaymentOrderRequest) (*paymentpbv1.CreatePaymentOrderResponse, error) {
	userID, err := contextutils.GetUserID(ctx)
	if err != nil {
		log.Printf("Failed to get user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "user ID not found: %v", err)
	}

	paymentOrderResponse, err := h.paymentUsecase.CreatePaymentOrder(ctx, userID, req.PlanId)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrSubscriptionPlanNotFound):
			return nil, status.Errorf(codes.NotFound, "Subscription plan not found")
		case errors.Is(err, appErrors.ErrSubscriptionPlanNotActive):
			return nil, status.Errorf(codes.FailedPrecondition, "Subscription plan is not active")
		default:
			log.Printf("Failed to create payment order: %v", err)
			return nil, status.Errorf(codes.Internal, "Something went wrong, please try again")
		}
	}

	expiresAtPb := timestamppb.New(paymentOrderResponse.ExpiresAt)

	return &paymentpbv1.CreatePaymentOrderResponse{
		RazorpayOrderId: paymentOrderResponse.RazorpayOrderID,
		Amount:          paymentOrderResponse.Amount,
		Currency:        paymentOrderResponse.Currency,
		PlanId:          paymentOrderResponse.PlanID,
		ExpiresAt:       expiresAtPb,
	}, nil
}

func (h *PaymentHandler) ShowPaymentPage(ctx context.Context, req *paymentpbv1.ShowPaymentPageRequest) (*paymentpbv1.ShowPaymentPageResponse, error) {
	response, err := h.paymentUsecase.ShowPaymentPage(ctx, req.RazorpayOrderId)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrPaymentNotFound):
			return nil, status.Errorf(codes.NotFound, "Payment not found")
		case errors.Is(err, appErrors.ErrPaymentAlreadyCompleted):
			return nil, status.Errorf(codes.FailedPrecondition, "Payment already completed")
		case errors.Is(err, appErrors.ErrPaymentExpired):
			return nil, status.Errorf(codes.FailedPrecondition, "Payment has expired")
		case errors.Is(err, appErrors.ErrSubscriptionPlanNotFound):
			return nil, status.Errorf(codes.NotFound, "Subscription plan not found")
		case errors.Is(err, appErrors.ErrSubscriptionPlanNotActive):
			return nil, status.Errorf(codes.FailedPrecondition, "Subscription plan is not active")
		default:
			log.Printf("Failed to get payment page data: %v", err)
			return nil, status.Errorf(codes.Internal, "Something went wrong, please try again")
		}
	}

	return &paymentpbv1.ShowPaymentPageResponse{
		RazorpayOrderId:    response.RazorpayOrderID,
		RazorpayKeyId:      response.RazorpayKeyID,
		PlanId:             response.PlanID,
		Amount:             response.Amount,
		DisplayAmount:      response.DisplayAmount,
		PlanDurationInDays: response.PlanDurationInDays,
	}, nil
}

func (h *PaymentHandler) VerifyPayment(ctx context.Context, req *paymentpbv1.VerifyPaymentRequest) (*paymentpbv1.VerifyPaymentResponse, error) {
	log.Printf("[PAYMENT] Verifying payment Handler called for the razorpay order %s", req.RazorpayOrderId)

	var verifyPaymentRequest entity.VerifyPaymentRequest
	verifyPaymentRequest.RazorpayOrderID = req.RazorpayOrderId
	verifyPaymentRequest.RazorpayPaymentID = req.RazorpayPaymentId
	verifyPaymentRequest.RazorpaySignature = req.RazorpaySignature

	verifyPaymentResponse, err := h.paymentUsecase.VerifyPayment(ctx, &verifyPaymentRequest)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrPaymentNotFound):
			return nil, status.Errorf(codes.NotFound, "Payment not found")
		case errors.Is(err, appErrors.ErrPaymentAlreadyCompleted):
			return nil, status.Errorf(codes.FailedPrecondition, "Payment already completed")
		case errors.Is(err, appErrors.ErrPaymentExpired):
			return nil, status.Errorf(codes.FailedPrecondition, "Payment has expired")
		case errors.Is(err, appErrors.ErrPaymentSignatureInvalid):
			return nil, status.Errorf(codes.InvalidArgument, "Payment verification failed")
		case errors.Is(err, appErrors.ErrSubscriptionPlanNotFound):
			return nil, status.Errorf(codes.NotFound, "Subscription plan not found")
		case errors.Is(err, appErrors.ErrSubscriptionPlanNotActive):
			return nil, status.Errorf(codes.FailedPrecondition, "Subscription plan is not active")
		default:
			log.Printf("Failed to verify payment for the razorpay order %s: %v", verifyPaymentRequest.RazorpayOrderID, err)
			return nil, status.Errorf(codes.Internal, "Something went wrong, please try again")
		}
	}

	subscriptionStartDatePb := timestamppb.New(verifyPaymentResponse.SubscriptionStartDate)
	subscriptionEndDatePb := timestamppb.New(verifyPaymentResponse.SubscriptionEndDate)

	return &paymentpbv1.VerifyPaymentResponse{
		SubscriptionId:        verifyPaymentResponse.SubscriptionID,
		SubscriptionStartDate: subscriptionStartDatePb,
		SubscriptionEndDate:   subscriptionEndDatePb,
		SubscriptionStatus:    verifyPaymentResponse.SubscriptionStatus,
	}, nil
}

func (h *PaymentHandler) CreateOrUpdateSubscriptionPlan(ctx context.Context, req *paymentpbv1.CreateOrUpdateSubscriptionPlanRequest) (*paymentpbv1.CreateOrUpdateSubscriptionPlanResponse, error) {
	updateRequest := entity.UpdateSubscriptionPlanRequest{
		ID: req.Id,
	}

	// Convert optional fields from protobuf wrappers to pointers
	if req.DurationDays != nil {
		durationDays := int(req.DurationDays.Value)
		updateRequest.DurationDays = &durationDays
	}

	if req.Amount != nil {
		amount := req.Amount.Value
		updateRequest.Amount = &amount
	}

	if req.Currency != nil {
		currency := req.Currency.Value
		updateRequest.Currency = &currency
	}

	if req.Description != nil {
		description := req.Description.Value
		updateRequest.Description = &description
	}

	if req.IsActive != nil {
		isActive := req.IsActive.Value
		updateRequest.IsActive = &isActive
	}

	err := h.subscriptionUsecase.CreateOrUpdateSubscriptionPlans(ctx, updateRequest)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrInvalidSubscriptionPlanID):
			return nil, status.Errorf(codes.InvalidArgument, "Invalid subscription plan ID")
		case errors.Is(err, appErrors.ErrInvalidSubscriptionPlanDurationDays):
			return nil, status.Errorf(codes.InvalidArgument, "Invalid subscription plan duration days")
		case errors.Is(err, appErrors.ErrInvalidSubscriptionPlanAmount):
			return nil, status.Errorf(codes.InvalidArgument, "Invalid subscription plan amount")
		case errors.Is(err, appErrors.ErrInvalidCurrency):
			return nil, status.Errorf(codes.InvalidArgument, "Invalid currency")
		case errors.Is(err, appErrors.ErrMissingRequiredFields):
			return nil, status.Errorf(codes.InvalidArgument, "Missing required fields for plan creation")
		default:
			log.Printf("Failed to create or update subscription plan: %v", err)
			return nil, status.Errorf(codes.Internal, "Something went wrong, please try again")
		}
	}

	return &paymentpbv1.CreateOrUpdateSubscriptionPlanResponse{
		Success: true,
	}, nil
}

func (h *PaymentHandler) GetSubscriptionPlan(ctx context.Context, req *paymentpbv1.GetSubscriptionPlanRequest) (*paymentpbv1.GetSubscriptionPlanResponse, error) {
	plan, err := h.subscriptionUsecase.GetSubscriptionPlan(ctx, req.PlanId)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrSubscriptionPlanNotFound):
			return nil, status.Errorf(codes.NotFound, "Subscription plan not found")
		default:
			log.Printf("Failed to get subscription plan: %v", err)
			return nil, status.Errorf(codes.Internal, "Something went wrong, please try again")
		}
	}

	createdAtPb := timestamppb.New(plan.CreatedAt)
	updatedAtPb := timestamppb.New(plan.UpdatedAt)

	subscriptionPlan := &paymentpbv1.SubscriptionPlan{
		Id:           plan.ID,
		DurationDays: int32(plan.DurationDays),
		Amount:       plan.Amount,
		Currency:     plan.Currency,
		Description:  plan.Description,
		IsActive:     plan.IsActive,
		CreatedAt:    createdAtPb,
		UpdatedAt:    updatedAtPb,
	}

	return &paymentpbv1.GetSubscriptionPlanResponse{
		Plan: subscriptionPlan,
	}, nil
}

func (h *PaymentHandler) GetActiveSubscriptionPlans(ctx context.Context, req *paymentpbv1.GetActiveSubscriptionPlansRequest) (*paymentpbv1.GetActiveSubscriptionPlansResponse, error) {
	plans, err := h.subscriptionUsecase.GetActiveSubscriptionPlans(ctx)
	if err != nil {
		log.Printf("Failed to get active subscription plans: %v", err)
		return nil, status.Errorf(codes.Internal, "Something went wrong, please try again")
	}

	var subscriptionPlans []*paymentpbv1.SubscriptionPlan
	for _, plan := range plans {
		createdAtPb := timestamppb.New(plan.CreatedAt)
		updatedAtPb := timestamppb.New(plan.UpdatedAt)

		subscriptionPlan := &paymentpbv1.SubscriptionPlan{
			Id:           plan.ID,
			DurationDays: int32(plan.DurationDays),
			Amount:       plan.Amount,
			Currency:     plan.Currency,
			Description:  plan.Description,
			IsActive:     plan.IsActive,
			CreatedAt:    createdAtPb,
			UpdatedAt:    updatedAtPb,
		}
		subscriptionPlans = append(subscriptionPlans, subscriptionPlan)
	}

	return &paymentpbv1.GetActiveSubscriptionPlansResponse{
		Plans: subscriptionPlans,
	}, nil
}

func (h *PaymentHandler) GetActiveSubscriptionByUserID(ctx context.Context, req *paymentpbv1.GetActiveSubscriptionByUserIDRequest) (*paymentpbv1.GetActiveSubscriptionByUserIDResponse, error) {
	userID, err := contextutils.GetUserID(ctx)
	if err != nil {
		log.Printf("Failed to get user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "user ID not found: %v", err)
	}

	subscription, err := h.subscriptionUsecase.GetActiveSubscriptionByUserID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrActiveSubscriptionNotFound):
			return nil, status.Errorf(codes.NotFound, "No active subscription found")
		default:
			log.Printf("Failed to get active subscription for user %s: %v", userID, err)
			return nil, status.Errorf(codes.Internal, "Something went wrong, please try again")
		}
	}

	startDatePb := timestamppb.New(subscription.StartDate)
	endDatePb := timestamppb.New(subscription.EndDate)
	createdAtPb := timestamppb.New(subscription.CreatedAt)
	updatedAtPb := timestamppb.New(subscription.UpdatedAt)

	return &paymentpbv1.GetActiveSubscriptionByUserIDResponse{
		SubscriptionId: subscription.ID,
		PlanId:         subscription.PlanID,
		StartDate:      startDatePb,
		EndDate:        endDatePb,
		Status:         subscription.Status,
		CreatedAt:      createdAtPb,
		UpdatedAt:      updatedAtPb,
	}, nil
}

func (h *PaymentHandler) GetPaymentHistory(ctx context.Context, req *paymentpbv1.GetPaymentHistoryRequest) (*paymentpbv1.GetPaymentHistoryResponse, error) {
	userID, err := contextutils.GetUserID(ctx)
	if err != nil {
		log.Printf("Failed to get user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "user ID not found: %v", err)
	}

	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Failed to parse user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "user ID not found: %v", err)
	}

	paymentHistory, err := h.paymentUsecase.GetPaymentHistory(ctx, userIDUUID)
	if err != nil {
		log.Printf("Failed to get payment history for user %s: %v", userID, err)
		return nil, status.Errorf(codes.Internal, "Something went wrong, please try again")
	}

	var paymentHistoryItems []*paymentpbv1.PaymentHistoryItem
	for _, payment := range paymentHistory {
		createdAtPb := timestamppb.New(payment.CreatedAt)

		paymentHistoryItem := &paymentpbv1.PaymentHistoryItem{
			Id:              payment.ID,
			PlanId:          payment.PlanID,
			RazorpayOrderId: payment.RazorpayOrderID,
			Amount:          payment.Amount,
			Currency:        payment.Currency,
			Status:          payment.Status,
			PaymentMethod:   payment.PaymentMethod,
			CreatedAt:       createdAtPb,
		}
		paymentHistoryItems = append(paymentHistoryItems, paymentHistoryItem)
	}

	return &paymentpbv1.GetPaymentHistoryResponse{
		PaymentHistory: paymentHistoryItems,
	}, nil
}

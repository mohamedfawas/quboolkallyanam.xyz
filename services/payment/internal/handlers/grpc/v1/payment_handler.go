package v1

import (
	"context"

	"github.com/google/uuid"
	paymentpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/payment/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/usecase"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PaymentHandler struct {
	paymentpbv1.UnimplementedPaymentServiceServer
	paymentUsecase      usecase.PaymentUsecase
	subscriptionUsecase usecase.SubscriptionUsecase
	logger              *zap.Logger
}

func NewPaymentHandler(
	paymentUsecase usecase.PaymentUsecase,
	subscriptionUsecase usecase.SubscriptionUsecase,
	logger *zap.Logger,
) *PaymentHandler {
	return &PaymentHandler{
		paymentUsecase:      paymentUsecase,
		subscriptionUsecase: subscriptionUsecase,
		logger:              logger,
	}
}

func (h *PaymentHandler) CreatePaymentOrder(ctx context.Context, req *paymentpbv1.CreatePaymentOrderRequest) (*paymentpbv1.CreatePaymentOrderResponse, error) {
	contextData, err := contextutils.ExtractGrpcContextData(ctx)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, contextData.RequestID),
		zap.String(constants.ContextKeyUserID, contextData.UserID),
	)

	paymentOrderResponse, err := h.paymentUsecase.CreatePaymentOrder(ctx, contextData.UserID, req.PlanId)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to create payment order", zap.Error(err))
		}
		return nil, err
	}

	log.Info("Payment order created successfully",
		zap.String("razorpay_order_id", paymentOrderResponse.RazorpayOrderID),
		zap.String("plan_id", paymentOrderResponse.PlanID),
	)

	return &paymentpbv1.CreatePaymentOrderResponse{
		RazorpayOrderId: paymentOrderResponse.RazorpayOrderID,
		Amount:          paymentOrderResponse.Amount,
		Currency:        paymentOrderResponse.Currency,
		PlanId:          paymentOrderResponse.PlanID,
		ExpiresAt:       timestamppb.New(paymentOrderResponse.ExpiresAt),
	}, nil
}


func (h *PaymentHandler) ShowPaymentPage(ctx context.Context, req *paymentpbv1.ShowPaymentPageRequest) (*paymentpbv1.ShowPaymentPageResponse, error) {
	reqCtx, err := contextutils.ExtractRequestIDFromGrpcContext(ctx)
	if err != nil {
		h.logger.Error("Failed to extract request ID", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, reqCtx.RequestID),
	)

	response, err := h.paymentUsecase.ShowPaymentPage(ctx, req.RazorpayOrderId)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to get payment page data",
				zap.String("razorpay_order_id", req.RazorpayOrderId),
				zap.Error(err))
		}
		return nil, err
	}

	log.Info("Payment page data retrieved",
		zap.String("razorpay_order_id", response.RazorpayOrderID),
	)

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
	reqCtx, err := contextutils.ExtractRequestIDFromGrpcContext(ctx)
	if err != nil {
		h.logger.Error("Failed to extract request ID", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, reqCtx.RequestID),
		zap.String("razorpay_order_id", req.RazorpayOrderId),
	)

	var verifyPaymentRequest entity.VerifyPaymentRequest
	verifyPaymentRequest.RazorpayOrderID = req.RazorpayOrderId
	verifyPaymentRequest.RazorpayPaymentID = req.RazorpayPaymentId
	verifyPaymentRequest.RazorpaySignature = req.RazorpaySignature

	verifyPaymentResponse, err := h.paymentUsecase.VerifyPayment(ctx, &verifyPaymentRequest)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to verify payment", zap.Error(err))
		}
		return nil, err
	}

	log.Info("Payment verified successfully",
		zap.String("subscription_id", verifyPaymentResponse.SubscriptionID),
	)

	return &paymentpbv1.VerifyPaymentResponse{
		SubscriptionId:        verifyPaymentResponse.SubscriptionID,
		SubscriptionStartDate: timestamppb.New(verifyPaymentResponse.SubscriptionStartDate),
		SubscriptionEndDate:   timestamppb.New(verifyPaymentResponse.SubscriptionEndDate),
		SubscriptionStatus:    verifyPaymentResponse.SubscriptionStatus,
	}, nil
}


func (h *PaymentHandler) CreateOrUpdateSubscriptionPlan(ctx context.Context, req *paymentpbv1.CreateOrUpdateSubscriptionPlanRequest) (*paymentpbv1.CreateOrUpdateSubscriptionPlanResponse, error) {
	contextData, err := contextutils.ExtractGrpcContextData(ctx)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, contextData.RequestID),
		zap.String(constants.ContextKeyUserID, contextData.UserID),
	)

	updateRequest := entity.UpdateSubscriptionPlanRequest{
		ID: req.Id,
	}

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

	if err := h.subscriptionUsecase.CreateOrUpdateSubscriptionPlans(ctx, updateRequest); err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to create or update subscription plan", zap.Error(err))
		}
		return nil, err
	}

	log.Info("Subscription plan created/updated successfully",
		zap.String("plan_id", req.Id),
	)

	return &paymentpbv1.CreateOrUpdateSubscriptionPlanResponse{Success: true}, nil
}


func (h *PaymentHandler) GetSubscriptionPlan(ctx context.Context, req *paymentpbv1.GetSubscriptionPlanRequest) (*paymentpbv1.GetSubscriptionPlanResponse, error) {
	reqCtx, err := contextutils.ExtractRequestIDFromGrpcContext(ctx)
	if err != nil {
		h.logger.Error("Failed to extract request ID", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, reqCtx.RequestID),
	)

	plan, err := h.subscriptionUsecase.GetSubscriptionPlan(ctx, req.PlanId)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to get subscription plan", zap.Error(err))
		}
		return nil, err
	}

	log.Info("Subscription plan fetched successfully",
		zap.String("plan_id", plan.ID),
	)

	return &paymentpbv1.GetSubscriptionPlanResponse{
		Plan: &paymentpbv1.SubscriptionPlan{
			Id:           plan.ID,
			DurationDays: int32(plan.DurationDays),
			Amount:       plan.Amount,
			Currency:     plan.Currency,
			Description:  plan.Description,
			IsActive:     plan.IsActive,
			CreatedAt:    timestamppb.New(plan.CreatedAt),
			UpdatedAt:    timestamppb.New(plan.UpdatedAt),
		},
	}, nil
}

func (h *PaymentHandler) GetActiveSubscriptionPlans(ctx context.Context, req *paymentpbv1.GetActiveSubscriptionPlansRequest) (*paymentpbv1.GetActiveSubscriptionPlansResponse, error) {
	reqCtx, err := contextutils.ExtractRequestIDFromGrpcContext(ctx)
	if err != nil {
		h.logger.Error("Failed to extract request ID", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, reqCtx.RequestID),
	)

	plans, err := h.subscriptionUsecase.GetActiveSubscriptionPlans(ctx)
	if err != nil {
		log.Error("Failed to get active subscription plans", zap.Error(err))
		return nil, err
	}

	var subscriptionPlans []*paymentpbv1.SubscriptionPlan
	for _, plan := range plans {
		subscriptionPlans = append(subscriptionPlans, &paymentpbv1.SubscriptionPlan{
			Id:           plan.ID,
			DurationDays: int32(plan.DurationDays),
			Amount:       plan.Amount,
			Currency:     plan.Currency,
			Description:  plan.Description,
			IsActive:     plan.IsActive,
			CreatedAt:    timestamppb.New(plan.CreatedAt),
			UpdatedAt:    timestamppb.New(plan.UpdatedAt),
		})
	}

	log.Info("Active subscription plans fetched successfully",
		zap.Int("count", len(subscriptionPlans)),
	)

	return &paymentpbv1.GetActiveSubscriptionPlansResponse{
		Plans: subscriptionPlans,
	}, nil
}


func (h *PaymentHandler) GetActiveSubscriptionByUserID(ctx context.Context, req *paymentpbv1.GetActiveSubscriptionByUserIDRequest) (*paymentpbv1.GetActiveSubscriptionByUserIDResponse, error) {
	contextData, err := contextutils.ExtractGrpcContextData(ctx)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, contextData.RequestID),
		zap.String(constants.ContextKeyUserID, contextData.UserID),
	)

	subscription, err := h.subscriptionUsecase.GetActiveSubscriptionByUserID(ctx, contextData.UserID)
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to get active subscription", zap.Error(err))
		}
		return nil, err
	}

	log.Info("Active subscription fetched successfully")

	return &paymentpbv1.GetActiveSubscriptionByUserIDResponse{
		SubscriptionId: subscription.ID,
		PlanId:         subscription.PlanID,
		StartDate:      timestamppb.New(subscription.StartDate),
		EndDate:        timestamppb.New(subscription.EndDate),
		Status:         subscription.Status,
		CreatedAt:      timestamppb.New(subscription.CreatedAt),
		UpdatedAt:      timestamppb.New(subscription.UpdatedAt),
	}, nil
}


func (h *PaymentHandler) GetPaymentHistory(ctx context.Context, req *paymentpbv1.GetPaymentHistoryRequest) (*paymentpbv1.GetPaymentHistoryResponse, error) {
	contextData, err := contextutils.ExtractGrpcContextData(ctx)
	if err != nil {
		h.logger.Error("Failed to extract context data", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, contextData.RequestID),
		zap.String(constants.ContextKeyUserID, contextData.UserID),
	)

	userIDUUID, err := uuid.Parse(contextData.UserID)
	if err != nil {
		log.Error("Failed to parse user ID", zap.Error(err))
		return nil, err
	}

	paymentHistory, err := h.paymentUsecase.GetPaymentHistory(ctx, userIDUUID)
	if err != nil {
		log.Error("Failed to get payment history", zap.Error(err))
		return nil, err
	}

	var paymentHistoryItems []*paymentpbv1.PaymentHistoryItem
	for _, payment := range paymentHistory {
		paymentHistoryItems = append(paymentHistoryItems, &paymentpbv1.PaymentHistoryItem{
			Id:              payment.ID,
			PlanId:          payment.PlanID,
			RazorpayOrderId: payment.RazorpayOrderID,
			Amount:          payment.Amount,
			Currency:        payment.Currency,
			Status:          payment.Status,
			PaymentMethod:   payment.PaymentMethod,
			CreatedAt:       timestamppb.New(payment.CreatedAt),
		})
	}

	log.Info("Payment history fetched successfully",
		zap.Int("count", len(paymentHistoryItems)),
	)

	return &paymentpbv1.GetPaymentHistoryResponse{
		PaymentHistory: paymentHistoryItems,
	}, nil
}


func (h *PaymentHandler) GetCompletedPaymentDetails(ctx context.Context, req *paymentpbv1.GetCompletedPaymentDetailsRequest) (*paymentpbv1.GetCompletedPaymentDetailsResponse, error) {
	reqCtx, err := contextutils.ExtractRequestIDFromGrpcContext(ctx)
	if err != nil {
		h.logger.Error("Failed to extract request ID", zap.Error(err))
		return nil, err
	}
	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, reqCtx.RequestID),
	)

	payments, err := h.paymentUsecase.GetCompletedPaymentDetails(ctx, int(req.Page), int(req.Limit))
	if err != nil {
		if !apperrors.IsAppError(err) {
			log.Error("Failed to get completed payment details", zap.Error(err))
		}
		return nil, err
	}

	resp := &paymentpbv1.GetCompletedPaymentDetailsResponse{}
	for _, p := range payments {
		resp.Payments = append(resp.Payments, &paymentpbv1.CompletedPaymentDetail{
			Id:              p.ID,
			UserId:          p.UserID,
			PlanId:          p.PlanID,
			RazorpayOrderId: p.RazorpayOrderID,
			Amount:          p.Amount,
			Currency:        p.Currency,
			Status:          p.Status,
			PaymentMethod:   p.PaymentMethod,
			CreatedAt:       timestamppb.New(p.CreatedAt),
			UpdatedAt:       timestamppb.New(p.UpdatedAt),
		})
	}

	log.Info("Completed payments fetched successfully", zap.Int("count", len(resp.Payments)))
	return resp, nil
}
package v1

import (
	paymentpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/payment/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

// /////////////////////////// Create Payment Order //////////////////////////////
func MapCreatePaymentOrderRequest(req dto.PaymentOrderRequest) *paymentpbv1.CreatePaymentOrderRequest {
	return &paymentpbv1.CreatePaymentOrderRequest{
		PlanId: req.PlanID,
	}
}

func MapCreatePaymentOrderResponse(resp *paymentpbv1.CreatePaymentOrderResponse) *dto.PaymentOrderResponse {
	return &dto.PaymentOrderResponse{
		RazorpayOrderID: resp.RazorpayOrderId,
		Amount:          resp.Amount,
		Currency:        resp.Currency,
		PlanID:          resp.PlanId,
		ExpiresAt:       resp.ExpiresAt.AsTime(),
	}
}

// /////////////////////////// Show Payment Page //////////////////////////////
func MapShowPaymentPageRequest(req dto.ShowPaymentPageRequest) *paymentpbv1.ShowPaymentPageRequest {
	return &paymentpbv1.ShowPaymentPageRequest{
		RazorpayOrderId: req.RazorpayOrderID,
	}
}

func MapShowPaymentPageResponse(resp *paymentpbv1.ShowPaymentPageResponse) *dto.ShowPaymentPageResponse {
	return &dto.ShowPaymentPageResponse{
		RazorpayOrderID:    resp.RazorpayOrderId,
		RazorpayKeyID:      resp.RazorpayKeyId,
		PlanID:             resp.PlanId,
		Amount:             resp.Amount,
		DisplayAmount:      resp.DisplayAmount,
		PlanDurationInDays: int(resp.PlanDurationInDays),
	}
}

// /////////////////////////// Verify Payment //////////////////////////////
func MapVerifyPaymentRequest(req dto.VerifyPaymentRequest) *paymentpbv1.VerifyPaymentRequest {
	return &paymentpbv1.VerifyPaymentRequest{
		RazorpayOrderId:   req.RazorpayOrderID,
		RazorpayPaymentId: req.RazorpayPaymentID,
		RazorpaySignature: req.RazorpaySignature,
	}
}

func MapVerifyPaymentResponse(resp *paymentpbv1.VerifyPaymentResponse) *dto.VerifyPaymentResponse {
	return &dto.VerifyPaymentResponse{
		SubscriptionID:        resp.SubscriptionId,
		SubscriptionStartDate: resp.SubscriptionStartDate.AsTime(),
		SubscriptionEndDate:   resp.SubscriptionEndDate.AsTime(),
		SubscriptionStatus:    resp.SubscriptionStatus,
	}
}

// /////////////////////////// Create Or Update Subscription Plan //////////////////////////////
func MapCreateOrUpdateSubscriptionPlanRequest(req dto.UpdateSubscriptionPlanRequest) *paymentpbv1.CreateOrUpdateSubscriptionPlanRequest {
	grpcReq := &paymentpbv1.CreateOrUpdateSubscriptionPlanRequest{
		Id: req.ID,
	}

	// Convert optional fields to protobuf wrappers
	if req.DurationDays != nil {
		grpcReq.DurationDays = &wrapperspb.Int32Value{Value: int32(*req.DurationDays)}
	}

	if req.Amount != nil {
		grpcReq.Amount = &wrapperspb.DoubleValue{Value: *req.Amount}
	}

	if req.Currency != nil {
		grpcReq.Currency = &wrapperspb.StringValue{Value: *req.Currency}
	}

	if req.Description != nil {
		grpcReq.Description = &wrapperspb.StringValue{Value: *req.Description}
	}

	if req.IsActive != nil {
		grpcReq.IsActive = &wrapperspb.BoolValue{Value: *req.IsActive}
	}

	return grpcReq
}

func MapCreateOrUpdateSubscriptionPlanResponse(resp *paymentpbv1.CreateOrUpdateSubscriptionPlanResponse) *dto.CreateOrUpdateSubscriptionPlanResponse {
	return &dto.CreateOrUpdateSubscriptionPlanResponse{
		Success: resp.Success,
	}
}

// /////////////////////////// Get Subscription Plan //////////////////////////////
func MapGetSubscriptionPlanResponse(resp *paymentpbv1.GetSubscriptionPlanResponse) *dto.SubscriptionPlan {
	if resp.Plan == nil {
		return nil
	}

	return &dto.SubscriptionPlan{
		ID:           resp.Plan.Id,
		DurationDays: int(resp.Plan.DurationDays),
		Amount:       resp.Plan.Amount,
		Currency:     resp.Plan.Currency,
		Description:  resp.Plan.Description,
		IsActive:     resp.Plan.IsActive,
		CreatedAt:    resp.Plan.CreatedAt.AsTime(),
		UpdatedAt:    resp.Plan.UpdatedAt.AsTime(),
	}
}

// /////////////////////////// Get Active Subscription Plans //////////////////////////////
func MapGetActiveSubscriptionPlansResponse(resp *paymentpbv1.GetActiveSubscriptionPlansResponse) []*dto.SubscriptionPlan {
	var plans []*dto.SubscriptionPlan

	for _, plan := range resp.Plans {
		subscriptionPlan := &dto.SubscriptionPlan{
			ID:           plan.Id,
			DurationDays: int(plan.DurationDays),
			Amount:       plan.Amount,
			Currency:     plan.Currency,
			Description:  plan.Description,
			IsActive:     plan.IsActive,
			CreatedAt:    plan.CreatedAt.AsTime(),
			UpdatedAt:    plan.UpdatedAt.AsTime(),
		}
		plans = append(plans, subscriptionPlan)
	}

	return plans
}

// /////////////////////////// Get Active Subscription By User ID //////////////////////////////
func MapGetActiveSubscriptionByUserIDResponse(resp *paymentpbv1.GetActiveSubscriptionByUserIDResponse) *dto.ActiveSubscription {
	return &dto.ActiveSubscription{
		SubscriptionID: resp.SubscriptionId,
		PlanID:         resp.PlanId,
		StartDate:      resp.StartDate.AsTime(),
		EndDate:        resp.EndDate.AsTime(),
		Status:         resp.Status,
		CreatedAt:      resp.CreatedAt.AsTime(),
		UpdatedAt:      resp.UpdatedAt.AsTime(),
	}
}

// /////////////////////////// Get Payment History //////////////////////////////
func MapGetPaymentHistoryResponse(resp *paymentpbv1.GetPaymentHistoryResponse) []*dto.GetPaymentHistoryResponse {
	var paymentHistory []*dto.GetPaymentHistoryResponse

	for _, payment := range resp.PaymentHistory {
		paymentHistory = append(paymentHistory, &dto.GetPaymentHistoryResponse{
			ID:              payment.Id,
			PlanID:          payment.PlanId,
			RazorpayOrderID: payment.RazorpayOrderId,
			Amount:          payment.Amount,
			Currency:        payment.Currency,
			Status:          payment.Status,
			PaymentMethod:   payment.PaymentMethod,
			CreatedAt:       payment.CreatedAt.AsTime(),
		})
	}

	return paymentHistory
}

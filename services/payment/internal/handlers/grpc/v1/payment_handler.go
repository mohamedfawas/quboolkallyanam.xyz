package v1

import (
	"context"
	"errors"
	"log"

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
	paymentUsecase usecase.PaymentUsecase
}

func NewPaymentHandler(paymentUsecase usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{paymentUsecase: paymentUsecase}
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

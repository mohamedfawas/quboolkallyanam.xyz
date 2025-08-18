package payments

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	paymentEvents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/payment"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/timeutil"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
	"gorm.io/gorm"
)

func (u *paymentUsecase) VerifyPayment(ctx context.Context, req *entity.VerifyPaymentRequest) (*entity.VerifyPaymentResponse, error) {

	if err := u.razorpayService.VerifySignature(req.RazorpayOrderID,
		req.RazorpayPaymentID,
		req.RazorpaySignature); err != nil {
		// Check if it's a signature mismatch vs other errors
		if strings.Contains(err.Error(), "signature mismatch") {
			return nil, apperrors.ErrPaymentSignatureInvalid
		}
		return nil, apperrors.ErrPaymentProcessingFailed
	}

	verifyPaymentResponse := &entity.VerifyPaymentResponse{
		SubscriptionID:        "",
		SubscriptionStartDate: time.Time{},
		SubscriptionEndDate:   time.Time{},
		SubscriptionStatus:    "",
	}

	var paymentVerifiedEvent paymentEvents.PaymentVerified

	err := u.txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
		payment, err := u.paymentRepository.GetPaymentDetailsByRazorpayOrderIDTx(ctx, tx, req.RazorpayOrderID)
		if err != nil {
			return err
		}
		if payment == nil {
			return apperrors.ErrPaymentNotFound
		}

		if payment.Status == constants.PaymentStatusCompleted {
			return apperrors.ErrPaymentAlreadyCompleted
		}

		if time.Now().UTC().After(payment.ExpiresAt) {
			return apperrors.ErrPaymentExpired
		}

		plan, err := u.subscriptionPlanRepository.GetPlanByIDTx(ctx, tx, payment.PlanID)
		if err != nil {
			return err
		}
		if plan == nil {
			return apperrors.ErrSubscriptionPlanNotFound
		}

		payment.RazorpayPaymentID = req.RazorpayPaymentID
		payment.RazorpaySignature = req.RazorpaySignature
		payment.Status = constants.PaymentStatusCompleted
		payment.UpdatedAt = time.Now().UTC()

		if err := u.paymentRepository.UpdatePaymentTx(ctx, tx, payment); err != nil {
			return err
		}

		existingSubscription, err := u.subscriptionRepository.GetActiveSubscriptionByUserIDTx(ctx, tx, payment.UserID)
		if err != nil {
			return err
		}

		// Deactivate existing subscription if found
		if existingSubscription != nil {
			existingSubscription.Status = constants.SubscriptionStatusCancelled
			existingSubscription.UpdatedAt = time.Now().UTC()
			if err := u.subscriptionRepository.UpdateSubscriptionTx(ctx, tx, existingSubscription); err != nil {
				return err
			}
		}

		// Create new subscription
		now := time.Now().UTC()
		subscription := &entity.Subscription{
			UserID:    payment.UserID,
			PlanID:    payment.PlanID,
			StartDate: now,
			EndDate:   now.AddDate(0, 0, plan.DurationDays),
			Status:    constants.SubscriptionStatusActive,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := u.subscriptionRepository.CreateSubscriptionTx(ctx, tx, subscription); err != nil {
			return err
		}

		// Get the newly created subscription to populate response
		newSubscription, err := u.subscriptionRepository.GetActiveSubscriptionByUserIDTx(ctx, tx, payment.UserID)
		if err != nil {
			return err
		}

		verifyPaymentResponse.SubscriptionID = strconv.FormatInt(newSubscription.ID, 10)
		verifyPaymentResponse.SubscriptionStartDate = timeutil.ToIST(newSubscription.StartDate)
		verifyPaymentResponse.SubscriptionEndDate = timeutil.ToIST(newSubscription.EndDate)
		verifyPaymentResponse.SubscriptionStatus = newSubscription.Status

		paymentVerifiedEvent = paymentEvents.PaymentVerified{
			UserID:              payment.UserID,
			SubscriptionID:      strconv.FormatInt(newSubscription.ID, 10),
			SubscriptionEndDate: newSubscription.EndDate,
			PlanID:              payment.PlanID,
			PaymentID:           payment.RazorpayPaymentID,
			Timestamp:           time.Now().UTC(),
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if u.eventPublisher != nil {
		if err := u.eventPublisher.PublishPaymentVerified(ctx, paymentVerifiedEvent); err != nil {
			// Note: We don't return error here as the payment was successful
		}
	}

	return verifyPaymentResponse, nil
}
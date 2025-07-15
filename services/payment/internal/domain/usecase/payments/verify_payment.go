package payments

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/timeutil"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/entity"
	"gorm.io/gorm"
)

func (u *paymentUsecase) VerifyPayment(ctx context.Context, req *entity.VerifyPaymentRequest) (*entity.VerifyPaymentResponse, error) {

	if err := u.razorpayService.VerifySignature(req.RazorpayOrderID,
		req.RazorpayPaymentID,
		req.RazorpaySignature); err != nil {
		log.Printf("failed to verify payment signature: %v", err)
		return nil, err
	}

	verifyPaymentResponse := &entity.VerifyPaymentResponse{
		SubscriptionID:        "",
		SubscriptionStartDate: time.Time{},
		SubscriptionEndDate:   time.Time{},
		SubscriptionStatus:    "",
	}

	err := u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		payment, err := u.paymentRepository.GetPaymentDetailsByRazorpayOrderID(txCtx, req.RazorpayOrderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return appErrors.ErrPaymentNotFound
			}
			log.Printf("failed to get payment details by Razorpay Order ID: %v", err)
			return err
		}

		if payment.Status == constants.PaymentStatusCompleted {
			return appErrors.ErrPaymentAlreadyCompleted
		}

		if time.Now().UTC().After(payment.ExpiresAt) {
			return appErrors.ErrPaymentExpired
		}

		plan, err := u.subscriptionPlanRepository.GetPlanByID(txCtx, payment.PlanID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return appErrors.ErrSubscriptionPlanNotFound
			}
			log.Printf("failed to get subscription plan: %v", err)
			return err
		}

		payment.RazorpayPaymentID = req.RazorpayPaymentID
		payment.RazorpaySignature = req.RazorpaySignature
		payment.Status = constants.PaymentStatusCompleted
		payment.UpdatedAt = time.Now().UTC()

		if err := u.paymentRepository.UpdatePayment(txCtx, payment); err != nil {
			log.Printf("failed to update payment: %v", err)
			return err
		}

		existingSubscription, err := u.subscriptionRepository.GetActiveSubscriptionByUserID(txCtx, payment.UserID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("failed to check existing subscription: %v", err)
			return err
		}

		// Deactivate existing subscription if found
		if existingSubscription != nil {
			existingSubscription.Status = constants.SubscriptionStatusCancelled
			existingSubscription.UpdatedAt = time.Now().UTC()
			if err := u.subscriptionRepository.UpdateSubscription(txCtx, existingSubscription); err != nil {
				log.Printf("failed to deactivate existing subscription: %v", err)
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

		if err := u.subscriptionRepository.CreateSubscription(txCtx, subscription); err != nil {
			log.Printf("failed to create subscription: %v", err)
			return err
		}

		// Get the newly created subscription to populate response
		newSubscription, err := u.subscriptionRepository.GetActiveSubscriptionByUserID(txCtx, payment.UserID)
		if err != nil {
			log.Printf("failed to get active subscription: %v", err)
			return err
		}

		verifyPaymentResponse.SubscriptionID = strconv.FormatInt(newSubscription.ID, 10)
		verifyPaymentResponse.SubscriptionStartDate = timeutil.ToIST(newSubscription.StartDate)
		verifyPaymentResponse.SubscriptionEndDate = timeutil.ToIST(newSubscription.EndDate)
		verifyPaymentResponse.SubscriptionStatus = newSubscription.Status

		log.Printf("Payment verified and subscription created successfully for user %s", payment.UserID)
		return nil
	})

	if err != nil {
		log.Printf("failed to verify payment for razorpay order %s: %v", req.RazorpayOrderID, err)
		return nil, err
	}

	return verifyPaymentResponse, nil
}

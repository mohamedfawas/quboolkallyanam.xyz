package payments

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/payment/razorpay"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/repository"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/usecase"
)

type paymentUsecase struct {
	paymentRepository          repository.PaymentsRepository
	subscriptionPlanRepository repository.SubscriptionPlansRepository
	razorpayService            *razorpay.Service
}

func NewPaymentUsecase(
	paymentRepository repository.PaymentsRepository,
	subscriptionPlanRepository repository.SubscriptionPlansRepository,
	razorpayService *razorpay.Service,
) usecase.PaymentUsecase {
	return &paymentUsecase{paymentRepository: paymentRepository,
		subscriptionPlanRepository: subscriptionPlanRepository,
		razorpayService:            razorpayService,
	}
}

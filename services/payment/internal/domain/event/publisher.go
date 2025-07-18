package event

import (
	"context"

	paymentEvents "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/events/payment"
)

type EventPublisher interface {
	PublishPaymentVerified(ctx context.Context, event paymentEvents.PaymentVerified) error
	// PublishPaymentFailed(ctx context.Context, event PaymentFailed) error : just examples to understand how i thought about desinging events
}

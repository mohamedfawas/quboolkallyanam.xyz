package razorpay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/google/uuid"
	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
	razorpay "github.com/razorpay/razorpay-go"
)

type Service struct {
	client    *razorpay.Client
	keyID     string
	secretKey string
}

func NewService(keyID, secretKey string) *Service {
	client := razorpay.NewClient(keyID, secretKey)
	return &Service{
		client:    client,
		keyID:     keyID,
		secretKey: secretKey,
	}
}

// KeyID returns the public Key ID for front-end initialization
func (s *Service) KeyID() string {
	return s.keyID
}

// CreateOrder creates a new Razorpay order and returns its ID
func (s *Service) CreateOrder(amount float64, currency string) (string, error) {
	amountInPaise := int64(amount * 100)
	receipt := uuid.New().String()
	data := map[string]interface{}{
		"amount":          amountInPaise, // amount in paise
		"currency":        currency,      // e.g., "INR"
		"receipt":         receipt,
		"payment_capture": 1,
	}

	resp, err := s.client.Order.Create(data, nil)
	if err != nil {
		log.Printf("razorpay: create order failed: %v", err)
		return "", appErrors.ErrRazorpayOrderCreation
	}

	id, ok := resp["id"].(string)
	if !ok {
		log.Printf("razorpay: invalid order id in response")
		return "", appErrors.ErrRazorpayOrderCreation
	}

	return id, nil
}

// VerifySignature checks HMAC-SHA256 signature for payment integrity
func (s *Service) VerifySignature(orderID, paymentID, signature string) error {
	msg := orderID + "|" + paymentID
	mac := hmac.New(sha256.New, []byte(s.secretKey))
	mac.Write([]byte(msg))
	expected := hex.EncodeToString(mac.Sum(nil))
	if expected != signature {
		log.Printf("razorpay: signature mismatch: expected %s, got %s", expected, signature)
		return appErrors.ErrPaymentSignatureInvalid
	}
	return nil
}

// VerifyWebhookSignature validates Razorpay webhook payload signature
func (s *Service) VerifyWebhookSignature(signature string, payload []byte) error {
	mac := hmac.New(sha256.New, []byte(s.secretKey))
	mac.Write(payload)
	expected := hex.EncodeToString(mac.Sum(nil))
	if expected != signature {
		log.Printf("razorpay: webhook signature mismatch: expected %s, got %s", expected, signature)
		return fmt.Errorf("razorpay: webhook signature mismatch")
	}
	return nil
}

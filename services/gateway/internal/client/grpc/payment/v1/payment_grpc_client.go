package v1

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	paymentpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/payment/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/client"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

type paymentGRPCClient struct {
	conn   *grpc.ClientConn
	client paymentpbv1.PaymentServiceClient
}

func NewPaymentGRPCClient(
	ctx context.Context,
	address string,
	useTLS bool,
	tlsConfig *tls.Config) (client.PaymentClient, error) {
	_, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var opts []grpc.DialOption
	if useTLS {
		creds := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		creds := insecure.NewCredentials()
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	cc, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client to %s: %w", address, err)
	}

	return &paymentGRPCClient{
		conn:   cc,
		client: paymentpbv1.NewPaymentServiceClient(cc),
	}, nil
}

func (c *paymentGRPCClient) CreatePaymentOrder(ctx context.Context, req dto.PaymentOrderRequest) (*dto.PaymentOrderResponse, error) {
	userID, ok := ctx.Value(constants.ContextKeyUserID).(string)
	if !ok {
		return nil, fmt.Errorf("user ID not found in context")
	}
	ctx = contextutils.SetUserContext(ctx, userID)

	grpcReq := MapCreatePaymentOrderRequest(req)
	grpcResp, err := c.client.CreatePaymentOrder(ctx, grpcReq)
	if err != nil {
		log.Printf("CreatePaymentOrder error in payment grpc client: %v", err)
		return nil, err
	}
	return MapCreatePaymentOrderResponse(grpcResp), nil
}

func (c *paymentGRPCClient) ShowPaymentPage(ctx context.Context, req dto.ShowPaymentPageRequest) (*dto.ShowPaymentPageResponse, error) {
	grpcReq := MapShowPaymentPageRequest(req)
	grpcResp, err := c.client.ShowPaymentPage(ctx, grpcReq)
	if err != nil {
		log.Printf("ShowPaymentPage error in payment grpc client: %v", err)
		return nil, err
	}
	return MapShowPaymentPageResponse(grpcResp), nil
}

func (c *paymentGRPCClient) VerifyPayment(ctx context.Context, req dto.VerifyPaymentRequest) (*dto.VerifyPaymentResponse, error) {
	grpcReq := MapVerifyPaymentRequest(req)
	grpcResp, err := c.client.VerifyPayment(ctx, grpcReq)
	if err != nil {
		log.Printf("VerifyPayment error in payment grpc client: %v", err)
		return nil, err
	}
	return MapVerifyPaymentResponse(grpcResp), nil
}

func (c *paymentGRPCClient) Close() error {
	return c.conn.Close()
}

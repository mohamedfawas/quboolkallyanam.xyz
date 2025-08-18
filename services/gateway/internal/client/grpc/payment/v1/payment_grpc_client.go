package v1

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	paymentpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/payment/v1"
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
	var err error
	ctx, err = contextutils.PrepareGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	grpcReq := MapCreatePaymentOrderRequest(req)
	grpcResp, err := c.client.CreatePaymentOrder(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapCreatePaymentOrderResponse(grpcResp), nil
}


func (c *paymentGRPCClient) ShowPaymentPage(ctx context.Context, req dto.ShowPaymentPageRequest) (*dto.ShowPaymentPageResponse, error) {
	var err error
	ctx, err = contextutils.PrepareRequestIDForGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	grpcReq := MapShowPaymentPageRequest(req)
	grpcResp, err := c.client.ShowPaymentPage(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapShowPaymentPageResponse(grpcResp), nil
}


func (c *paymentGRPCClient) VerifyPayment(ctx context.Context, req dto.VerifyPaymentRequest) (*dto.VerifyPaymentResponse, error) {
	var err error
	ctx, err = contextutils.PrepareRequestIDForGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	grpcReq := MapVerifyPaymentRequest(req)
	grpcResp, err := c.client.VerifyPayment(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapVerifyPaymentResponse(grpcResp), nil
}


func (c *paymentGRPCClient) CreateOrUpdateSubscriptionPlan(ctx context.Context, req dto.UpdateSubscriptionPlanRequest) (*dto.CreateOrUpdateSubscriptionPlanResponse, error) {
	var err error
	ctx, err = contextutils.PrepareGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	grpcReq := MapCreateOrUpdateSubscriptionPlanRequest(req)
	grpcResp, err := c.client.CreateOrUpdateSubscriptionPlan(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapCreateOrUpdateSubscriptionPlanResponse(grpcResp), nil
}


func (c *paymentGRPCClient) GetSubscriptionPlan(ctx context.Context, planID string) (*dto.SubscriptionPlan, error) {
	var err error
	ctx, err = contextutils.PrepareRequestIDForGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	grpcReq := &paymentpbv1.GetSubscriptionPlanRequest{
		PlanId: planID,
	}
	grpcResp, err := c.client.GetSubscriptionPlan(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapGetSubscriptionPlanResponse(grpcResp), nil
}

func (c *paymentGRPCClient) GetActiveSubscriptionPlans(ctx context.Context) ([]*dto.SubscriptionPlan, error) {
	var err error
	ctx, err = contextutils.PrepareRequestIDForGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	grpcReq := &paymentpbv1.GetActiveSubscriptionPlansRequest{}
	grpcResp, err := c.client.GetActiveSubscriptionPlans(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapGetActiveSubscriptionPlansResponse(grpcResp), nil
}


func (c *paymentGRPCClient) GetActiveSubscriptionByUserID(ctx context.Context) (*dto.ActiveSubscription, error) {
	var err error
	ctx, err = contextutils.PrepareGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	grpcReq := &paymentpbv1.GetActiveSubscriptionByUserIDRequest{}
	grpcResp, err := c.client.GetActiveSubscriptionByUserID(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapGetActiveSubscriptionByUserIDResponse(grpcResp), nil
}


func (c *paymentGRPCClient) GetPaymentHistory(ctx context.Context) ([]*dto.GetPaymentHistoryResponse, error) {
	var err error
	ctx, err = contextutils.PrepareGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	grpcReq := &paymentpbv1.GetPaymentHistoryRequest{}
	grpcResp, err := c.client.GetPaymentHistory(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapGetPaymentHistoryResponse(grpcResp), nil
}


func (c *paymentGRPCClient) GetCompletedPaymentDetails(ctx context.Context, req dto.GetCompletedPaymentDetailsRequest) (*dto.GetCompletedPaymentDetailsResponse, error) {
	var err error
	ctx, err = contextutils.PrepareGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	grpcReq := &paymentpbv1.GetCompletedPaymentDetailsRequest{
		Page:  req.Page,
		Limit: req.Limit,
	}
	grpcResp, err := c.client.GetCompletedPaymentDetails(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	return MapGetCompletedPaymentDetailsResponse(grpcResp), nil
}

func (c *paymentGRPCClient) Close() error {
	return c.conn.Close()
}

package chat

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	chatpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/chat/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/client"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

type chatGRPCClient struct {
	conn   *grpc.ClientConn
	client chatpbv1.ChatServiceClient
}

func NewChatGRPCClient(ctx context.Context,
	address string,
	useTLS bool,
	tlsConfig *tls.Config) (client.ChatClient, error) {
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

	return &chatGRPCClient{
		conn:   cc,
		client: chatpbv1.NewChatServiceClient(cc),
	}, nil
}

func (c *chatGRPCClient) CreateConversation(ctx context.Context, req dto.CreateConversationRequest) (*dto.CreateConversationResponse, error) {
	userID, ok := ctx.Value(constants.ContextKeyUserID).(string)
	if !ok {
		return nil, errors.New("user ID not found in context")
	}
	ctx = contextutils.SetUserContext(ctx, userID)

	grpcReq := MapCreateConversationRequest(req)
	grpcResp, err := c.client.CreateConversation(ctx, grpcReq)
	if err != nil {
		log.Printf("CreateConversation error in chat grpc client: %v", err)
		return nil, err
	}
	return MapCreateConversationResponse(grpcResp), nil
}

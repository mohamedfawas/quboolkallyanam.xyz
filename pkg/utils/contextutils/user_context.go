package contextutils

import (
	"context"
	"fmt"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"google.golang.org/grpc/metadata"
)

func SetUserContext(ctx context.Context, userID string) context.Context {
	return appendToOutgoingContext(ctx, constants.ContextKeyUserID, userID)
}

func SetRequestIDContext(ctx context.Context, requestID string) context.Context {
	return appendToOutgoingContext(ctx, constants.ContextKeyRequestID, requestID)
}

func GetUserID(ctx context.Context) (string, error) {
	return extractFromIncomingContext(ctx, constants.ContextKeyUserID)
}

func GetRequestID(ctx context.Context) (string, error) {
	return extractFromIncomingContext(ctx, constants.ContextKeyRequestID)
}

func appendToOutgoingContext(ctx context.Context, key, value string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		md = md.Copy()
	}
	md.Set(key, value)
	return metadata.NewOutgoingContext(ctx, md)
}

func extractFromIncomingContext(ctx context.Context, key string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("no metadata found in context")
	}

	values := md.Get(key)
	if len(values) == 0 {
		return "", fmt.Errorf("%s not found in metadata", key)
	}

	return values[0], nil
}


type GrpcContextData struct {
	RequestID string
	UserID    string
}


func ExtractGrpcContextData(ctx context.Context) (*GrpcContextData, error) {
	requestID, err := GetRequestID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get request ID from context: %w", err)
	}

	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID from context: %w", err)
	}

	return &GrpcContextData{
		RequestID: requestID,
		UserID:    userID,
	}, nil
}

func PrepareGrpcContext(ctx context.Context) (context.Context, error) {
	userID, ok := ctx.Value(constants.ContextKeyUserID).(string)
	if !ok {
		return nil, fmt.Errorf("user ID context missing")
	}

	requestID, ok := ctx.Value(constants.ContextKeyRequestID).(string)
	if !ok {
		return nil, fmt.Errorf("request ID context missing")
	}

	ctx = SetUserContext(ctx, userID)
	ctx = SetRequestIDContext(ctx, requestID)

	return ctx, nil
}

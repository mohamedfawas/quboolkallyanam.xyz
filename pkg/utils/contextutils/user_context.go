package contextutils

import (
	"context"
	"fmt"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"google.golang.org/grpc/metadata"
)

func SetUserContext(ctx context.Context, userID string) context.Context {
	md := metadata.New(map[string]string{
		constants.ContextKeyUserID: userID,
	})
	return metadata.NewOutgoingContext(ctx, md)
}

func GetUserID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("no metadata found in context")
	}

	userIDs := md.Get(constants.ContextKeyUserID)
	if len(userIDs) == 0 {
		return "", fmt.Errorf("user ID not found in metadata")
	}

	return userIDs[0], nil
}

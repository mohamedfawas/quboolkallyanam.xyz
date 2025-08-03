package v1

import (
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/user/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
)

func (h *UserHandler) RecordMatchAction(ctx context.Context, req *userpbv1.RecordMatchActionRequest) (*userpbv1.RecordMatchActionResponse, error) {
	userID, err := contextutils.GetUserID(ctx)
	if err != nil {
		log.Printf("Failed to get user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "user ID not found: %v", err)
	}

	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Failed to parse user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	if req.Action == "" {
		log.Printf("Action is required")
		return nil, status.Errorf(codes.InvalidArgument, "action is required")
	}

	if req.TargetProfileId == 0 {
		log.Printf("Target profile ID is required")
		return nil, status.Errorf(codes.InvalidArgument, "target profile ID is required")
	}

	success, err := h.matchMakingUsecase.RecordMatchAction(ctx, userIDUUID, uint(req.TargetProfileId), req.Action)
	if err != nil {
		log.Printf("Failed to record match action: %v", err)
		return nil, err
	}

	return &userpbv1.RecordMatchActionResponse{Success: success}, nil
}

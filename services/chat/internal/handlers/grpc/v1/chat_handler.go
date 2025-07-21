package v1

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/usecase"

	chatpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/chat/v1"
	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
)

type ChatHandler struct {
	chatpbv1.UnimplementedChatServiceServer
	chatUsecase usecase.ChatUsecase
}

func NewChatHandler(chatUsecase usecase.ChatUsecase) *ChatHandler {
	return &ChatHandler{chatUsecase: chatUsecase}
}

func (h *ChatHandler) CreateConversation(ctx context.Context, req *chatpbv1.CreateConversationRequest) (*chatpbv1.CreateConversationResponse, error) {
	userID, err := contextutils.GetUserID(ctx)
	if err != nil {
		log.Printf("Failed to get user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "user ID not found: %v", err)
	}

	conversation, err := h.chatUsecase.CreateConversation(ctx, userID, req.PartnerProfileId)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrUserNotFound):
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
		default:
			return nil, status.Errorf(codes.Internal, "failed to create conversation: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to create conversation: %v", err)
	}

	participantIDs := make([]string, len(conversation.ParticipantIDs))
	for i, id := range conversation.ParticipantIDs {
		participantIDs[i] = string(id)
	}

	return &chatpbv1.CreateConversationResponse{
		ConversationId: conversation.ID.Hex(),
		ParticipantIds: participantIDs,
		CreatedAt:      timestamppb.New(conversation.CreatedAt),
		UpdatedAt:      timestamppb.New(conversation.UpdatedAt),
	}, nil
}

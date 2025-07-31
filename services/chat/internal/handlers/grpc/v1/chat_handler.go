package v1

import (
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/usecase"

	chatpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/chat/v1"
	appErrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"go.uber.org/zap"
)

type ChatHandler struct {
	chatpbv1.UnimplementedChatServiceServer
	chatUsecase usecase.ChatUsecase
	logger      *zap.Logger
}

func NewChatHandler(
	chatUsecase usecase.ChatUsecase,
	logger *zap.Logger) *ChatHandler {

	return &ChatHandler{
		chatUsecase: chatUsecase,
		logger:      logger}
}

func (h *ChatHandler) CreateConversation(
	ctx context.Context,
	req *chatpbv1.CreateConversationRequest) (*chatpbv1.CreateConversationResponse, error) {

	requestID, err := contextutils.GetRequestID(ctx)
	if err != nil {
		h.logger.Error("Failed to get request ID From Context",
			zap.Error(err))
		return nil, status.Errorf(codes.Internal, constants.InteralServerErrorMessage)
	}

	userID, err := contextutils.GetUserID(ctx)
	if err != nil {
		h.logger.Error("Failed to get user ID From Context", zap.Error(err))
		return nil, status.Errorf(codes.Internal, constants.InteralServerErrorMessage)
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, requestID),
		zap.String(constants.ContextKeyUserID, userID),
	)

	conversation, err := h.chatUsecase.CreateConversation(ctx, userID, req.PartnerProfileId)
	if err != nil {
		if !appErrors.IsAppError(err) {
			log.Error("Failed to create conversation",
				zap.Error(err))
		}
		return nil, err
	}

	log.Info("Conversation created successfully",
		zap.String("conversation_id", conversation.ConversationID.Hex()))

	return &chatpbv1.CreateConversationResponse{
		ConversationId:   conversation.ConversationID.Hex(),
		ParticipantNames: conversation.Participants,
		CreatedAt:        timestamppb.New(conversation.CreatedAt),
	}, nil
}

func (h *ChatHandler) SendMessage(ctx context.Context, req *chatpbv1.SendMessageRequest) (*chatpbv1.SendMessageResponse, error) {
	userID, err := contextutils.GetUserID(ctx)
	if err != nil {
		log.Printf("Failed to get user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "user ID not found: %v", err)
	}

	message, err := h.chatUsecase.SendMessage(ctx, req.ConversationId, userID, req.Content)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to send message: %v", err)
	}

	return &chatpbv1.SendMessageResponse{
		MessageId:      message.ID.Hex(),
		ConversationId: message.ConversationID.Hex(),
		SenderId:       string(message.SenderID),
		Content:        message.Content,
		SentAt:         timestamppb.New(message.SentAt),
	}, nil
}

func (h *ChatHandler) GetConversation(ctx context.Context, req *chatpbv1.GetConversationRequest) (*chatpbv1.GetConversationResponse, error) {
	conversation, err := h.chatUsecase.GetConversationByID(ctx, req.ConversationId)
	if err != nil {
		log.Printf("Failed to get conversation: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get conversation: %v", err)
	}

	participantIDs := make([]string, len(conversation.ParticipantIDs))
	for i, id := range conversation.ParticipantIDs {
		participantIDs[i] = string(id)
	}

	return &chatpbv1.GetConversationResponse{
		ConversationId: conversation.ID.Hex(),
		ParticipantIds: participantIDs,
		CreatedAt:      timestamppb.New(conversation.CreatedAt),
		UpdatedAt:      timestamppb.New(conversation.UpdatedAt),
	}, nil
}

func (h *ChatHandler) GetUserConversations(ctx context.Context, req *chatpbv1.GetUserConversationsRequest) (*chatpbv1.GetUserConversationsResponse, error) {
	userID, err := contextutils.GetUserID(ctx)
	if err != nil {
		log.Printf("Failed to get user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "user ID not found: %v", err)
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Failed to parse user ID as UUID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format: %v", err)
	}

	conversations, paginationData, err := h.chatUsecase.GetUserConversations(ctx, userUUID, int(req.Limit), int(req.Offset))
	if err != nil {
		log.Printf("Failed to get user conversations: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get user conversations: %v", err)
	}

	// Convert domain entities to proto messages
	conversationInfos := make([]*chatpbv1.ConversationInfo, len(conversations))
	for i, conv := range conversations {
		participantIDs := make([]string, len(conv.ParticipantIDs))
		for j, id := range conv.ParticipantIDs {
			participantIDs[j] = string(id)
		}

		conversationInfo := &chatpbv1.ConversationInfo{
			ConversationId: conv.ID.Hex(),
			ParticipantIds: participantIDs,
			CreatedAt:      timestamppb.New(conv.CreatedAt),
			UpdatedAt:      timestamppb.New(conv.UpdatedAt),
		}

		// Handle optional last_message_at field
		if conv.LastMessageAt != nil {
			conversationInfo.LastMessageAt = timestamppb.New(*conv.LastMessageAt)
		}

		conversationInfos[i] = conversationInfo
	}

	paginationInfo := &chatpbv1.PaginationInfo{
		TotalCount: paginationData.TotalCount,
		Limit:      int32(paginationData.Limit),
		Offset:     int32(paginationData.Offset),
		HasMore:    paginationData.HasMore,
	}

	return &chatpbv1.GetUserConversationsResponse{
		Conversations: conversationInfos,
		Pagination:    paginationInfo,
	}, nil
}

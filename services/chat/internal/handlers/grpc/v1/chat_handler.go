package v1

import (
	"context"

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
		return nil, err
	}

	userID, err := contextutils.GetUserID(ctx)
	if err != nil {
		h.logger.Error("Failed to get user ID From Context", zap.Error(err))
		return nil, err
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

func (h *ChatHandler) SendMessage(
	ctx context.Context,
	req *chatpbv1.SendMessageRequest) (*chatpbv1.SendMessageResponse, error) {

	requestID, err := contextutils.GetRequestID(ctx)
	if err != nil {
		h.logger.Error("Failed to get request ID From Context",
			zap.Error(err))
		return nil, err
	}

	userID, err := contextutils.GetUserID(ctx)
	if err != nil {
		h.logger.Error("Failed to get user ID From Context", zap.Error(err))
		return nil, err
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, requestID),
		zap.String(constants.ContextKeyUserID, userID),
	)

	message, err := h.chatUsecase.SendMessage(ctx, req.ConversationId, userID, req.Content)
	if err != nil {
		if !appErrors.IsAppError(err) {
			log.Error("Failed to send message",
				zap.Error(err))
		}
		return nil, err
	}

	log.Info("Message sent successfully",
		zap.String("message_id", message.MessageID.Hex()))

	return &chatpbv1.SendMessageResponse{
		MessageId:      message.MessageID.Hex(),
		ConversationId: message.ConversationID.Hex(),
		SenderId:       string(message.SenderID),
		SenderName:     message.SenderName,
		Content:        message.Content,
		SentAt:         timestamppb.New(message.SentAt),
	}, nil
}

func (h *ChatHandler) GetConversation(ctx context.Context, req *chatpbv1.GetConversationRequest) (*chatpbv1.GetConversationResponse, error) {
	requestID, err := contextutils.GetRequestID(ctx)
	if err != nil {
		h.logger.Error("Failed to get request ID From Context",
			zap.Error(err))
		return nil, err
	}

	userID, err := contextutils.GetUserID(ctx)
	if err != nil {
		h.logger.Error("Failed to get user ID From Context", zap.Error(err))
		return nil, err
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, requestID),
		zap.String(constants.ContextKeyUserID, userID),
	)

	conversation, err := h.chatUsecase.GetConversationByID(ctx, req.ConversationId)
	if err != nil {
		if !appErrors.IsAppError(err) {
			log.Error("Failed to get conversation",
				zap.Error(err))
		}
		return nil, err
	}

	log.Info("Conversation retrieved successfully",
		zap.String("conversation_id", conversation.ConversationID.Hex()))

	return &chatpbv1.GetConversationResponse{
		ConversationId: conversation.ConversationID.Hex(),
		ParticipantIds: conversation.ParticipantIDs,
		CreatedAt:      timestamppb.New(conversation.CreatedAt),
		UpdatedAt:      timestamppb.New(conversation.UpdatedAt),
	}, nil
}

func (h *ChatHandler) GetMessagesByConversationId(
	ctx context.Context,
	req *chatpbv1.GetMessagesByConversationIdRequest) (*chatpbv1.GetMessagesByConversationIdResponse, error) {

	requestID, err := contextutils.GetRequestID(ctx)
	if err != nil {
		h.logger.Error("Failed to get request ID From Context",
			zap.Error(err))
		return nil, err
	}

	userID, err := contextutils.GetUserID(ctx)
	if err != nil {
		h.logger.Error("Failed to get user ID From Context", zap.Error(err))
		return nil, err
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, requestID),
		zap.String(constants.ContextKeyUserID, userID),
	)

	response, err := h.chatUsecase.GetMessagesByConversationID(ctx, req.ConversationId, userID, req.Limit, req.Offset)
	if err != nil {
		if !appErrors.IsAppError(err) {
			log.Error("Failed to get messages by conversation ID",
				zap.Error(err),
				zap.String("conversation_id", req.ConversationId))
		}
		return nil, err
	}

	protoMessages := make([]*chatpbv1.MessageInfo, 0, len(response.Messages))
	for _, msg := range response.Messages {
		protoMsg := &chatpbv1.MessageInfo{
			MessageId:  msg.MessageID,
			SenderId:   msg.SenderID,
			SenderName: msg.SenderName,
			Content:    msg.Content,
			SentAt:     timestamppb.New(msg.SentAt),
		}
		protoMessages = append(protoMessages, protoMsg)
	}

	protoPagination := &chatpbv1.PaginationInfo{
		TotalCount: response.Pagination.TotalCount,
		Limit:      response.Pagination.Limit,
		Offset:     response.Pagination.Offset,
		HasMore:    response.Pagination.HasMore,
	}

	log.Info("Messages retrieved successfully",
		zap.String("conversation_id", req.ConversationId),
		zap.Int("message_count", len(protoMessages)),
		zap.Int64("total_count", response.Pagination.TotalCount))

	return &chatpbv1.GetMessagesByConversationIdResponse{
		Messages:   protoMessages,
		Pagination: protoPagination,
	}, nil
}

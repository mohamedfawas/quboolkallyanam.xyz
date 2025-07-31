package chat

import (
	chatpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/chat/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/pagination"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

// /////////////////////////// Create Conversation //////////////////////////////
func MapCreateConversationRequest(req dto.CreateConversationRequest) *chatpbv1.CreateConversationRequest {
	return &chatpbv1.CreateConversationRequest{
		PartnerProfileId: req.PartnerUserProfileID,
	}
}

func MapCreateConversationResponse(resp *chatpbv1.CreateConversationResponse) *dto.CreateConversationResponse {
	return &dto.CreateConversationResponse{
		ConversationID:   resp.ConversationId,
		ParticipantNames: resp.ParticipantNames,
		CreatedAt:        resp.CreatedAt.AsTime(),
	}
}

// /////////////////////////// Send Message //////////////////////////////
func MapSendMessageRequest(req dto.SendMessageRequest) *chatpbv1.SendMessageRequest {
	return &chatpbv1.SendMessageRequest{
		ConversationId: req.ConversationID,
		Content:        req.Content,
	}
}

func MapSendMessageResponse(resp *chatpbv1.SendMessageResponse) *dto.SendMessageResponse {
	return &dto.SendMessageResponse{
		MessageID:      resp.MessageId,
		ConversationID: resp.ConversationId,
		SenderID:       resp.SenderId,
		Content:        resp.Content,
		SentAt:         resp.SentAt.AsTime(),
	}
}

// /////////////////////////// Get Conversation //////////////////////////////
func MapGetConversationRequest(req dto.GetConversationRequest) *chatpbv1.GetConversationRequest {
	return &chatpbv1.GetConversationRequest{
		ConversationId: req.ConversationID,
	}
}

func MapGetConversationResponse(resp *chatpbv1.GetConversationResponse) *dto.GetConversationResponse {
	return &dto.GetConversationResponse{
		ConversationID: resp.ConversationId,
		ParticipantIDs: resp.ParticipantIds,
		CreatedAt:      resp.CreatedAt.AsTime(),
		UpdatedAt:      resp.UpdatedAt.AsTime(),
	}
}

// /////////////////////////// Get User Conversations //////////////////////////////
func MapGetUserConversationsRequest(req dto.GetUserConversationsRequest) *chatpbv1.GetUserConversationsRequest {
	return &chatpbv1.GetUserConversationsRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
	}
}

func MapGetUserConversationsResponse(resp *chatpbv1.GetUserConversationsResponse) *dto.GetUserConversationsResponse {
	conversations := make([]dto.ConversationInfo, len(resp.Conversations))
	for i, conv := range resp.Conversations {
		conversationInfo := dto.ConversationInfo{
			ConversationID: conv.ConversationId,
			ParticipantIDs: conv.ParticipantIds,
			CreatedAt:      conv.CreatedAt.AsTime(),
			UpdatedAt:      conv.UpdatedAt.AsTime(),
		}

		// Handle optional last_message_at field
		if conv.LastMessageAt != nil {
			lastMessageAt := conv.LastMessageAt.AsTime()
			conversationInfo.LastMessageAt = &lastMessageAt
		}

		conversations[i] = conversationInfo
	}

	return &dto.GetUserConversationsResponse{
		Conversations: conversations,
		Pagination: pagination.PaginationData{
			TotalCount: resp.Pagination.TotalCount,
			Limit:      int(resp.Pagination.Limit),
			Offset:     int(resp.Pagination.Offset),
			HasMore:    resp.Pagination.HasMore,
		},
	}
}

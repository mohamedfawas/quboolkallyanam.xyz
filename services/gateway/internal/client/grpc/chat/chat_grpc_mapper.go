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
		SenderName:     resp.SenderName,
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

// /////////////////////////// Get Messages By Conversation ID //////////////////////////////
func MapGetMessagesByConversationIdRequest(req dto.GetMessagesByConversationIdRequest) *chatpbv1.GetMessagesByConversationIdRequest {
	return &chatpbv1.GetMessagesByConversationIdRequest{
		ConversationId: req.ConversationID,
		Limit:          req.Limit,
		Offset:         req.Offset,
	}
}

func MapGetMessagesByConversationIdResponse(resp *chatpbv1.GetMessagesByConversationIdResponse) *dto.GetMessagesByConversationIdResponse {
	messages := make([]dto.MessageInfo, 0, len(resp.Messages))
	for _, msg := range resp.Messages {
		messages = append(messages, dto.MessageInfo{
			MessageID:  msg.MessageId,
			SenderID:   msg.SenderId,
			SenderName: msg.SenderName,
			Content:    msg.Content,
			SentAt:     msg.SentAt.AsTime(),
		})
	}

	return &dto.GetMessagesByConversationIdResponse{
		Messages: messages,
		Pagination: pagination.PaginationData{
			TotalCount: resp.Pagination.TotalCount,
			Limit:      int(resp.Pagination.Limit),
			Offset:     int(resp.Pagination.Offset),
			HasMore:    resp.Pagination.HasMore,
		},
	}
}

package chat

import (
	chatpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/chat/v1"
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
		ConversationID: resp.ConversationId,
		ParticipantIDs: resp.ParticipantIds,
		CreatedAt:      resp.CreatedAt.AsTime(),
		UpdatedAt:      resp.UpdatedAt.AsTime(),
	}
}

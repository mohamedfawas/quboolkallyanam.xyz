package dto

import "time"

type CreateConversationRequest struct {
	PartnerUserProfileID int64 `json:"partner_user_profile_id"`
}

type CreateConversationResponse struct {
	ConversationID string    `json:"conversation_id"`
	ParticipantIDs []string  `json:"participant_ids"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

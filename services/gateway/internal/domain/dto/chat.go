package dto

import (
	"time"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/pagination"
)

type CreateConversationRequest struct {
	PartnerUserProfileID int64 `json:"partner_user_profile_id"`
}

type CreateConversationResponse struct {
	ConversationID   string    `json:"conversation_id"`
	ParticipantNames []string  `json:"participant_names"`
	CreatedAt        time.Time `json:"created_at"`
}

type SendMessageRequest struct {
	ConversationID string `json:"conversation_id" binding:"required"`
	Content        string `json:"content" binding:"required"`
}

type SendMessageResponse struct {
	MessageID      string    `json:"message_id"`
	ConversationID string    `json:"conversation_id"`
	SenderID       string    `json:"sender_id"`
	Content        string    `json:"content"`
	SentAt         time.Time `json:"sent_at"`
}

type GetConversationRequest struct {
	ConversationID string `json:"conversation_id"`
}

type GetConversationResponse struct {
	ConversationID string    `json:"conversation_id"`
	ParticipantIDs []string  `json:"participant_ids"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}


type GetUserConversationsRequest struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ConversationInfo struct {
	ConversationID string     `json:"conversation_id"`
	ParticipantIDs []string   `json:"participant_ids"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	LastMessageAt  *time.Time `json:"last_message_at,omitempty"`
}

type GetUserConversationsResponse struct {
	Conversations []ConversationInfo `json:"conversations"`
	Pagination    pagination.PaginationData     `json:"pagination"`
}

type WebSocketMessage struct {
	Type           string    `json:"type"` // "message", "error", etc.
	ConversationID string    `json:"conversation_id,omitempty"`
	MessageID      string    `json:"message_id,omitempty"`
	SenderID       string    `json:"sender_id,omitempty"`
	Content        string    `json:"content,omitempty"`
	SentAt         time.Time `json:"sent_at,omitempty"`
}

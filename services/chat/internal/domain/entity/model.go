package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// If you had used a new type (type A B),
// MongoDB wouldn't serialize it correctly unless you wrote extra logic
type ConversationID = primitive.ObjectID
type MessageID = primitive.ObjectID

type UserID string

type Conversation struct {
	ID             ConversationID `bson:"_id,omitempty" json:"id"`
	ParticipantIDs []UserID       `bson:"participant_ids" json:"participant_ids"`
	CreatedAt      time.Time      `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time      `bson:"updated_at" json:"updated_at"`
	LastMessageAt  *time.Time     `bson:"last_message_at,omitempty" json:"last_message_at,omitempty"`
}

type Message struct {
	ID             MessageID      `bson:"_id,omitempty" json:"id"`
	ConversationID ConversationID `bson:"conversation_id" json:"conversation_id"`
	SenderID       UserID         `bson:"sender_id" json:"sender_id"`
	Content        string         `bson:"content" json:"content"`
	SentAt         time.Time      `bson:"sent_at" json:"sent_at"`
}

func NewConversationID() ConversationID {
	return primitive.NewObjectID()
}

func NewMessageID() MessageID {
	return primitive.NewObjectID()
}

type CreateConversationResponse struct {
	ConversationID ConversationID `json:"conversation_id"`
	Participants   []string       `json:"participants"`
	CreatedAt      time.Time      `json:"created_at"`
}

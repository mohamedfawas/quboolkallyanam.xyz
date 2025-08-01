package validation

import "go.mongodb.org/mongo-driver/bson/primitive"

func IsValidConversationID(conversationID string) bool {
	_, err := primitive.ObjectIDFromHex(conversationID)
	return err == nil
}
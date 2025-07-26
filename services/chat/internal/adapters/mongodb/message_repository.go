package mongodb

import (
	"context"
	"fmt"
	"log"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/mongodb"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/chat/internal/domain/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type messageRepository struct {
	db *mongo.Database
}

func NewMessageRepository(client *mongodb.Client) repository.MessageRepository {
	return &messageRepository{
		db: client.Database,
	}
}

func (r *messageRepository) CreateMessage(ctx context.Context, message *entity.Message) error {
	coll := r.db.Collection(constants.MongoDBCollectionMessages)
	if _, err := coll.InsertOne(ctx, message); err != nil {
		log.Printf("Error inserting message: %v", err)
		return fmt.Errorf("mongodb: insert message: %w", err)
	}
	return nil
}
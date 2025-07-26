package chat

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/security/jwt"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin
		// In production, you should be more restrictive
		return true
	},
}

type ConnectionManager struct {
	connections map[string]*websocket.Conn // userID -> connection
	mutex       sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*websocket.Conn),
	}
}

func (cm *ConnectionManager) AddConnection(userID string, conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Close existing connection if any
	if existingConn, exists := cm.connections[userID]; exists {
		existingConn.Close()
	}

	cm.connections[userID] = conn
	log.Printf("User %s connected via WebSocket", userID)
}

func (cm *ConnectionManager) RemoveConnection(userID string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if conn, exists := cm.connections[userID]; exists {
		conn.Close()
		delete(cm.connections, userID)
		log.Printf("User %s disconnected from WebSocket", userID)
	}
}

func (cm *ConnectionManager) GetConnection(userID string) (*websocket.Conn, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	conn, exists := cm.connections[userID]
	return conn, exists
}

func (cm *ConnectionManager) BroadcastToUsers(userIDs []string, message interface{}) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	for _, userID := range userIDs {
		if conn, exists := cm.connections[userID]; exists {
			if err := conn.WriteJSON(message); err != nil {
				log.Printf("Failed to send message to user %s: %v", userID, err)
				// Connection might be dead, remove it
				go cm.RemoveConnection(userID)
			}
		}
	}
}

// Global connection manager instance
var connManager = NewConnectionManager()

func (h *ChatHandler) HandleWebSocket(c *gin.Context) {
	// Authenticate WebSocket connection using query parameter
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	// You need to inject JWT manager from server - for now using hardcoded config
	// This should be passed from the handler initialization
	jwtManager := jwt.NewJWTManager(jwt.JWTConfig{
		SecretKey:          "your-256-bit-secret-replace-in-production", // Should match your config
		AccessTokenMinutes: 15,
		RefreshTokenDays:   7,
		Issuer:             "qubool-kallyanam",
	})

	userID, _, err := jwtManager.ExtractUserIDAndRole(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// Add connection to manager
	connManager.AddConnection(userID, conn)
	defer connManager.RemoveConnection(userID)

	// Handle incoming messages
	for {
		var wsMessage dto.WebSocketMessage
		if err := conn.ReadJSON(&wsMessage); err != nil {
			log.Printf("Error reading WebSocket message: %v", err)
			break
		}

		// Process the message
		if err := h.handleIncomingMessage(c.Request.Context(), userID, wsMessage); err != nil {
			log.Printf("Error processing message: %v", err)
			// Send error back to client
			errorMsg := dto.WebSocketMessage{
				Type:    "error",
				Content: fmt.Sprintf("Failed to send message: %v", err),
			}
			conn.WriteJSON(errorMsg)
		}
	}
}

func (h *ChatHandler) handleIncomingMessage(ctx context.Context, senderID string, wsMessage dto.WebSocketMessage) error {
	switch wsMessage.Type {
	case "message":
		// Send message through chat service
		sendMessageReq := dto.SendMessageRequest{
			ConversationID: wsMessage.ConversationID,
			Content:        wsMessage.Content,
		}

		// Add sender ID to context
		ctx = context.WithValue(ctx, constants.ContextKeyUserID, senderID)

		response, err := h.chatUsecase.SendMessage(ctx, sendMessageReq)
		if err != nil {
			return err
		}

		// Broadcast message to all participants
		broadcastMessage := dto.WebSocketMessage{
			Type:           "message",
			ConversationID: response.ConversationID,
			MessageID:      response.MessageID,
			SenderID:       response.SenderID,
			Content:        response.Content,
			SentAt:         response.SentAt,
		}

		// Get conversation participants and broadcast
		conversationReq := dto.GetConversationRequest{
			ConversationID: wsMessage.ConversationID,
		}
		conversation, err := h.chatUsecase.GetConversation(ctx, conversationReq)
		if err != nil {
			return err
		}

		connManager.BroadcastToUsers(conversation.ParticipantIDs, broadcastMessage)

		return nil
	default:
		return fmt.Errorf("unknown message type: %s", wsMessage.Type)
	}
}

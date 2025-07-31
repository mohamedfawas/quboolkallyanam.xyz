package chat

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
	"go.uber.org/zap"
)

// upgrader converts a normal HTTP request into a WebSocket connection.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin
		// In production, you should be more restrictive
		return true
	},
}

// ConnectionManager keeps track of active WebSocket connections by user ID.
type ConnectionManager struct {
	connections map[string]*websocket.Conn // map of userID to WebSocket
	mutex       sync.RWMutex               // read-write lock around the map
	logger      *zap.Logger
}

func NewConnectionManager(logger *zap.Logger) *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*websocket.Conn),
		logger:      logger,
	}
}

// AddConnection registers a new WebSocket for a user.
// If they already had one, we close the old connection.
func (cm *ConnectionManager) AddConnection(userID string, conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Close existing connection if any
	if existingConn, exists := cm.connections[userID]; exists {
		existingConn.Close()
	}

	cm.connections[userID] = conn
	cm.logger.Info("WebSocket connection established",
		zap.String(constants.ContextKeyUserID, userID))
}

// RemoveConnection cleans up when a user disconnects or an error happens.
func (cm *ConnectionManager) RemoveConnection(userID string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if conn, exists := cm.connections[userID]; exists {
		conn.Close()                   // close the WebSocket
		delete(cm.connections, userID) // remove from map
		cm.logger.Info("WebSocket connection closed",
			zap.String(constants.ContextKeyUserID, userID))
	}
}

// GetConnection retrieves a user's WebSocket if it exists
func (cm *ConnectionManager) GetConnection(userID string) (*websocket.Conn, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	conn, exists := cm.connections[userID]
	return conn, exists
}

// SendToUser writes a JSON message to the given user's WebSocket
func (cm *ConnectionManager) SendToUser(userID string, message interface{}) error {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	conn, exists := cm.connections[userID]
	if !exists {
		return fmt.Errorf("user %s not connected", userID)
	}

	if err := conn.WriteJSON(message); err != nil {
		cm.logger.Error("Failed to send WebSocket message",
			zap.String(constants.ContextKeyUserID, userID),
			zap.Error(err))
		// Connection might be dead, remove it
		go cm.RemoveConnection(userID)
		return err
	}

	return nil
}

func (h *ChatHandler) authenticateWebSocket(c *gin.Context) (string, error) {
	authHeader := c.GetHeader(constants.HeaderAuthorization)
	if authHeader == "" || !strings.HasPrefix(authHeader, constants.BearerTokenPrefix) {
		return "", fmt.Errorf("missing or invalid Authorization header")
	}

	token := strings.TrimPrefix(authHeader, constants.BearerTokenPrefix)
	userID, role, err := h.jwtManager.ExtractUserIDAndRole(token)
	if err != nil {
		return "", fmt.Errorf("invalid token")
	}

	if role != constants.RolePremiumUser {
		return "", fmt.Errorf("unauthorized")
	}

	return userID, nil
}

// HandleWebSocket is the main entry point when a client hits the WebSocket endpoint
func (h *ChatHandler) HandleWebSocket(c *gin.Context) {
	requestID, exists := c.Get(constants.ContextKeyRequestID)
	if !exists {
		h.logger.Error("request ID context missing")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "request ID context missing"})
		return
	}
	ctx := context.WithValue(c.Request.Context(), constants.ContextKeyRequestID, requestID)

	// Authenticate WebSocket connection
	userID, err := h.authenticateWebSocket(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
		return
	}

	log := h.logger.With(
		zap.String(constants.ContextKeyRequestID, requestID.(string)),
		zap.String(constants.ContextKeyUserID, userID),
	)

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error("Failed to upgrade WebSocket connection", zap.Error(err))
		return
	}
	defer conn.Close()

	// Add connection to manager
	h.connManager.AddConnection(userID, conn)
	defer h.connManager.RemoveConnection(userID)

	log.Info("WebSocket handler started")

	// Handle incoming messages
	h.handleWebSocketMessages(ctx, userID, conn, log)
}

// handleWebSocketMessages listens for incoming JSON messages and processes them
func (h *ChatHandler) handleWebSocketMessages(ctx context.Context, userID string, conn *websocket.Conn, logger *zap.Logger) {
	for {
		var wsMessage dto.WebSocketMessage
		// ReadJSON blocks until a message arrives or an error happens
		if err := conn.ReadJSON(&wsMessage); err != nil {
			logger.Error("Error reading WebSocket message", zap.Error(err))
			break
		}

		// Delegate to handler based on message Type
		if err := h.handleIncomingMessage(ctx, userID, wsMessage, logger); err != nil {
			logger.Error("Error processing WebSocket message",
				zap.Error(err))

			// Send error back to client using standardized format
			h.sendErrorToClient(conn, err, logger)
		}
	}
}

func (h *ChatHandler) sendErrorToClient(conn *websocket.Conn, err error, logger *zap.Logger) {
	errorMsg := dto.WebSocketMessage{
		Type:    "error",
		Content: fmt.Sprintf("Failed to send message: %v", err),
	}

	// Attempt to write the error; log if it fails
	if writeErr := conn.WriteJSON(errorMsg); writeErr != nil {
		logger.Error("Failed to send error message to WebSocket client",
			zap.Error(writeErr))
	}
}

// handleIncomingMessage decides what to do based on the message Type field
func (h *ChatHandler) handleIncomingMessage(ctx context.Context, senderID string, wsMessage dto.WebSocketMessage, logger *zap.Logger) error {
	switch wsMessage.Type {
	case "message":
		return h.handleChatMessage(ctx, senderID, wsMessage, logger)
	default:
		return fmt.Errorf("unknown message type: %s", wsMessage.Type)
	}
}

// handleChatMessage handles a chat 'message' type by saving and forwarding it
func (h *ChatHandler) handleChatMessage(ctx context.Context, senderID string, wsMessage dto.WebSocketMessage, logger *zap.Logger) error {
	sendMessageReq := dto.SendMessageRequest{
		ConversationID: wsMessage.ConversationID,
		Content:        wsMessage.Content,
	}

	// Add sender ID to context using existing utility
	ctx = context.WithValue(ctx, constants.ContextKeyUserID, senderID)

	response, err := h.chatUsecase.SendMessage(ctx, sendMessageReq)
	if err != nil {
		return err
	}

	// Send message to the other participant (simplified for 1-on-1 chat)
	return h.sendToOtherParticipant(ctx, response, senderID, logger)
}

func (h *ChatHandler) sendToOtherParticipant(ctx context.Context, response *dto.SendMessageResponse, senderID string, logger *zap.Logger) error {
	broadcastMessage := dto.WebSocketMessage{
		Type:           "message",
		ConversationID: response.ConversationID,
		MessageID:      response.MessageID,
		SenderID:       response.SenderID,
		Content:        response.Content,
		SentAt:         response.SentAt,
	}

	// Get conversation to find the other participant
	conversationReq := dto.GetConversationRequest{
		ConversationID: response.ConversationID,
	}

	conversation, err := h.chatUsecase.GetConversation(ctx, conversationReq)
	if err != nil {
		return err
	}

	// Find the other participant (since it's always 1-on-1)
	otherParticipantID := h.findOtherParticipant(conversation.ParticipantIDs, senderID)
	if otherParticipantID == "" {
		return fmt.Errorf("other participant not found in conversation")
	}

	// Send to the other participant only
	if err := h.connManager.SendToUser(otherParticipantID, broadcastMessage); err != nil {
		logger.Debug("Other participant not connected or failed to send",
			zap.Error(err))
		// This is not an error - the other user might just be offline
	}

	return nil
}

func (h *ChatHandler) findOtherParticipant(participantIDs []string, currentUserID string) string {
	for _, participantID := range participantIDs {
		if participantID != currentUserID {
			return participantID
		}
	}
	return ""
}

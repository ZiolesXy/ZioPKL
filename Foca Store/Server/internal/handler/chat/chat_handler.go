package chat

import (
	"log"
	"net/http"
	"voca-store/internal/domain/dto/request"
	"voca-store/internal/domain/dto/response"
	"voca-store/internal/domain/models"
	"voca-store/internal/helper"
	"voca-store/internal/service"
	"voca-store/internal/websocket"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	service service.ChatService
	hub     *websocket.Hub
}

func NewChatHandler(svc service.ChatService, hub *websocket.Hub) *ChatHandler {
	return &ChatHandler{service: svc, hub: hub}
}

// === HTTP HANDLERS ===

func (h *ChatHandler) CreateChatRequest(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	var req request.CreateChatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	session, err := h.service.CreateChatRequest(c.Request.Context(), userID, req.Message)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, "Chat request created successfully", response.ToChatSessionResponse(*session))
}

func (h *ChatHandler) AcceptChatRequest(c *gin.Context) {
	var req request.AcceptChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	sessionUID := c.Param("session_uid")

	session, err := h.service.AcceptChatRequest(c.Request.Context(), sessionUID, req.AdminID)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, "Chat request accepted", response.ToChatSessionResponse(*session))
}

func (h *ChatHandler) GetPendingChatRequests(c *gin.Context) {
	sessions, err := h.service.GetPendingChatRequests(c.Request.Context())
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch pending requests")
		return
	}

	var resp []response.ChatSessionResponse
	for _, s := range sessions {
		resp = append(resp, response.ToChatSessionResponse(s))
	}

	response.SuccessResponse(c, "Pending chat requests retrieved", helper.WrapListIfNeeded(resp))
}

func (h *ChatHandler) GetAllSessions(c *gin.Context) {
	sessions, err := h.service.GetAllChatSessions(c.Request.Context())
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch chat sessions")
		return
	}

	var resp []response.ChatSessionResponse
	for _, s := range sessions {
		resp = append(resp, response.ToChatSessionResponse(s))
	}

	response.SuccessResponse(c, "All chat sessions retrieved", helper.WrapListIfNeeded(resp))
}

func (h *ChatHandler) GetActiveSession(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	session, err := h.service.GetUserActiveSession(c.Request.Context(), userID)
	if err != nil {
		response.ErrorResponse(c, http.StatusNotFound, "No active chat session")
		return
	}
	response.SuccessResponse(c, "Active session retrieved", response.ToChatSessionResponse(*session))
}

func (h *ChatHandler) GetSessionByUID(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	role, _ := c.Get("role")
	sessionUID := c.Param("session_uid")

	session, err := h.service.GetSessionByUID(c.Request.Context(), sessionUID, userID, role.(string))
	if err != nil {
		response.ErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	response.SuccessResponse(c, "Session retrieved", response.ToChatSessionResponse(*session))
}

func (h *ChatHandler) GetChatHistory(c *gin.Context) {
	var req request.GetChatHistoryRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	session, err := h.service.GetSessionByUID(c.Request.Context(), req.SessionUID, 
		c.GetUint("user_id"), c.GetString("role"))
	if err != nil {
		response.ErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	messages, err := h.service.GetChatHistory(c.Request.Context(), session.ID, req.Limit, req.Offset)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch messages")
		return
	}

	var resp []response.ChatMessageResponse
	for _, msg := range messages {
		senderRole := "user"
		if session.AdminID != nil && msg.SenderID == *session.AdminID {
			senderRole = "admin"
		}
		resp = append(resp, response.ToChatMessageResponse(msg, senderRole))
	}

	response.SuccessResponse(c, "Chat history retrieved", helper.WrapListIfNeeded(resp))
}

func (h *ChatHandler) MarkMessagesRead(c *gin.Context) {
	sessionUID := c.Param("session_uid")
	userID := c.GetUint("user_id")

	if err := h.service.MarkMessagesRead(c.Request.Context(), sessionUID, userID); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, "Messages marked as read", nil)
}

func (h *ChatHandler) CloseSession(c *gin.Context) {
	sessionUID := c.Param("session_uid")
	closerID := c.GetUint("user_id")

	if err := h.service.CloseChatSession(c.Request.Context(), sessionUID, closerID); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, "Chat session closed", nil)
}

// === WEBSOCKET HANDLER (FIXED) ===
func (h *ChatHandler) WebSocketHandler(c *gin.Context) {
	sessionUID := c.Param("session_uid")
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	session, err := h.service.GetSessionByUID(c.Request.Context(), sessionUID, userID, role)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if session.Status != models.ChatSessionActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session is not active"})
		return
	}

	// === FIX: Gunakan exported Upgrader dari websocket package ===
	conn, err := websocket.UpgradeConnection(c.Writer, c.Request)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := websocket.NewClient(h.hub, conn, sessionUID, userID, role)
	
	// === FIX: Gunakan exported method RegisterClient ===
	h.hub.RegisterClient(client)

	if role == "Admin" {
		h.hub.RegisterAdmin(client)
	}

	go client.WritePump()
	go client.ReadPump()
}
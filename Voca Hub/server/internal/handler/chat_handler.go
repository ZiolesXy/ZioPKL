package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"server/internal/helper"
	"server/internal/service"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

func (h *ChatHandler) GetHistory(c *gin.Context) {
	otherID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}
	user := helper.MustCurrentUser(c)
	history, err := h.chatService.GetHistory(user.ID, uint(otherID))
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "chat history fetched", history)
}

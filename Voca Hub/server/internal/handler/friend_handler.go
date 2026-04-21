package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"server/internal/domain/dto"
	"server/internal/helper"
	"server/internal/service"
)

type FriendHandler struct {
	friendService *service.FriendService
}

func NewFriendHandler(friendService *service.FriendService) *FriendHandler {
	return &FriendHandler{friendService: friendService}
}

func (h *FriendHandler) AddFriend(c *gin.Context) {
	var req dto.AddFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user := helper.MustCurrentUser(c)
	if err := h.friendService.AddFriend(user.ID, req.FriendID); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	helper.Success(c, http.StatusCreated, "friend request sent", nil)
}

func (h *FriendHandler) AcceptFriend(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	user := helper.MustCurrentUser(c)
	if err := h.friendService.AcceptFriend(user.ID, uint(id)); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "friend request accepted", nil)
}

func (h *FriendHandler) RejectFriend(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	user := helper.MustCurrentUser(c)
	if err := h.friendService.RejectFriend(user.ID, uint(id)); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "friend request rejected", nil)
}

func (h *FriendHandler) ListFriends(c *gin.Context) {
	user := helper.MustCurrentUser(c)
	friends, err := h.friendService.ListFriends(user.ID)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "friends fetched", helper.WrapListIfNeeded(friends))
}

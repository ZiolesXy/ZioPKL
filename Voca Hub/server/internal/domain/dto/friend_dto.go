package dto

import "server/internal/domain/models"

type AddFriendRequest struct {
	FriendID uint `json:"friend_id" binding:"required"`
}

type PendingFriendRequestResponse struct {
	ID        uint        `json:"id"`
	UserID    uint        `json:"user_id"`
	FriendID  uint        `json:"friend_id"`
	Status    string      `json:"status"`
	Requester models.User `json:"requester"`
}

func BuildPendingFriendRequestResponses(relations []models.Friend) []PendingFriendRequestResponse {
	result := make([]PendingFriendRequestResponse, 0, len(relations))
	for _, relation := range relations {
		result = append(result, BuildPendingFriendRequestResponse(relation))
	}
	return result
}

func BuildPendingFriendRequestResponse(relation models.Friend) PendingFriendRequestResponse {
	return PendingFriendRequestResponse{
		ID:        relation.ID,
		UserID:    relation.UserID,
		FriendID:  relation.FriendID,
		Status:    relation.Status,
		Requester: relation.User,
	}
}

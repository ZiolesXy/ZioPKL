package dto

type AddFriendRequest struct {
	FriendID uint `json:"friend_id" binding:"required"`
}

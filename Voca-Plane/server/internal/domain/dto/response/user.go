package response

import (
	"time"
	"voca-plane/internal/domain/models"

	"gorm.io/gorm"
)

type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	IsBanned  bool      `json:"is_banned"`
	Reason    string    `json:"reason,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

func ToUserResponse(u models.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		IsBanned: u.IsBanned,
		Reason: u.BanReason,
		DeletedAt: u.DeletedAt,
	}
}

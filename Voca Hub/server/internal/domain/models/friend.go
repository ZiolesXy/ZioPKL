package models

type Friend struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UserID   uint   `gorm:"not null;index:idx_friend_pair,unique" json:"user_id"`
	FriendID uint   `gorm:"not null;index:idx_friend_pair,unique" json:"friend_id"`
	Status   string `gorm:"size:50;not null" json:"status"`
	User     User   `gorm:"foreignKey:UserID" json:"user"`
	Friend   User   `gorm:"foreignKey:FriendID" json:"friend"`
}

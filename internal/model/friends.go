package model

type FriendRequest struct {
	ID          uint64 `gorm:"primaryKey" json:"id"`
	SenderID    uint64
	RecipientID uint64
	Sender      *User `gorm:"foreignKey:SenderID" json:"sender"`
	Recipient   *User `gorm:"foreignKey:RecipientID" json:"recipient"`
	Pending     bool  `gorm:"default:true" json:"pending"`
	Accepted    bool  `gorm:"default:false" json:"accepted"`
}

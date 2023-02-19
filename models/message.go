package models

import (
	"time"
)

type Message struct {
	ID         int64     `json:"id,-"`
	ToUserID   int64     `json:"to_user_id,-"`
	FromUserID int64     `json:"from_user_id,-"`
	Content    string    `gorm:"type:varchar(999) not null" json:"content,-"`
	CreatedAt  time.Time `gorm:"column:create_time" json:"create_time,-"`
}

func GetMessagesById(userId int64, friendId int64) []Message {
	var messageList []Message
	Db.Where("to_user_id = ? AND from_user_id = ?", userId, friendId).
		Or("to_user_id = ? AND from_user_id = ?", friendId, userId).
		Find(&messageList)
	return messageList
}

func AddMessage(fromUserId int64, toUserId int64, content string) error {
	message := Message{
		ToUserID:   toUserId,
		FromUserID: fromUserId,
		Content:    content,
	}
	result := Db.Create(&message)
	return result.Error
}

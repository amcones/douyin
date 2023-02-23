package models

type Message struct {
	ID         int64  `json:"id,-"`
	ToUserID   int64  `gorm:"index:idx_user" json:"to_user_id,-"`
	FromUserID int64  `gorm:"index:idx_user" json:"from_user_id,-"`
	Content    string `gorm:"type:text not null" json:"content,-"`
	CreateTime int64  `gorm:"autoCreateTime" json:"create_time"`
}

func GetMessagesById(userId int64, friendId int64, preMsgTime int64) []Message {
	var messageList []Message
	Db.Where("to_user_id = ? AND from_user_id = ? AND create_time > ?", userId, friendId, preMsgTime).
		Or("to_user_id = ? AND from_user_id = ?  AND create_time > ?", friendId, userId, preMsgTime).
		Find(&messageList)
	return messageList
}

func AddMessage(fromUserId int64, toUserId int64, content string) error {
	newMessage := Message{
		ToUserID:   toUserId,
		FromUserID: fromUserId,
		Content:    content,
	}
	result := Db.Create(&newMessage)
	return result.Error
}

func (message *Message) GetLatestMessagesById(userId int64, friendId int64) {
	Db.Where("to_user_id = ? AND from_user_id = ?", userId, friendId).
		Or("to_user_id = ? AND from_user_id = ?", friendId, userId).
		Last(&message)
}

package model

type ChatRoom struct {
	RoomId string `json:"room_id" binding:"required"`
	Name   string `json:"name"`
}

func (ChatRoom) TableName() string {
	return "chat_room"
}

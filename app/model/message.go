package model

type Message struct {
	Id      uint64 `json:"id" binding:"required"`
	IsSelf  bool   `json:"is_self"`
	IsGroup bool   `json:"is_group"`
	Type    uint32 `json:"type"`
	Ts      uint32 `json:"ts"`
	RoomId  string `json:"room_id"`
	Content string `json:"content"`
	Sender  string `json:"sender"`
	Sign    string `json:"sign"`
	Thumb   string `json:"thumb"`
	Extra   string `json:"extra"`
	Xml     string `json:"xml"`
	Remark  string `json:"remark"`
}

func (Message) TableName() string {
	return "message"
}

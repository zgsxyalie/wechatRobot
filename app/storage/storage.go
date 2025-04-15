package storage

import (
	"os"
	"path/filepath"
	"strings"
	"time"
	"wechatRobot/app/config"
	"wechatRobot/app/model"
	"wechatRobot/grpc/wcf"
)

type Database struct{}

func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) SyncContacts(contacts map[string]Contact) {
	for wxId, con := range contacts {
		if strings.Contains(wxId, "@chatroom") {
			d.syncChatRoom(wxId, con.Name)
			continue
		}
		d.syncContact(wxId, con)
	}
}

func (d *Database) syncChatRoom(roomId, name string) {
	if model.DB.Where("room_id = ?", roomId).First(&model.ChatRoom{}).RowsAffected == 0 {
		model.DB.Create(&model.ChatRoom{RoomId: roomId, Name: name})
	}
}

func (d *Database) syncContact(wxId string, con Contact) {
	contact := model.Contact{
		WxId:     wxId,
		Code:     con.Code,
		Name:     con.Name,
		Remark:   con.Remark,
		Country:  con.Country,
		Province: con.Province,
		City:     con.City,
		Gender:   con.Gender,
	}

	if model.DB.Where("wx_id = ?", wxId).First(&model.Contact{}).RowsAffected == 0 {
		model.DB.Create(&contact)
	} else if con.Remark != "" {
		model.DB.Model(&model.Contact{}).Where("wx_id = ?", wxId).Updates(&model.Contact{Remark: con.Remark})
	}
}

func (d *Database) SaveMessage(msg *wcf.WxMsg) {
	model.DB.Create(&model.Message{
		Id:      msg.Id,
		RoomId:  msg.Roomid,
		Sender:  msg.Sender,
		Content: msg.Content,
		Sign:    msg.Sign,
		Xml:     msg.Xml,
		IsGroup: msg.IsGroup,
		IsSelf:  msg.IsSelf,
		Type:    msg.Type,
		Ts:      msg.Ts,
		Extra:   msg.Extra,
	})
}

func (d *Database) UpdateMessageImage(msgId uint64, image string) {
	model.DB.Model(&model.Message{}).Where("id = ?", msgId).Updates(&model.Contact{Remark: image})
}

type Contact struct {
	WxId     string
	Code     string
	Remark   string
	Name     string
	Country  string
	Province string
	City     string
	Gender   int32
}

type ContactManager struct {
	client *wcf.Client
}

func NewContactManager(client *wcf.Client) *ContactManager {
	return &ContactManager{client: client}
}

func (cm *ContactManager) Initialize() map[string]Contact {
	contacts := make(map[string]Contact, 500)
	for _, v := range cm.client.GetContacts() {
		if strings.Contains(v.Wxid, "@openim") || strings.Contains(v.Wxid, "gh_") {
			continue
		}
		contacts[v.Wxid] = Contact{
			WxId:     v.Wxid,
			Code:     v.Code,
			Remark:   v.Remark,
			Name:     v.Name,
			Country:  v.Country,
			Province: v.Province,
			City:     v.City,
			Gender:   v.Gender,
		}
	}
	return contacts
}

func (cm *ContactManager) GetContacts() map[string]Contact {
	return cm.Initialize()
}

type ChatHistory struct {
	msg     *wcf.WxMsg
	history []map[string]string
}

var userChatMap = make(map[string][]map[string]string)

func GetChatHistory(msg *wcf.WxMsg) *ChatHistory {
	key := msg.Roomid + msg.Sender
	if history, exists := userChatMap[key]; exists {
		return &ChatHistory{msg: msg, history: history}
	}
	return &ChatHistory{msg: msg, history: make([]map[string]string, 0)}
}

func (ch *ChatHistory) BuildMessages(content string) []map[string]string {
	if len(ch.history) == 0 {
		return []map[string]string{
			{"role": "system", "content": config.AppConfig.AI.Prompt},
			{"role": "user", "content": content},
		}
	}
	return append(ch.history, map[string]string{"role": "user", "content": content})
}

func (ch *ChatHistory) Append(role, content string) {
	ch.history = append(ch.history, map[string]string{"role": role, "content": content})
	userChatMap[ch.msg.Roomid+ch.msg.Sender] = ch.history
}

func (ch *ChatHistory) Clear() {
	delete(userChatMap, ch.msg.Roomid+ch.msg.Sender)
}

func GetStoragePath(msg *wcf.WxMsg) string {
	target, err := filepath.Abs("storage")
	if err != nil {
		if self, err := os.Executable(); err == nil {
			target = filepath.Dir(self)
		}
	}

	senderId := msg.Sender
	if msg.IsGroup {
		senderId = msg.Roomid
	}

	target = filepath.Join(target, "static", senderId, time.Now().Format("2006-01-02"))
	_ = os.MkdirAll(target, os.ModePerm)
	return target
}

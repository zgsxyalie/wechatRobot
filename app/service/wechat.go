package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"strings"
	"time"
	"wechatRobot/app/config"
	"wechatRobot/app/entity"
	"wechatRobot/app/sdk"
	"wechatRobot/app/storage"
	"wechatRobot/grpc/wcf"
)

type WeChatService struct {
	db     *storage.Database
	client *wcf.Client
}

func NewWeChatService(db *storage.Database) *WeChatService {
	return &WeChatService{
		db:     db,
		client: sdk.GetClient(),
	}
}

func (s *WeChatService) WaitForLogin() error {
	return sdk.WaitForLogin(s.client)
}

func (s *WeChatService) GetCurrentUser() string {
	return s.client.GetUserInfo().Name
}

func (s *WeChatService) InitializeContacts() {
	contacts := storage.NewContactManager(s.client).Initialize()
	s.db.SyncContacts(contacts)
}

func (s *WeChatService) StartMessageListener() {
	s.client.EnableRecvTxt()
	_ = s.client.OnMSG(func(msg *wcf.WxMsg) {
		if msg.IsSelf {
			return
		}

		content := s.sanitizeContent(msg.Content)
		go s.handleAutoReply(msg)
		go s.processMessage(msg, content)
	})
}

func (s *WeChatService) sanitizeContent(content string) string {
	if strings.Contains(content, "<?xml") && strings.Contains(content, "<msg>") {
		return "xml内容"
	}

	if strings.Contains(content, "<msg>") && strings.Contains(content, "<emoji") {
		return "表情包"
	}

	return content
}

func (s *WeChatService) processMessage(msg *wcf.WxMsg, content string) {
	contacts := storage.NewContactManager(s.client).GetContacts()
	s.logMessage(msg, content, contacts)
	s.db.SaveMessage(msg)

	if msg.Type == 3 {
		s.handleAttachment(msg)
	}
}

func (s *WeChatService) logMessage(msg *wcf.WxMsg, content string, contacts map[string]storage.Contact) {
	if contact, exists := contacts[msg.Sender]; exists && !msg.IsGroup {
		log.Printf("%s->%s:%s", contact.Name, contact.Remark, content)
	}
	if contact, exists := contacts[msg.Roomid]; exists && msg.IsGroup {
		log.Printf("%s->%s:%s", contact.Name, msg.Sender, content)
	}
}

func (s *WeChatService) handleAttachment(msg *wcf.WxMsg) {
	time.Sleep(2 * time.Second)
	s.downloadAttachment(msg)
	s.decryptImage(msg)
}

func (s *WeChatService) downloadAttachment(msg *wcf.WxMsg) {
	log.Printf("Downloading attachment: %s", msg.Extra)
	for i := 0; i < 3; i++ {
		if s.client.DownloadAttach(msg.Id, msg.Thumb, msg.Extra) == 0 {
			return
		}
		time.Sleep(time.Second)
		log.Println("Retrying attachment download...")
	}
}

func (s *WeChatService) decryptImage(msg *wcf.WxMsg) {
	target := storage.GetStoragePath(msg)
	for i := 0; i < 60; i++ {
		if image := s.client.DecryptImage(msg.Extra, target); image != "" {
			log.Printf("Download successful: %s", image)
			s.db.UpdateMessageImage(msg.Id, image)
			return
		}
		time.Sleep(time.Second)
		log.Println("Retrying image decryption...")
	}
}

func (s *WeChatService) handleAutoReply(msg *wcf.WxMsg) {
	if !strings.Contains(msg.Content, config.AppConfig.AI.Trigger) {
		return
	}

	content, err := s.getAIResponse(msg)
	if err != nil {
		s.sendErrorResponse(msg, err)
		return
	}

	s.sendReply(msg, content)
}

func (s *WeChatService) getAIResponse(msg *wcf.WxMsg) (string, error) {
	chatHistory := storage.GetChatHistory(msg)
	if strings.Contains(msg.Content, "清空会话") {
		chatHistory.Clear()
		return "已清空", nil
	}

	messages := chatHistory.BuildMessages(msg.Content)

	resp, err := resty.New().R().
		SetDebug(config.AppConfig.Debug).
		SetAuthToken(config.AppConfig.API.Key).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model":       config.AppConfig.API.Model,
			"messages":    messages,
			"temperature": config.AppConfig.API.Temperature,
		}).
		Post(config.AppConfig.API.URL)

	if err != nil {
		return "", err
	}

	var result entity.GptResp
	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "无法回答你的问题", nil
	}

	content := result.Choices[0].Message.Content
	storage.GetChatHistory(msg).Append("assistant", content)
	return content, nil
}

func (s *WeChatService) sendErrorResponse(msg *wcf.WxMsg, err error) {
	var restErr *resty.ResponseError
	if errors.As(err, &restErr) && restErr.Response.StatusCode() == 429 {
		s.client.SendTxt("限速了,明天再来吧", msg.Roomid, []string{msg.Sender})
	}
}

func (s *WeChatService) sendReply(msg *wcf.WxMsg, content string) {
	if msg.IsGroup {
		memberName := s.client.GetAliasInChatRoom(msg.Roomid, msg.Sender)
		s.client.SendTxt(fmt.Sprintf("@%s \n\n %s", memberName, content), msg.Roomid, []string{msg.Sender})
	} else {
		s.client.SendTxt(content, msg.Sender, []string{})
	}
}

package ws

import (
	"encoding/json"
)

// ChatRoom is a chat room
type Room struct {
	// room id
	ID string
	// room name
	Name string
	// room clients
	Clients map[string]*Client
}

type Message struct {
	// 消息类型
	Type string `json:"type"`
	// 消息内容
	Content string `json:"content"`
	// 消息发送者
	Sender string `json:"sender"`
	// 消息接收者
	Receiver string `json:"receiver"`
}

//描述的英文

func ResponseMessage(Type, str, sender, receiver string) []byte {
	msg := &Message{
		Type:     "400",
		Content:  str,
		Sender:   sender,
		Receiver: receiver,
	}
	if Type != "" {
		msg.Type = Type
	}

	msgByte, _ := json.Marshal(msg)
	return msgByte
}
func LogionSuccessMessage(uuid string) []byte {
	msg := &Message{
		Type:     "100",
		Content:  uuid,
		Sender:   "system",
		Receiver: uuid,
	}
	msgByte, _ := json.Marshal(msg)
	return msgByte
}

//返回在线用户列表
func LoginUserListMessage() []byte {
	content := ""
	for _, v := range manager.Clients {
		content += v.UserId.String() + ","
	}

	msg := &Message{
		Type:     "101",
		Content:  content,
		Sender:   "system",
		Receiver: "all",
	}
	msgByte, _ := json.Marshal(msg)
	return msgByte
}

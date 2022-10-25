package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func init() {
	log.Println("ws init")
	go manager.Start()
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsServe(c *gin.Context) {
	//upgrade := websocket.Upgrader{
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade ws:", err)
		c.JSON(400, gin.H{
			"code":    500,
			"message": "upgrade ws error",
		})
		return
	}
	defer ws.Close()
	//并发控制
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)

	defer cancel()

	_, msgByte, err := ws.ReadMessage()
	if err != nil {
		log.Println("read ws message error:", err)
		ws.WriteJSON("read ws message error:" + string(ResponseMessage("400", "读取ws消息错误", "system", "")))
		ws.Close()
		return
	}
	log.Println("msg:", string(msgByte))
	message := &Message{}
	err = json.Unmarshal(msgByte, message)
	if err != nil {
		log.Println("unmarshal msg err:", err)
		ws.WriteJSON("unmarshal msg:" + string(ResponseMessage("400", "解析消息错误", "system", "")))
		ws.Close()
		return
	}
	//消息类型
	client := &Client{}
	if message.Type == "100" {
		// 登录
		// 生成用户id
		client.UserId = uuid.New()
		client.UserName = message.Content
		client.Conn = ws
		client.Send = make(chan []byte, 100)
		log.Println("client Login:", client.UserId)
		// 注册用户
		manager.Register <- client
		log.Println("回显消息")
		err := ws.WriteMessage(websocket.TextMessage, LogionSuccessMessage(client.UserId.String()))
		if err != nil {
			log.Println("write message error:", err)
			ws.Close()
			return
		}
	} else {
		log.Println("client not login")
		ws.WriteJSON("message type error:" + string(ResponseMessage("400", "消息类型错误", "system", "")))
		ws.Close()
		return
	}

	go client.Read(ctx, cancel)
	// 发送消息
	go client.Write(ctx, cancel)

	<-ctx.Done()
	log.Println("ws serve done")
	manager.Unregister <- client
	ws.Close()
}

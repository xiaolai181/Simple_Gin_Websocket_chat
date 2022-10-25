package ws

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	// 用户id
	UserId uuid.UUID
	// 用户名
	UserName string
	// 用户ws连接
	Conn *websocket.Conn
	// 发送消息
	Send chan []byte
	// 接收消息
	Receive chan []byte
}

func (c *Client) Read(ctx context.Context, cancel context.CancelFunc) {
	defer func() {
		// 注销用户
		log.Println("send关闭:", c.UserId)
		manager.Unregister <- c
		cancel()
	}()
	log.Println("read start")
	for {
		select {
		case <-ctx.Done():
			log.Println("read conn close")
			return
		default:
			_, msgByte, err := c.Conn.ReadMessage()
			if err != nil {
				log.Println("read ws message error:", err)
				return
			}
			if msgByte != nil {
				// byte to Message
				msg := &Message{}
				err = json.Unmarshal(msgByte, msg)
				if err != nil {
					log.Println("read  unmarshal msg:", err)
					continue
				}
				//msg.Receiver to uuid.UUID
				log.Println("read Receiver:", msg.Receiver)
				receiver, err := uuid.Parse(msg.Receiver)
				if err != nil || receiver == uuid.Nil {
					log.Println("read  parse uuid:", err)
					c.Send <- []byte(ResponseMessage("500", "接收者id错误", "system", c.UserId.String()))
					continue
				}
				//send to receiver
				if _, found := manager.Clients[receiver]; found {
					log.Println("Send 接收者在线:", c.UserId.String())
					manager.Clients[receiver].Send <- []byte(ResponseMessage("200", msg.Content, c.UserId.String(), receiver.String()))
				} else {
					c.Send <- []byte(ResponseMessage("500", "接收者不在线"+receiver.String(), "system", c.UserId.String()))
					continue
				}
			}
		}
	}
}

func (c *Client) Write(ctx context.Context, cancel context.CancelFunc) {
	defer func() {
		log.Println("write关闭:", c.UserId)
		manager.Unregister <- c
		cancel()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				log.Println("send channel close，read message error")
				continue
			}
			log.Println("send:", string(message))
			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("write message error:", err)
				return
			}
		case <-ctx.Done():
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			log.Println("wirte goroutine conn close")
			return
		}
	}
}

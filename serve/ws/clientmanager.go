package ws

import (
	"context"
	"log"

	"github.com/google/uuid"
)

type ClientManager struct {
	// 所有用户
	Clients map[uuid.UUID]*Client
	// 注册用户
	Register chan *Client
	// 注销用户
	Unregister chan *Client
	// 广播消息
	Broadcast chan []byte
	// 广播在线用户的信号
	BroadcastUser chan int
}

var manager = &ClientManager{
	Clients:       make(map[uuid.UUID]*Client),
	Broadcast:     make(chan []byte),
	BroadcastUser: make(chan int),
	Register:      make(chan *Client, 10),
	Unregister:    make(chan *Client, 10),
}

func (m *ClientManager) Start() {
	log.Println("ws start")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go register(ctx)
	go unregister(ctx)
	go BroadcastOnlineUser(ctx)
	<-ctx.Done()

}

//上线注册
func register(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn := <-manager.Register
			// 注册用户
			log.Println("注册:", conn.UserId)
			manager.Clients[conn.UserId] = conn
			manager.BroadcastUser <- 1
		}
	}
}

//注销登录
func unregister(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn := <-manager.Unregister
			// 注销用户
			if _, ok := manager.Clients[conn.UserId]; ok {
				log.Println("开始注销:", conn.UserId)
				delete(manager.Clients, conn.UserId)
				manager.BroadcastUser <- 1
				log.Println("注销成功:", conn.UserId)
			}
		}
	}

}

//广播当前在线用户
func BroadcastOnlineUser(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-manager.BroadcastUser:
			//广播在线用户
			log.Println("广播在线用户")
			for _, conn := range manager.Clients {
				conn.Send <- LoginUserListMessage()
			}
		}
	}
}

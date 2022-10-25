package main

import (
	"wsgin/serve/ws"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLFiles("index.html")
	r.GET("/", Iendx)
	r.GET("/ws", ws.WsServe)
	r.Run(":8080")
}

func Iendx(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

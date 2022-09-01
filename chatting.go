package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	Users = map[string]*websocket.Conn{}
)

func SocketSrv(c *gin.Context) {
	wsconn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("SocketSrv1", err)
		return
	}
	defer wsconn.Close()
	userID := ""
	err = wsconn.WriteMessage(1, []byte(`Masukkan ID yang kamu inginkan`))
	if err != nil {
		log.Println("SocketSrv2", err)
		return
	}
	for {
		mt, message, err := wsconn.ReadMessage()
		if err != nil {
			log.Println("SocketSrv3", err)
			break
		}
		if userID == "" {
			userID = string(message)
			Users[userID] = wsconn
			wsconn.WriteMessage(mt, []byte(fmt.Sprintf(`Selamat datang %s, sekarang anda sudah bisa menggunakan layanan ini`, userID)))
		} else {
			for key, value := range Users {
				if key == userID {
					continue
				}
				value.WriteMessage(mt, []byte(fmt.Sprintf(`[%s] %s`, userID, string(message))))
			}
		}
	}
	delete(Users, userID)
}

func main() {
	router := gin.Default()
	router.GET("/ws", SocketSrv)
	router.Run(":8083")
}

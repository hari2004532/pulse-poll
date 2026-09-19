package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")

		return origin == "http://localhost:5173"
	},
}

func HandleConnection(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		pollID := c.Param("id")

		conn, err := upgrader.Upgrade(
			c.Writer,
			c.Request,
			nil,
		)

		if err != nil {
			return
		}

		client := &Client{
			Conn:   conn,
			PollID: pollID,
		}

		hub.AddClient(client)

		defer hub.RemoveClient(client)

		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}
}

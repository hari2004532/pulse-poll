package websocket

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

func NewUpgrader(frontendURL string) websocket.Upgrader {
    return websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool {
            origin := r.Header.Get("Origin")

            return origin == frontendURL
        },
    }
}

func HandleConnection(hub *Hub, frontendURL string) gin.HandlerFunc {
    upgrader := NewUpgrader(frontendURL)

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

// chat-server/handlers/websocket.go
package handlers

import (
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
	"layer7/chat-server/models"

)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // Allowing all origins for development
    },
}

type WSHandler struct {
    room *models.Room
}

func NewWSHandler() *WSHandler {
    return &WSHandler{
        room: models.NewRoom(),
    }
}

func (h *WSHandler) HandleConnections(c *gin.Context) {
    ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Printf("Error upgrading connection: %v", err)
        return
    }

    username := c.Query("username")
    if username == "" {
        username = "anonymous"
    }

    client := &models.Client{
        Username: username,
        Send:     make(chan *models.Message, 256),
    }

    h.room.Register <- client

    // Notify others that user joined
    h.room.Broadcast <- &models.Message{
        Type:     "join",
        Username: username,
        Content:  "joined the chat",
    }

    go h.writePump(client, ws)
    h.readPump(client, ws)
}

func (h *WSHandler) readPump(client *models.Client, ws *websocket.Conn) {
    defer func() {
        h.room.Unregister <- client
        ws.Close()
    }()

    for {
        var msg models.Message
        err := ws.ReadJSON(&msg)
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
                log.Printf("error: %v", err)
            }
            break
        }
        msg.Username = client.Username
        h.room.Broadcast <- &msg
    }
}

func (h *WSHandler) writePump(client *models.Client, ws *websocket.Conn) {
    defer ws.Close()

    for {
        message, ok := <-client.Send
        if !ok {
            ws.WriteMessage(websocket.CloseMessage, []byte{})
            return
        }

        err := ws.WriteJSON(message)
        if err != nil {
            return
        }
    }
}

func (h *WSHandler) Run() {
    for {
        select {
        case client := <-h.room.Register:
            h.room.Clients[client] = true
            log.Printf("Client registered: %s", client.Username)

        case client := <-h.room.Unregister:
            if _, ok := h.room.Clients[client]; ok {
                delete(h.room.Clients, client)
                close(client.Send)
                h.room.Broadcast <- &models.Message{
                    Type:     "leave",
                    Username: client.Username,
                    Content:  "left the chat",
                }
            }

        case message := <-h.room.Broadcast:
            for client := range h.room.Clients {
                client.Send <- message
            }
        }
    }
}
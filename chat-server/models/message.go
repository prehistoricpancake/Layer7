// chat-server/models/message.go
package models

type Message struct {
    Type     string `json:"type"`     // "chat", "join", "leave"
    Content  string `json:"content"`
    Username string `json:"username"`
}

type Client struct {
    Username string
    Send     chan *Message
}

type Room struct {
    Clients    map[*Client]bool
    Broadcast  chan *Message
    Register   chan *Client
    Unregister chan *Client
}

func NewRoom() *Room {
    return &Room{
        Clients:    make(map[*Client]bool),
        Broadcast:  make(chan *Message),
        Register:   make(chan *Client),
        Unregister: make(chan *Client),
    }
}
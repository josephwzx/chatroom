package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		type TempMessage struct {
			Username string `json:"username"`
			Message  string `json:"message"`
		}
		var temp TempMessage
		var msg Message
		if err := json.Unmarshal(p, &temp); err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}
		msg.Sender = temp.Username
		msg.Content = temp.Message
		// time is current
		msg.CreatedAt = time.Now()

		fmt.Printf("Message Received: %+v\n", msg)

		var id int64
		id, err = SaveMessage(msg.Content, msg.Sender)
		msg.ID = id
		if err != nil {
			log.Printf("Error saving message: %v", err)
			continue
		}

		// Broadcast message to all clients in the pool
		c.Pool.Broadcast <- msg
	}
}

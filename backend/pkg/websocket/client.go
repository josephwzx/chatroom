package websocket

import (
	"fmt"
	"log"

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
		messageType, p, err := c.Conn.ReadMessage()
		log.Printf("Message Type: %v", messageType)
		if err != nil {
			log.Println(err)
			return
		}
		// Assume Body of Message is type string. If it's not, you will need to adjust.
		messageContent := string(p)
		fmt.Printf("Message Received: %+v\n", messageContent)

		// Save message to database, using the exported SaveMessage function
		err = SaveMessage(messageContent, c.ID)
		if err != nil {
			log.Printf("Error saving message from client %v: %v", c.ID, err)
			continue // Or handle the error as you see fit
		}

		// Broadcast message to all clients in the pool
		c.Pool.Broadcast <- Message{Content: messageContent}
	}
}

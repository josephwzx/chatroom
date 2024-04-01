package websocket

import (
	"database/sql"
	"log"
	"time"
)

// Db is your database connection, initialized elsewhere in your application.
var Db *sql.DB

func SetDatabaseConnection(database *sql.DB) {
	Db = database
}

// Message struct represents a chat message
type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Sender    string    `json:"sender"`
	CreatedAt time.Time `json:"created_at"`
}

// SaveMessage stores a message in the database
func SaveMessage(message, sender string) error {
	_, err := Db.Exec("INSERT INTO messages (content, sender) VALUES ($1, $2)", message, sender)
	if err != nil {
		log.Printf("Error saving message to database: %v", err)
		return err
	}
	return nil
}

// GetChatHistory retrieves the chat history from the database
func GetChatHistory() ([]Message, error) {
	var history []Message

	rows, err := Db.Query("SELECT id, content, sender, created_at FROM messages")
	if err != nil {
		log.Printf("Error retrieving chat history: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.Content, &msg.Sender, &msg.CreatedAt); err != nil {
			log.Printf("Error scanning message: %v", err)
			return nil, err
		}
		history = append(history, msg)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
		return nil, err
	}

	return history, nil
}

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/josephwzx/chatroom/pkg/websocket"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	var err error
	connStr := "user=joseph dbname=chatroom sslmode=disable password=zx240915"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	websocket.SetDatabaseConnection(db)
}

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*") // Adjust accordingly for security in production
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})

	http.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/history endpoint hit") // Debug logging
		enableCors(&w)
		history, err := websocket.GetChatHistory()
		if err != nil {
			log.Printf("Error retrieving chat history: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(history); err != nil {
			log.Printf("Error encoding chat history to JSON: %v", err)
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}
	})

}

func main() {
	fmt.Println("Chatroom backend started!")
	initDB()
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}

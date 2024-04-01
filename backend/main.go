package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/josephwzx/chatroom/pkg/auth"
	"github.com/josephwzx/chatroom/pkg/vote"
	"github.com/josephwzx/chatroom/pkg/websocket"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	var err error
	connStr := "host=10.0.0.122 user=joseph dbname=chatroom sslmode=disable password=zx240915"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	auth.InitializeAuth(db)
	websocket.SetDatabaseConnection(db)
	vote.SetDatabaseConnection(db)
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
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func voteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, err := auth.AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var v vote.Vote
	err = json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	var userID = claims.Username
	v.UserID = userID

	fmt.Printf("Vote: %+v\n", v)

	var upvoteCount, downvoteCount int
	upvoteCount, downvoteCount, err = vote.CastVote(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// to-be-modified
	response := map[string]interface{}{
		"message":        "Vote successfully recorded",
		"upvote_count":   upvoteCount,
		"downvote_count": downvoteCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})

	http.HandleFunc("/login", corsMiddleware(auth.LoginUser))
	http.HandleFunc("/register", corsMiddleware(auth.RegisterUser))
	http.HandleFunc("/vote", corsMiddleware(voteHandler))

	http.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/history endpoint hit")
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

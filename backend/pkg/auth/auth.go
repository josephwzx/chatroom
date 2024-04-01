package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB // Assume db is initialized and connected from outside

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func InitializeAuth(database *sql.DB) {
	db = database
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", creds.Username, string(hashedPassword))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var storedCreds Credentials
	row := db.QueryRow("SELECT username, password_hash FROM users WHERE username=$1", creds.Username)
	err = row.Scan(&storedCreds.Username, &storedCreds.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := GenerateToken(storedCreds.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

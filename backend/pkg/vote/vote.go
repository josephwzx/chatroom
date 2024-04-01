package vote

import (
	"database/sql"
	"log"
)

var db *sql.DB

func SetDatabaseConnection(database *sql.DB) {
	db = database
}

type VoteType string

const (
	Upvote   VoteType = "upvote"
	Downvote VoteType = "downvote"
)

type Vote struct {
	ID        int64  `json:"id"`
	UserID    string `json:"user_id"`
	MessageID int64  `json:"message_id"`
	VoteType  string `json:"vote_type"`
}

func CastVote(v Vote) error {
	// Start a new transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return err
	}

	// Lock the message row to prevent concurrent updates
	var currentVersion, upvoteCount, downvoteCount int
	err = tx.QueryRow("SELECT version, upvotecount, downvotecount FROM messages WHERE id = $1 FOR UPDATE", v.MessageID).Scan(&currentVersion, &upvoteCount, &downvoteCount)
	if err != nil {
		log.Printf("Error locking message row or row not found: %v", err)
		tx.Rollback()
		return err
	}

	// Check for an existing vote by the user on the message
	var existingVoteType string
	err = tx.QueryRow("SELECT votetype FROM votes WHERE userid = $1 AND messageid = $2", v.UserID, v.MessageID).Scan(&existingVoteType)

	if err == sql.ErrNoRows {
		// Insert new vote
		_, err = tx.Exec("INSERT INTO votes (userid, messageid, votetype) VALUES ($1, $2, $3)", v.UserID, v.MessageID, v.VoteType)
		if err != nil {
			log.Println("Error inserting vote:", err)
			tx.Rollback()
			return err
		}
		// Adjust message counts based on the new vote
		if v.VoteType == string(Upvote) {
			upvoteCount++
		} else {
			downvoteCount++
		}
	} else if err == nil && existingVoteType != v.VoteType {
		// Update existing vote if the type has changed
		_, err = tx.Exec("UPDATE votes SET votetype = $1 WHERE userid = $2 AND messageid = $3", v.VoteType, v.UserID, v.MessageID)
		if err != nil {
			log.Println("Error updating existing vote:", err)
			tx.Rollback()
			return err
		}
		// Correctly adjust message vote counts based on the vote change
		if existingVoteType == string(Upvote) && v.VoteType == string(Downvote) {
			upvoteCount--
			downvoteCount++
		} else if existingVoteType == string(Downvote) && v.VoteType == string(Upvote) {
			downvoteCount--
			upvoteCount++
		}
	} else if err != nil {
		log.Println("Error querying existing vote:", err)
		tx.Rollback()
		return err
	}

	// Finally, update the message's vote counts and version
	_, err = tx.Exec("UPDATE messages SET upvotecount = $1, downvotecount = $2, version = version + 1 WHERE id = $3 AND version = $4", upvoteCount, downvoteCount, v.MessageID, currentVersion)
	if err != nil {
		log.Printf("Error updating message vote count: %v", err)
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	log.Println("Vote processed successfully.")
	return nil
}

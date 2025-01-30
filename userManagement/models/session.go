package models

import (
	"database/sql"
	"errors"
	"forum/utils"
	"log"
	"time"
)

// User struct represents the user data model
type Session struct {
	ID           int       `json:"id"`
	SessionToken string    `json:"session_token"`
	UserId       int       `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func InsertSession(session *Session) (*Session, error) {
	db := utils.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Generate UUID for the user if not already set
	if session.SessionToken == "" {
		uuidSessionTokenid, err := utils.GenerateUuid()
		if err != nil {
			return nil, err
		}
		session.SessionToken = uuidSessionTokenid
	}

	// Set session expiration time
	session.ExpiresAt = time.Now().Add(12 * time.Hour)

	insertQuery := `INSERT INTO sessions (session_token, user_id, expires_at) VALUES (?, ?, ?);`
	_, insertErr := db.Exec(insertQuery, session.SessionToken, session.UserId, session.ExpiresAt)
	if insertErr != nil {
		// Check if the error is a SQLite constraint violation
		if sqliteErr, ok := insertErr.(interface{ ErrorCode() int }); ok {
			if sqliteErr.ErrorCode() == 19 { // SQLite constraint violation error code
				return nil, sql.ErrNoRows // Return custom error to indicate a duplicate
			}
		}
		return nil, insertErr
	}
	return session, nil
}

func SelectSession(sessionToken string) (User, time.Time, error) {
	db := utils.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	var user User
	var expirationTime time.Time
	err := db.QueryRow(`SELECT 
							u.id as user_id, u.type as user_type, u.name as user_name, u.username as username, u.email as user_email, 
							expires_at 
						FROM sessions s
							INNER JOIN users u
								ON s.user_id = u.id
						WHERE session_token = ?`, sessionToken).Scan(&user.ID, &user.Type, &user.Name, &user.Username, &user.Email, &expirationTime)
	if err != nil {
		// Handle other database errors
		log.Fatal(err)
		return User{}, time.Time{}, errors.New("database error")
	}

	return user, expirationTime, nil
}

func DeleteSession(sessionToken string) error {

	db := utils.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes
	_, err := db.Exec(`UPDATE sessions
					SET expires_at = CURRENT_TIMESTAMP
					WHERE session_token = ?;`, sessionToken)
	if err != nil {
		// Handle other database errors
		log.Fatal(err)
		return errors.New("database error")
	}

	return nil

}

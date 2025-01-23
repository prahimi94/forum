package userManagement

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

func openDBConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func insertUser(username, email, password string) (int, error) {
	db := openDBConnection()
	defer db.Close() // Close the connection after the function finishes

	user_uuid, err := generateUuid()
	if err != nil {
		return -1, err
	}

	insertQuery := `INSERT INTO users (uuid, name, username, email, password) VALUES (?, ?, ?, ?, ?);`
	_, err = db.Exec(insertQuery, user_uuid, username, username, email, password)
	if err != nil {
		// Check if the error is a SQLite constraint violation
		if sqliteErr, ok := err.(interface{ ErrorCode() int }); ok {
			if sqliteErr.ErrorCode() == 19 { // SQLite constraint violation error code
				return -1, sql.ErrNoRows // Return custom error to indicate a duplicate
			}
		}
		fmt.Println("error is here")
		return -1, err
	}

	var userId int
	err = db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userId)
	if err != nil {
		// Handle other database errors
		log.Fatal(err)
	}

	return userId, nil
}

func authenticateUser(username, password string) (bool, int, error) {
	// Open SQLite database
	db := openDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Query to retrieve the hashed password stored in the database for the given username
	var userId int
	var storedHashedPassword string
	err := db.QueryRow("SELECT id, password FROM users WHERE username = ?", username).Scan(&userId, &storedHashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			// Username not found
			return false, -1, errors.New("username not found")
		}
		// Handle other database errors
		log.Fatal(err)
	}

	// Compare the entered password with the stored hashed password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password))
	if err != nil {
		// Password is incorrect
		return false, -1, errors.New("password is incorrect")
	}

	// Successful login if no errors occurred
	return true, userId, nil
}

func sessionInsert(userId int) (string, time.Time, error) {
	db := openDBConnection()
	defer db.Close() // Close the connection after the function finishes

	sessionToken, err := generateUuid()
	if err != nil {
		return "", time.Time{}, err
	}

	// Set session expiration time
	expirationTime := time.Now().Add(12 * time.Hour)

	insertQuery := `INSERT INTO sessions (session_token, user_id, expires_at) VALUES (?, ?, ?);`
	_, err = db.Exec(insertQuery, sessionToken, userId, expirationTime)
	if err != nil {
		// Check if the error is a SQLite constraint violation
		if sqliteErr, ok := err.(interface{ ErrorCode() int }); ok {
			if sqliteErr.ErrorCode() == 19 { // SQLite constraint violation error code
				return "", time.Time{}, sql.ErrNoRows // Return custom error to indicate a duplicate
			}
		}
		return "", time.Time{}, err
	}
	return sessionToken, expirationTime, nil
}

func generateUuid() (string, error) {
	// Create a Version 4 UUID.
	u2, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
		return "", err
	}
	log.Printf("generated Version 4 UUID %v", u2)

	return u2.String(), nil
}

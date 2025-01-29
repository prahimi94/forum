package models

import (
	"database/sql"
	"errors"
	"forum/utils"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User struct represents the user data model
type User struct {
	ID        int        `json:"id"`
	UUID      string     `json:"uuid"`
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy *int       `json:"updated_by"`
}

func InsertUser(user *User) (int, error) {
	db := utils.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Generate UUID for the user if not already set
	if user.UUID == "" {
		uuid, err := utils.GenerateUuid()
		if err != nil {
			return -1, err
		}
		user.UUID = uuid
	}

	insertQuery := `INSERT INTO users (uuid, name, username, email, password) VALUES (?, ?, ?, ?, ?);`
	result, insertErr := db.Exec(insertQuery, user.UUID, user.Username, user.Username, user.Email, user.Password)
	if insertErr != nil {
		// Check if the error is a SQLite constraint violation
		if sqliteErr, ok := insertErr.(interface{ ErrorCode() int }); ok {
			if sqliteErr.ErrorCode() == 19 { // SQLite constraint violation error code
				return -1, sql.ErrNoRows // Return custom error to indicate a duplicate
			}
		}
		return -1, insertErr
	}

	// Retrieve the last inserted ID
	userId, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return -1, err
	}

	return int(userId), nil
}

func AuthenticateUser(username, password string) (bool, int, error) {
	// Open SQLite database
	db := utils.OpenDBConnection()
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

package userManagement

import (
	"database/sql"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func openDBConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	return db
}

func insertUser(username, email, password string) error {
	db := openDBConnection()
	insertQuery := `INSERT INTO users (name, username, email, password) VALUES (? ,?, ?, ?);`
	_, err := db.Exec(insertQuery, username, username, email, password)
	if err != nil {
		// Check if the error is a SQLite constraint violation
		if sqliteErr, ok := err.(interface{ ErrorCode() int }); ok {
			if sqliteErr.ErrorCode() == 19 { // SQLite constraint violation error code
				return sql.ErrNoRows // Return custom error to indicate a duplicate
			}
		}
		return err
	}
	return nil
}

func authenticateUser(username, password string) (bool, error) {
	// Open SQLite database
	db := openDBConnection()

	// Query to retrieve the hashed password stored in the database for the given username
	var storedHashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedHashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			// Username not found
			return false, errors.New("Username not found")
		}
		// Handle other database errors
		log.Fatal(err)
	}

	// Compare the entered password with the stored hashed password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password))
	if err != nil {
		// Password is incorrect
		return false, errors.New("Password is incorrect")
	}

	// Successful login if no errors occurred
	return true, nil
}

package models

import (
	"database/sql"
	"forum/utils"
	"log"
	"time"
)

// Post struct represents the user data model
type Post struct {
	ID          int       `json:"id"`
	UUID        string    `json:"uuid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserId      int       `json:"user_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   int       `json:"updated_by"`
}

func openDBConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	if err != nil {
		log.Fatal(err)
	}

	// Enable foreign key constraints
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Failed to enable foreign key constraints:", err)
	}

	return db
}

func InsertPost(post *Post, categoryIds []int) (int, error) {
	db := openDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Start a transaction for atomicity
	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}

	post.UUID, err = utils.GenerateUuid()
	if err != nil {
		tx.Rollback() // Rollback if UUID generation fails
		return -1, err
	}

	insertQuery := `INSERT INTO posts (uuid, title, description, user_id) VALUES (?, ?, ?, ?);`
	result, insertErr := tx.Exec(insertQuery, post.UUID, post.Title, post.Description, post.UserId)
	if insertErr != nil {
		tx.Rollback() // Rollback on error
		// Check if the error is a SQLite constraint violation
		if sqliteErr, ok := insertErr.(interface{ ErrorCode() int }); ok {
			if sqliteErr.ErrorCode() == 19 { // SQLite constraint violation error code
				return -1, sql.ErrNoRows // Return custom error to indicate a duplicate
			}
		}
		return -1, insertErr
	}

	// Retrieve the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback() // Rollback on error
		log.Fatal(err)
		return -1, err
	}

	insertPostCategoriesErr := InsertPostCategories(int(lastInsertID), categoryIds, post.UserId, tx)
	if insertPostCategoriesErr != nil {
		return -1, insertPostCategoriesErr
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return -1, err
	}

	return int(lastInsertID), nil
}

func UpdatePost(post *Post, categories []int, user_id int) error {
	db := openDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Start a transaction for atomicity
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	updateQuery := `UPDATE posts
					SET title = ?,
						description = ?,
						updated_at = CURRENT_TIMESTAMP,
						updated_by = ?
					WHERE id = ?;`
	_, updateErr := tx.Exec(updateQuery, post.Title, post.Description, user_id, post.ID)
	if updateErr != nil {
		tx.Rollback() // Rollback on error
		// Check if the error is a SQLite constraint violation
		if sqliteErr, ok := updateErr.(interface{ ErrorCode() int }); ok {
			if sqliteErr.ErrorCode() == 19 { // SQLite constraint violation error code
				return sql.ErrNoRows // Return custom error to indicate a duplicate
			}
		}
		return updateErr
	}

	deletePostCategoriesErr := UpdateStatusPostCategories(post.ID, user_id, "delete", tx)
	if deletePostCategoriesErr != nil {
		return deletePostCategoriesErr
	}

	insertPostCategoriesErr := InsertPostCategories(post.ID, categories, user_id, tx)
	if insertPostCategoriesErr != nil {
		return insertPostCategoriesErr
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func UpdateStatusPost(post_id int, status string, user_id int) error {
	db := openDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Start a transaction for atomicity
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	updateQuery := `UPDATE posts
					SET status = ?,
						updated_at = CURRENT_TIMESTAMP,
						updated_by = ?
					WHERE id = ?;`
	_, updateErr := tx.Exec(updateQuery, status, user_id, post_id)
	if updateErr != nil {
		tx.Rollback() // Rollback on error
		// Check if the error is a SQLite constraint violation
		if sqliteErr, ok := updateErr.(interface{ ErrorCode() int }); ok {
			if sqliteErr.ErrorCode() == 19 { // SQLite constraint violation error code
				return sql.ErrNoRows // Return custom error to indicate a duplicate
			}
		}
		return updateErr
	}

	updateStatusPostCategories := UpdateStatusPostCategories(post_id, user_id, status, tx)
	if updateStatusPostCategories != nil {
		return updateStatusPostCategories
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func ReadPosts() {

}

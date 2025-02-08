package models

import (
	"database/sql"
	"time"
)

// Post File struct represents the user data model
type PostFile struct {
	ID               int       `json:"id"`
	PostId           int       `json:"post_id"`
	FileUploadedName *string   `json:"file_uploaded_name"`
	FileRealName     *string   `json:"file_real_name"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        int       `json:"created_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpdatedBy        int       `json:"updated_by"`
}

func InsertPostFiles(post_id int, uploadedFiles map[string]string, user_id int, tx *sql.Tx) error {
	// Prepare the bulk insert query for post_categories
	if len(uploadedFiles) > 0 {
		query := `INSERT INTO post_files (post_id, file_real_name, file_uploaded_name, created_by) VALUES `
		values := make([]any, 0, len(uploadedFiles)*3)

		for i := 0; i < len(uploadedFiles); i++ {
			if i > 0 {
				query += ", "
			}
			query += "(?, ?, ?, ?)"
			for key, value := range uploadedFiles {
				values = append(values, post_id, key, value, user_id)
			}
		}
		query += ";"

		// Execute the bulk insert query
		_, err := tx.Exec(query, values...)
		if err != nil {
			tx.Rollback() // Rollback on error
			return err
		}
	}
	return nil
}

func UpdateStatusPostFiles(post_id int, user_id int, status string, tx *sql.Tx) error {
	updateStatusQuery := `UPDATE post_files
					SET status = ?,
						updated_at = CURRENT_TIMESTAMP,
						updated_by = ?
					WHERE post_id = ?
					AND status != 'delete';`
	_, updateStatusErr := tx.Exec(updateStatusQuery, status, user_id, post_id)
	if updateStatusErr != nil {
		tx.Rollback() // Rollback on error
		// Check if the error is a SQLite constraint violation
		if sqliteErr, ok := updateStatusErr.(interface{ ErrorCode() int }); ok {
			if sqliteErr.ErrorCode() == 19 { // SQLite constraint violation error code
				return sql.ErrNoRows // Return custom error to indicate a duplicate
			}
		}
		return updateStatusErr
	}

	return nil
}

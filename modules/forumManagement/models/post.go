package models

import (
	"database/sql"
	"fmt"
	"forum/db"
	userManagementModels "forum/modules/userManagement/models"
	"forum/utils"
	"log"
	"sort"
	"time"
)

// Post struct represents the user data model
type Post struct {
	ID               int                       `json:"id"`
	UUID             string                    `json:"uuid"`
	Title            string                    `json:"title"`
	Description      string                    `json:"description"`
	UserId           int                       `json:"user_id"`
	Status           string                    `json:"status"`
	CreatedAt        time.Time                 `json:"created_at"`
	UpdatedAt        *time.Time                `json:"updated_at"`
	UpdatedBy        *int                      `json:"updated_by"`
	IsLikedByUser    bool                      `json:"liked"`
	IsDislikedByUser bool                      `json:"disliked"`
	NumberOfLikes    int                       `json:"number_of_likes"`
	NumberOfDislikes int                       `json:"number_of_dislikes"`
	User             userManagementModels.User `json:"user"`       // Embedded user data
	Categories       []Category                `json:"categories"` // List of categories related to the post
	PostFiles        []PostFile                `json:"post_files"` // List of files related to the post
}

func InsertPost(post *Post, categoryIds []int, uploadedFilesAddresses []string) (int, error) {
	db := db.OpenDBConnection()
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
		tx.Rollback() // Rollback on error
		return -1, insertPostCategoriesErr
	}

	insertPostFilesErr := InsertPostFiles(int(lastInsertID), uploadedFilesAddresses, post.UserId, tx)
	if insertPostFilesErr != nil {
		tx.Rollback() // Rollback on error
		return -1, insertPostFilesErr
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback() // Rollback on error
		return -1, err
	}

	return int(lastInsertID), nil
}

func UpdatePost(post *Post, categories []int, uploadedFilesAddresses []string, user_id int) error {
	db := db.OpenDBConnection()
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
		tx.Rollback() // Rollback on error
		return deletePostCategoriesErr
	}

	deletePostFilesErr := UpdateStatusPostFiles(post.ID, user_id, "delete", tx)
	if deletePostFilesErr != nil {
		tx.Rollback() // Rollback on error
		return deletePostFilesErr
	}

	insertPostCategoriesErr := InsertPostCategories(post.ID, categories, user_id, tx)
	if insertPostCategoriesErr != nil {
		tx.Rollback() // Rollback on error
		return insertPostCategoriesErr
	}

	insertPostFilesErr := InsertPostFiles(post.ID, uploadedFilesAddresses, user_id, tx)
	if insertPostFilesErr != nil {
		tx.Rollback() // Rollback on error
		return insertPostFilesErr
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback() // Rollback on error
		return err
	}

	return nil
}

func UpdateStatusPost(post_id int, status string, user_id int) error {
	db := db.OpenDBConnection()
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
		tx.Rollback() // Rollback on error
		return updateStatusPostCategories
	}

	UpdateStatusPostFiles := UpdateStatusPostFiles(post_id, user_id, status, tx)
	if UpdateStatusPostFiles != nil {
		tx.Rollback() // Rollback on error
		return UpdateStatusPostFiles
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback() // Rollback on error
		return err
	}

	return nil
}

func ReadAllPosts() ([]Post, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Query the records
	rows, selectError := db.Query(`
        SELECT p.id as post_id, p.uuid as post_uuid, p.title as post_title, p.description as post_description, p.status as post_status, p.created_at as post_created_at, p.updated_at as post_updated_at, p.updated_by as post_updated_by,
			u.id as user_id, u.name as user_name, u.username as user_username, u.email as user_email,
			c.id as category_id, c.name as category_name,
			IFNULL(pf.id, 0) as post_file_id, pf.file_address
		FROM posts p
			INNER JOIN users u
				ON p.user_id = u.id
			LEFT JOIN post_files pf
				ON p.id = pf.post_id
				AND pf.status = 'enable'
			LEFT JOIN post_categories pc
				ON p.id = pc.post_id
				AND pc.status = 'enable'
			LEFT JOIN categories c
				ON pc.category_id = c.id
				AND c.status = 'enable'
		WHERE p.status != 'delete'
			AND u.status != 'delete'
		ORDER BY p.id desc;
    `)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var posts []Post
	// Map to track posts by their ID to avoid duplicates
	postMap := make(map[int]*Post)

	for rows.Next() {
		var post Post
		var user userManagementModels.User
		var category Category
		var postFile PostFile

		// Scan the post and user data
		err := rows.Scan(
			&post.ID, &post.UUID, &post.Title, &post.Description, &post.Status,
			&post.CreatedAt, &post.UpdatedAt, &post.UpdatedBy, &post.UserId,
			&user.Name, &user.Username, &user.Email,
			&category.ID, &category.Name,
			&postFile.ID, &postFile.FileAddress,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Check if the post already exists in the map
		existingPost, found := postMap[post.ID]
		if !found {
			post.User = user
			post.Categories = []Category{}
			post.PostFiles = []PostFile{}
			postMap[post.ID] = &post
			existingPost = &post
		}

		// Ensure unique categories
		isCategoryAdded := false
		for _, c := range existingPost.Categories {
			if c.ID == category.ID {
				isCategoryAdded = true
				break
			}
		}
		if !isCategoryAdded && category.ID != 0 {
			existingPost.Categories = append(existingPost.Categories, category)
		}

		// Ensure unique post files
		isFileAdded := false
		for _, f := range existingPost.PostFiles {
			if f.ID == postFile.ID {
				isFileAdded = true
				break
			}
		}
		if !isFileAdded && postFile.ID != 0 {
			existingPost.PostFiles = append(existingPost.PostFiles, postFile)
		}
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	// Convert the map of posts into a slice
	for _, post := range postMap {
		posts = append(posts, *post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ID > posts[j].ID
	})

	return posts, nil
}

func ReadPostsByCategoryId(category_id int) ([]Post, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Query the records
	rows, selectError := db.Query(`
        SELECT p.id as post_id, p.uuid as post_uuid, p.title as post_title, p.description as post_description, p.status as post_status, p.created_at as post_created_at, p.updated_at as post_updated_at, p.updated_by as post_updated_by,
			u.id as user_id, u.name as user_name, u.username as user_username, u.email as user_email,
			c.id as category_id, c.name as category_name,
			IFNULL(pf.id, 0) as post_file_id, pf.file_address
		FROM posts p
			INNER JOIN users u
				ON p.user_id = u.id
			LEFT JOIN post_files pf
				ON p.id = pf.post_id
				AND pf.status = 'enable'
			INNER JOIN post_categories filterd_pc
				ON p.id = filterd_pc.post_id
				AND filterd_pc.status = 'enable'
				AND filterd_pc.category_id = ?
			INNER JOIN post_categories pc
				ON p.id = pc.post_id
				AND pc.status = 'enable'
			INNER JOIN categories c
				ON pc.category_id = c.id
				AND c.status = 'enable'
		WHERE p.status != 'delete'
			AND u.status != 'delete'
		ORDER BY p.id desc;
    `, category_id)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var posts []Post
	// Map to track posts by their ID to avoid duplicates
	postMap := make(map[int]*Post)

	for rows.Next() {
		var post Post
		var user userManagementModels.User
		var category Category
		var postFile PostFile

		// Scan the post and user data
		err := rows.Scan(
			&post.ID, &post.UUID, &post.Title, &post.Description, &post.Status,
			&post.CreatedAt, &post.UpdatedAt, &post.UpdatedBy, &post.UserId,
			&user.Name, &user.Username, &user.Email,
			&category.ID, &category.Name,
			&postFile.ID, &postFile.FileAddress,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Check if the post already exists in the map
		existingPost, found := postMap[post.ID]
		if !found {
			post.User = user
			post.Categories = []Category{}
			post.PostFiles = []PostFile{}
			postMap[post.ID] = &post
			existingPost = &post
		}

		// Ensure unique categories
		isCategoryAdded := false
		for _, c := range existingPost.Categories {
			if c.ID == category.ID {
				isCategoryAdded = true
				break
			}
		}
		if !isCategoryAdded && category.ID != 0 {
			existingPost.Categories = append(existingPost.Categories, category)
		}

		// Ensure unique post files
		isFileAdded := false
		for _, f := range existingPost.PostFiles {
			if f.ID == postFile.ID {
				isFileAdded = true
				break
			}
		}
		if !isFileAdded && postFile.ID != 0 {
			existingPost.PostFiles = append(existingPost.PostFiles, postFile)
		}
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	// Convert the map of posts into a slice
	for _, post := range postMap {
		posts = append(posts, *post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ID > posts[j].ID
	})

	return posts, nil
}

func FilterPosts(searchTerm string) ([]Post, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	searchPattern := "%" + searchTerm + "%" // Add wildcards for LIKE comparison

	// Query the records
	rows, selectError := db.Query(`
        SELECT p.id as post_id, p.uuid as post_uuid, p.title as post_title, p.description as post_description, p.status as post_status, p.created_at as post_created_at, p.updated_at as post_updated_at, p.updated_by as post_updated_by,
			u.id as user_id, u.name as user_name, u.username as user_username, u.email as user_email,
			c.id as category_id, c.name as category_name,
			IFNULL(pf.id, 0) as post_file_id, pf.file_address
		FROM posts p
			INNER JOIN users u
				ON p.user_id = u.id
			LEFT JOIN post_files pf
				ON p.id = pf.post_id
				AND pf.status = 'enable'
			LEFT JOIN post_categories pc
				ON p.id = pc.post_id
				AND pc.status = 'enable'
			LEFT JOIN categories c
				ON pc.category_id = c.id
				AND c.status = 'enable'
		WHERE p.status != 'delete'
			AND u.status != 'delete'
      		AND (p.title LIKE ? OR p.description LIKE ?)
		ORDER BY p.id desc;
    `, searchPattern, searchPattern)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var posts []Post
	// Map to track posts by their ID to avoid duplicates
	postMap := make(map[int]*Post)

	for rows.Next() {
		var post Post
		var user userManagementModels.User
		var category Category
		var postFile PostFile

		// Scan the post and user data
		err := rows.Scan(
			&post.ID, &post.UUID, &post.Title, &post.Description, &post.Status,
			&post.CreatedAt, &post.UpdatedAt, &post.UpdatedBy, &post.UserId,
			&user.Name, &user.Username, &user.Email,
			&category.ID, &category.Name,
			&postFile.ID, &postFile.FileAddress,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Check if the post already exists in the map
		existingPost, found := postMap[post.ID]
		if !found {
			post.User = user
			post.Categories = []Category{}
			post.PostFiles = []PostFile{}
			postMap[post.ID] = &post
			existingPost = &post
		}

		// Ensure unique categories
		isCategoryAdded := false
		for _, c := range existingPost.Categories {
			if c.ID == category.ID {
				isCategoryAdded = true
				break
			}
		}
		if !isCategoryAdded && category.ID != 0 {
			existingPost.Categories = append(existingPost.Categories, category)
		}

		// Ensure unique post files
		isFileAdded := false
		for _, f := range existingPost.PostFiles {
			if f.ID == postFile.ID {
				isFileAdded = true
				break
			}
		}
		if !isFileAdded && postFile.ID != 0 {
			existingPost.PostFiles = append(existingPost.PostFiles, postFile)
		}
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	// Convert the map of posts into a slice
	for _, post := range postMap {
		posts = append(posts, *post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ID > posts[j].ID
	})

	return posts, nil
}

func ReadPostsByUserId(userId int) ([]Post, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Query the records
	rows, selectError := db.Query(`
        SELECT p.id as post_id, p.uuid as post_uuid, p.title as post_title, p.description as post_description, p.status as post_status, p.created_at as post_created_at, p.updated_at as post_updated_at, p.updated_by as post_updated_by,
			u.id as user_id, u.name as user_name, u.username as user_username, u.email as user_email,
			c.id as category_id, c.name as category_name,
			IFNULL(pf.id, 0) as post_file_id, pf.file_address
		FROM posts p
			INNER JOIN users u
				ON p.user_id = u.id
				AND u.id = ?
			LEFT JOIN post_files pf
				ON p.id = pf.post_id
				AND pf.status = 'enable'
			LEFT JOIN post_categories pc
				ON p.id = pc.post_id
				AND pc.status = 'enable'
			LEFT JOIN categories c
				ON pc.category_id = c.id
				AND c.status = 'enable'
		WHERE p.status != 'delete'
			AND u.status != 'delete'
		ORDER BY p.id desc;
    `, userId)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var posts []Post
	// Map to track posts by their ID to avoid duplicates
	postMap := make(map[int]*Post)

	for rows.Next() {
		var post Post
		var user userManagementModels.User
		var category Category
		var postFile PostFile

		// Scan the post and user data
		err := rows.Scan(
			&post.ID, &post.UUID, &post.Title, &post.Description, &post.Status,
			&post.CreatedAt, &post.UpdatedAt, &post.UpdatedBy, &post.UserId,
			&user.Name, &user.Username, &user.Email,
			&category.ID, &category.Name,
			&postFile.ID, &postFile.FileAddress,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Check if the post already exists in the map
		existingPost, found := postMap[post.ID]
		if !found {
			post.User = user
			post.Categories = []Category{}
			post.PostFiles = []PostFile{}
			postMap[post.ID] = &post
			existingPost = &post
		}

		// Ensure unique categories
		isCategoryAdded := false
		for _, c := range existingPost.Categories {
			if c.ID == category.ID {
				isCategoryAdded = true
				break
			}
		}
		if !isCategoryAdded && category.ID != 0 {
			existingPost.Categories = append(existingPost.Categories, category)
		}

		// Ensure unique post files
		isFileAdded := false
		for _, f := range existingPost.PostFiles {
			if f.ID == postFile.ID {
				isFileAdded = true
				break
			}
		}
		if !isFileAdded && postFile.ID != 0 {
			existingPost.PostFiles = append(existingPost.PostFiles, postFile)
		}
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	// Convert the map of posts into a slice
	for _, post := range postMap {
		posts = append(posts, *post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ID > posts[j].ID
	})

	return posts, nil
}

func ReadPostsLikedByUserId(userId int) ([]Post, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Query the records
	rows, selectError := db.Query(`
        SELECT p.id as post_id, p.uuid as post_uuid, p.title as post_title, p.description as post_description, p.status as post_status, p.created_at as post_created_at, p.updated_at as post_updated_at, p.updated_by as post_updated_by,
			u.id as user_id, u.name as user_name, u.username as user_username, u.email as user_email,
			c.id as category_id, c.name as category_name,
			IFNULL(pf.id, 0) as post_file_id, pf.file_address
		FROM posts p
			INNER JOIN post_likes pl
				ON pl.post_id = p.id
				AND pl.status = 'enable'
			INNER JOIN users u
				ON p.user_id = u.id
			LEFT JOIN post_files pf
				ON p.id = pf.post_id
				AND pf.status = 'enable'
			INNER JOIN users liked_user
				ON pl.user_id = liked_user.id
				AND liked_user.id = ?
			LEFT JOIN post_categories pc
				ON p.id = pc.post_id
				AND pc.status = 'enable'
			LEFT JOIN categories c
				ON pc.category_id = c.id
				AND c.status = 'enable'
		WHERE p.status != 'delete'
			AND u.status != 'delete'
		ORDER BY p.id desc;
    `, userId)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var posts []Post
	// Map to track posts by their ID to avoid duplicates
	postMap := make(map[int]*Post)

	for rows.Next() {
		var post Post
		var user userManagementModels.User
		var category Category
		var postFile PostFile

		// Scan the post and user data
		err := rows.Scan(
			&post.ID, &post.UUID, &post.Title, &post.Description, &post.Status,
			&post.CreatedAt, &post.UpdatedAt, &post.UpdatedBy, &post.UserId,
			&user.Name, &user.Username, &user.Email,
			&category.ID, &category.Name,
			&postFile.ID, &postFile.FileAddress,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Check if the post already exists in the map
		existingPost, found := postMap[post.ID]
		if !found {
			post.User = user
			post.Categories = []Category{}
			post.PostFiles = []PostFile{}
			postMap[post.ID] = &post
			existingPost = &post
		}

		// Ensure unique categories
		isCategoryAdded := false
		for _, c := range existingPost.Categories {
			if c.ID == category.ID {
				isCategoryAdded = true
				break
			}
		}
		if !isCategoryAdded && category.ID != 0 {
			existingPost.Categories = append(existingPost.Categories, category)
		}

		// Ensure unique post files
		isFileAdded := false
		for _, f := range existingPost.PostFiles {
			if f.ID == postFile.ID {
				isFileAdded = true
				break
			}
		}
		if !isFileAdded && postFile.ID != 0 {
			existingPost.PostFiles = append(existingPost.PostFiles, postFile)
		}
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	// Convert the map of posts into a slice
	for _, post := range postMap {
		posts = append(posts, *post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ID > posts[j].ID
	})

	return posts, nil
}

func ReadPostById(postId int, checkLikeForUser int) (Post, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Query the records
	rows, selectError := db.Query(`
        SELECT p.id as post_id, p.uuid as post_uuid, p.title as post_title, p.description as post_description, p.status as post_status, p.created_at as post_created_at, p.updated_at as post_updated_at, p.updated_by as post_updated_by,
			(SELECT COUNT(DISTINCT id) from post_likes WHERE post_id = p.id AND status != 'delete' AND type = 'like') AS number_of_likes,
			(SELECT COUNT(DISTINCT id) from post_likes WHERE post_id = p.id AND status != 'delete' AND type = 'dislike') AS number_of_dislikes,
			u.id as user_id, u.name as user_name, u.username as user_username, u.email as user_email,
			c.id as category_id, c.name as category_name,
			IFNULL(pf.id, 0) as post_file_id, pf.file_address,
			CASE 
                WHEN EXISTS (SELECT 1 FROM post_likes WHERE post_id = p.id AND status != 'delete' AND type = 'like' AND user_id = ?) THEN 1
                ELSE 0
            END AS is_liked_by_user,
            CASE 
                WHEN EXISTS (SELECT 1 FROM post_likes WHERE post_id = p.id AND status != 'delete' AND type = 'dislike' AND user_id = ?) THEN 1
                ELSE 0
            END AS is_disliked_by_user
		FROM posts p
			INNER JOIN users u
				ON p.user_id = u.id
				AND p.id = ?
			LEFT JOIN post_files pf
				ON p.id = pf.post_id
				AND pf.status = 'enable'
			LEFT JOIN post_categories pc
				ON p.id = pc.post_id
				AND pc.status = 'enable'
			LEFT JOIN categories c
				ON pc.category_id = c.id
				AND c.status = 'enable'
		WHERE p.status != 'delete'
			AND u.status != 'delete';
    `, checkLikeForUser, checkLikeForUser, postId)
	if selectError != nil {
		return Post{}, selectError
	}
	defer rows.Close()

	var post Post
	var user userManagementModels.User

	// Scan the records
	for rows.Next() {
		var category Category
		var postFile PostFile

		err := rows.Scan(
			&post.ID, &post.UUID, &post.Title, &post.Description, &post.Status,
			&post.CreatedAt, &post.UpdatedAt, &post.UpdatedBy,
			&post.NumberOfLikes, &post.NumberOfDislikes,
			&post.UserId, &user.Name, &user.Username, &user.Email,
			&category.ID, &category.Name,
			&postFile.ID, &postFile.FileAddress,
			&post.IsLikedByUser, &post.IsDislikedByUser,
		)
		if err != nil {
			return Post{}, fmt.Errorf("error scanning row: %v", err)
		}

		// Assign user to post
		if post.UserId == 0 { // If this is the first time we're encountering the post
			post.User = user
		}

		// Ensure unique categories
		isCategoryAdded := false
		for _, c := range post.Categories {
			if c.ID == category.ID {
				isCategoryAdded = true
				break
			}
		}
		if !isCategoryAdded && category.ID != 0 {
			post.Categories = append(post.Categories, category)
		}

		// Ensure unique post files
		isFileAdded := false
		for _, f := range post.PostFiles {
			if f.ID == postFile.ID {
				isFileAdded = true
				break
			}
		}
		if !isFileAdded && postFile.ID != 0 {
			post.PostFiles = append(post.PostFiles, postFile)
		}
	}

	// If no rows were returned, the post doesn't exist
	if post.ID == 0 {
		return Post{}, fmt.Errorf("post with ID %d not found", postId)
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return Post{}, fmt.Errorf("row iteration error: %v", err)
	}

	return post, nil
}

func ReadPostByUUID(postUUID string, checkLikeForUser int) (Post, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Query the records
	rows, selectError := db.Query(`
        SELECT p.id as post_id, p.uuid as post_uuid, p.title as post_title, p.description as post_description, p.status as post_status, p.created_at as post_created_at, p.updated_at as post_updated_at, p.updated_by as post_updated_by,
			(SELECT COUNT(DISTINCT id) from post_likes WHERE post_id = p.id AND status != 'delete' AND type = 'like') AS number_of_likes,
			(SELECT COUNT(DISTINCT id) from post_likes WHERE post_id = p.id AND status != 'delete' AND type = 'dislike') AS number_of_dislikes,
			u.id as user_id, u.name as user_name, u.username as user_username, u.email as user_email,
			c.id as category_id, c.name as category_name,
			IFNULL(pf.id, 0) as post_file_id, pf.file_address,
			CASE 
                WHEN EXISTS (SELECT 1 FROM post_likes WHERE post_id = p.id AND status != 'delete' AND type = 'like' AND user_id = ?) THEN 1
                ELSE 0
            END AS is_liked_by_user,
            CASE 
                WHEN EXISTS (SELECT 1 FROM post_likes WHERE post_id = p.id AND status != 'delete' AND type = 'dislike' AND user_id = ?) THEN 1
                ELSE 0
            END AS is_disliked_by_user
		FROM posts p
			INNER JOIN users u
				ON p.user_id = u.id
				AND p.uuid = ?
			LEFT JOIN post_files pf
				ON p.id = pf.post_id
				AND pf.status = 'enable'
			LEFT JOIN post_categories pc
				ON p.id = pc.post_id
				AND pc.status = 'enable'
			LEFT JOIN categories c
				ON pc.category_id = c.id
				AND c.status = 'enable'
		WHERE p.status != 'delete'
			AND u.status != 'delete';
    `, checkLikeForUser, checkLikeForUser, postUUID)
	if selectError != nil {
		return Post{}, selectError
	}
	defer rows.Close()

	var post Post
	var user userManagementModels.User

	// Scan the records
	for rows.Next() {
		var category Category
		var postFile PostFile

		err := rows.Scan(
			&post.ID, &post.UUID, &post.Title, &post.Description, &post.Status,
			&post.CreatedAt, &post.UpdatedAt, &post.UpdatedBy,
			&post.NumberOfLikes, &post.NumberOfDislikes,
			&post.UserId, &user.Name, &user.Username, &user.Email,
			&category.ID, &category.Name,
			&postFile.ID, &postFile.FileAddress,
			&post.IsLikedByUser, &post.IsDislikedByUser,
		)
		if err != nil {
			return Post{}, fmt.Errorf("error scanning row: %v", err)
		}

		// Ensure unique categories
		isCategoryAdded := false
		for _, c := range post.Categories {
			if c.ID == category.ID {
				isCategoryAdded = true
				break
			}
		}
		if !isCategoryAdded && category.ID != 0 {
			post.Categories = append(post.Categories, category)
		}

		// Ensure unique post files
		isFileAdded := false
		for _, f := range post.PostFiles {
			if f.ID == postFile.ID {
				isFileAdded = true
				break
			}
		}
		if !isFileAdded && postFile.ID != 0 {
			post.PostFiles = append(post.PostFiles, postFile)
		}
	}

	// If no rows were returned, the post doesn't exist
	if post.ID == 0 {
		return Post{}, fmt.Errorf("post with UUID %s not found", postUUID)
	}

	post.User = user

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return Post{}, fmt.Errorf("row iteration error: %v", err)
	}

	return post, nil
}

func ReadPostByUserID(postId int, userID int) (Post, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes
	// Updated query to join comments with posts
	rows, selectError := db.Query(`
        SELECT p.id as post_id, p.uuid as post_uuid, p.title as post_title, p.description as post_description, p.status as post_status, p.created_at as post_created_at, p.updated_at as post_updated_at, p.updated_by as post_updated_by,
			p.user_id as post_user_id, u.id as user_id, u.name as user_name, u.username as user_username, u.email as user_email,
			c.id as category_id, c.name as category_name,
			IFNULL(pf.id, 0) as post_file_id, pf.file_address,
			COALESCE(pl.type, '')
		FROM posts p
			INNER JOIN users u
				ON p.user_id = u.id
				AND p.id = ?
			LEFT JOIN post_files pf
				ON p.id = pf.post_id
				AND pf.status = 'enable'
			LEFT JOIN post_categories pc
				ON p.id = pc.post_id
				AND pc.status = 'enable'
			LEFT JOIN categories c
				ON pc.category_id = c.id
				AND c.status = 'enable'
			LEFT JOIN post_likes pl
				ON p.id = pl.post_id AND pl.status != 'delete'	
		WHERE p.status != 'delete'
			AND u.status != 'delete';
    `, postId)
	if selectError != nil {
		return Post{}, selectError
	}
	defer rows.Close()

	var post Post
	var user userManagementModels.User

	// Scan the records
	for rows.Next() {
		fmt.Println(1)
		var category Category
		var postFile PostFile
		var Type string
		err := rows.Scan(
			&post.ID, &post.UUID, &post.Title, &post.Description, &post.Status,
			&post.CreatedAt, &post.UpdatedAt, &post.UpdatedBy, &post.UserId,
			&user.ID, &user.Name, &user.Username, &user.Email,
			&category.ID, &category.Name,
			&postFile.ID, &postFile.FileAddress,
			&Type,
		)
		if err != nil {
			return Post{}, fmt.Errorf("error scanning row: %v", err)
		}
		if user.ID == userID {
			if Type == "like" {
				post.IsLikedByUser = true
			} else if Type == "dislike" {
				post.IsDislikedByUser = true
			}
		}
		if Type == "like" {
			post.NumberOfLikes++
		} else if Type == "dislike" {
			post.NumberOfDislikes++
		}

		// Ensure unique categories
		isCategoryAdded := false
		for _, c := range post.Categories {
			if c.ID == category.ID {
				isCategoryAdded = true
				break
			}
		}
		if !isCategoryAdded && category.ID != 0 {
			post.Categories = append(post.Categories, category)
		}

		// Ensure unique post files
		isFileAdded := false
		for _, f := range post.PostFiles {
			if f.ID == postFile.ID {
				isFileAdded = true
				break
			}
		}
		if !isFileAdded && postFile.ID != 0 {
			post.PostFiles = append(post.PostFiles, postFile)
		}
	}

	post.User = user

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return Post{}, fmt.Errorf("row iteration error: %v", err)
	}

	fmt.Println("im here")
	return post, nil
}

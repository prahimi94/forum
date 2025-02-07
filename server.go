package main

import (
	"fmt"
	"forum/db"
	forumManagementControllers "forum/modules/forumManagement/controllers"
	userManagementControllers "forum/modules/userManagement/controllers"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var IsTest = false

func setupRoutes() {
	http.Handle("/css/", http.FileServer(http.Dir("assets/")))
	http.Handle("/js/", http.FileServer(http.Dir("assets/")))
	http.Handle("/img/", http.FileServer(http.Dir("assets/")))
	http.Handle("/uploads/", http.FileServer(http.Dir("static/")))

	// Register route handlers
	http.HandleFunc("/", forumManagementControllers.MainPageHandler)
	// http.HandleFunc("/home/", forumManagementControllers.HomePageHandler)
	http.HandleFunc("/auth/", userManagementControllers.AuthHandler)
	http.HandleFunc("/logout/", userManagementControllers.Logout)
	http.HandleFunc("/register", userManagementControllers.RegisterHandler) /*post method*/
	http.HandleFunc("/login", userManagementControllers.LoginHandler)       /*post method*/

	http.HandleFunc("/newPost/", forumManagementControllers.CreatePost)
	http.HandleFunc("/submitPost", forumManagementControllers.SubmitPost) /*post method*/
	http.HandleFunc("/myCreatedPosts/", forumManagementControllers.ReadMyCreatedPosts)
	http.HandleFunc("/myLikedPosts/", forumManagementControllers.ReadMyLikedPosts)
	http.HandleFunc("/post/", forumManagementControllers.ReadPost)
	http.HandleFunc("/posts/", forumManagementControllers.ReadPostsByCategory)
	http.HandleFunc("/filterPosts/", forumManagementControllers.FilterPosts)
	http.HandleFunc("/likePost", forumManagementControllers.LikePost)
	http.HandleFunc("/editPost/", forumManagementControllers.EditPost)
	http.HandleFunc("/updatePost", forumManagementControllers.UpdatePost) /*post method*/
	http.HandleFunc("/deletePost", forumManagementControllers.DeletePost) /*post method*/

	http.HandleFunc("/likeComment", forumManagementControllers.LikeComment)
	http.HandleFunc("/submitComment", forumManagementControllers.SubmitComment) /*post method*/
	http.HandleFunc("/updateComment", forumManagementControllers.UpdateComment) /*post method*/
	http.HandleFunc("/deleteComment", forumManagementControllers.DeleteComment) /*post method*/
}

func testInit() (bool, string, int) {
	if IsTest {
		db := db.OpenDBConnection()
		testSessionToken := "92fdfae0-8e51-4f49-915a-c77add3e101f"

		testUserId := -1
		db.QueryRow(`SELECT id
					FROM users
					WHERE uuid = ?`, testSessionToken).Scan(&testUserId)
		if testUserId == -1 {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
			if err != nil {
				log.Fatalf("Error: %v", err)
				os.Exit(1)
			}

			insertTestUserQuery := `INSERT INTO users(uuid, type, name, username, password, email)
									VALUES (?, 'test_user', 'test', 'test', ?, 'test@test.com');`
			insertTestUserResult, insertTestUserStatusErr := db.Exec(insertTestUserQuery, testSessionToken, hashedPassword)
			if insertTestUserStatusErr != nil {
				log.Fatalf("Error: %v", insertTestUserStatusErr)
				os.Exit(1)
			}

			lastInsertID, err := insertTestUserResult.LastInsertId()
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			testUserId = int(lastInsertID)
		}

		sessionId := -1
		var expirationTime time.Time
		db.QueryRow(`SELECT id, expires_at
					FROM sessions
					WHERE session_token = ?`, testSessionToken).Scan(&sessionId, &expirationTime)
		fmt.Println(sessionId)
		fmt.Println(expirationTime)
		if sessionId == -1 { //if sessionId == -1 means that this session didnt found
			sessionExpiresAt := time.Now().Add(1 * time.Hour)

			insertSessionQuery := `INSERT INTO sessions(session_token, user_id, expires_at)
							VALUES (?, ?, ?);`
			_, insertSessionStatusErr := db.Exec(insertSessionQuery, testSessionToken, testUserId, sessionExpiresAt)
			if insertSessionStatusErr != nil {
				log.Fatalf("Error: %v", insertSessionStatusErr)
				os.Exit(1)
			}
		} else if time.Now().After(expirationTime) {
			sessionExpiresAt := time.Now().Add(1 * time.Hour)

			updateSessionQuery := `UPDATE sessions SET expires_at = ? WHERE session_token = ?;`
			_, updateSessionStatusErr := db.Exec(updateSessionQuery, sessionExpiresAt, testSessionToken)
			if updateSessionStatusErr != nil {
				log.Fatalf("Error: %v", updateSessionStatusErr)
				os.Exit(1)
			}
		}

		db.Close()
		return true, testSessionToken, testUserId
	}
	return false, "", -1
}

func finishTest() {
	if IsTest {
		db := db.OpenDBConnection()
		testSessionToken := "92fdfae0-8e51-4f49-915a-c77add3e101f"

		testUserId := -1
		db.QueryRow(`SELECT id
					FROM users
					WHERE uuid = ?
					AND username = 'test'
					AND email = 'test@test.com'`, testSessionToken).Scan(&testUserId)
		test2UserId := -1
		db.QueryRow(`SELECT id
					FROM users
					WHERE username = 'test2'
					AND email = 'test2@test2.com'`).Scan(&test2UserId)

		deletePostCategoryQuery := `DELETE FROM post_categories
									WHERE post_id IN (
									SELECT id FROM posts WHERE user_id IN (?, ?)
									);`
		_, deletePostCategoryStatusErr := db.Exec(deletePostCategoryQuery, testUserId, test2UserId)
		if deletePostCategoryStatusErr != nil {
			log.Fatalf("Error: %v", deletePostCategoryStatusErr)
			os.Exit(1)
		}

		deletePostQuery := `DELETE FROM posts WHERE user_id IN (?, ?);`
		_, deletePostStatusErr := db.Exec(deletePostQuery, testUserId, test2UserId)
		if deletePostStatusErr != nil {
			log.Fatalf("Error: %v", deletePostStatusErr)
			os.Exit(1)
		}

		deleteSessionQuery := `DELETE FROM sessions WHERE user_id IN (?, ?);`
		_, deleteSessionStatusErr := db.Exec(deleteSessionQuery, testUserId, test2UserId)
		if deleteSessionStatusErr != nil {
			log.Fatalf("Error: %v", deleteSessionStatusErr)
			os.Exit(1)
		}

		deleteUserQuery := `DELETE FROM users WHERE id IN (?, ?);`
		_, deleteUserStatusErr := db.Exec(deleteUserQuery, testUserId, test2UserId)
		if deleteUserStatusErr != nil {
			log.Fatalf("Error: %v", deleteUserStatusErr)
			os.Exit(1)
		}

		db.Close()
	}
}

// main initializes the HTTP server, registers routes, and starts listening for incoming requests.
func main() {
	// remove fresh version of database
	// if err := db.ExecuteSQLFile("db/forum.sql"); err != nil {
	// 	log.Fatalf("Error: %v", err)
	// 	os.Exit(1)
	// }

	setupRoutes()

	//start the server on port 8080
	log.Println("Starting server on: http://localhost:8080")
	log.Println("Status ok: ", http.StatusOK)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

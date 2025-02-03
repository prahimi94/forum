package main

import (
	forumManagementControllers "forum/modules/forumManagement/controllers"
	userManagementControllers "forum/modules/userManagement/controllers"
	"log"
	"net/http"
)

// main initializes the HTTP server, registers routes, and starts listening for incoming requests.
func main() {
	http.Handle("/css/", http.FileServer(http.Dir("assets/")))
	http.Handle("/js/", http.FileServer(http.Dir("assets/")))
	http.Handle("/img/", http.FileServer(http.Dir("assets/")))

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

	//start the server on port 8080
	log.Println("Starting server on: http://localhost:8080")
	log.Println("Status ok: ", http.StatusOK)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

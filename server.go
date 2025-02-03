package main

import (
	forumManagementControllers "forum/forumManagement/controllers"
	userManagementControllers "forum/userManagement/controllers"
	"log"
	"net/http"
)

// main initializes the HTTP server, registers routes, and starts listening for incoming requests.
func main() {
	http.Handle("/static/", http.FileServer(http.Dir("frontend/public/")))
	http.Handle("/img/", http.FileServer(http.Dir("frontend/public/")))

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
	http.HandleFunc("/likeComment", forumManagementControllers.LikeComment)
	http.HandleFunc("/likePost", forumManagementControllers.LikePost)

	http.HandleFunc("/submitComment", forumManagementControllers.SubmitComment) /*post method*/
	http.HandleFunc("/editPost/", forumManagementControllers.EditPost)
	http.HandleFunc("/updatePost", forumManagementControllers.UpdatePost) /*post method*/
	http.HandleFunc("/deletePost", forumManagementControllers.DeletePost) /*post method*/

	http.HandleFunc("/updateComment", forumManagementControllers.UpdateComment) /*post method*/
	http.HandleFunc("/deleteComment", forumManagementControllers.DeleteComment) /*post method*/

	// forumManagement.InsertPost("first post", "first post description", []int{1, 2}, 1)
	// forumManagement.UpdatePost(1, "first post", "first post description updated", []int{2, 3}, 1)
	// forumManagement.UpdateStatusPost(1, "enable", 1)
	// post, err := forumManagementModels.ReadPostById(2)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// } else {
	// 	fmt.Println("Post:", post)
	// 	if len(post.Categories) == 0 {
	// 		fmt.Println("No categories associated with this post.")
	// 	}
	// }

	// categories, err := forumManagementModels.ReadAllCategories()
	// fmt.Println(categories)
	// fmt.Println(err)

	//start the server on port 8080
	log.Println("Starting server on: http://localhost:8080")
	log.Println("Status ok: ", http.StatusOK)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

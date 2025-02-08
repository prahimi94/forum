package routes

import (
	"net/http"

	"forum/middlewares"
	forumManagementControllers "forum/modules/forumManagement/controllers"
	userManagementControllers "forum/modules/userManagement/controllers"
)

func SetupRoutes() {
	http.Handle("/css/", http.FileServer(http.Dir("assets/")))
	http.Handle("/js/", http.FileServer(http.Dir("assets/")))
	http.Handle("/img/", http.FileServer(http.Dir("assets/")))
	http.Handle("/uploads/", http.FileServer(http.Dir("static/")))

	// Public routes
	http.HandleFunc("/", forumManagementControllers.MainPageHandler)
	http.HandleFunc("/auth/", userManagementControllers.AuthHandler)
	http.HandleFunc("/logout/", userManagementControllers.Logout)
	http.HandleFunc("/register", userManagementControllers.RegisterHandler) /*post method*/
	http.HandleFunc("/login", userManagementControllers.LoginHandler)       /*post method*/
	http.HandleFunc("/post/", forumManagementControllers.ReadPost)
	http.HandleFunc("/posts/", forumManagementControllers.ReadPostsByCategory)
	http.HandleFunc("/filterPosts/", forumManagementControllers.FilterPosts)

	// Protected routes (require authentication)
	http.HandleFunc("/newPost/", middlewares.AuthMiddleware(forumManagementControllers.CreatePost))
	http.HandleFunc("/submitPost", middlewares.AuthMiddleware(forumManagementControllers.SubmitPost)) /*post method*/
	http.HandleFunc("/editPost/", middlewares.AuthMiddleware(forumManagementControllers.EditPost))
	http.HandleFunc("/updatePost", middlewares.AuthMiddleware(forumManagementControllers.UpdatePost)) /*post method*/
	http.HandleFunc("/deletePost", middlewares.AuthMiddleware(forumManagementControllers.DeletePost)) /*post method*/
	http.HandleFunc("/myCreatedPosts/", middlewares.AuthMiddleware(forumManagementControllers.ReadMyCreatedPosts))
	http.HandleFunc("/myLikedPosts/", middlewares.AuthMiddleware(forumManagementControllers.ReadMyLikedPosts))
	http.HandleFunc("/likePost", middlewares.AuthMiddleware(forumManagementControllers.LikePost))

	http.HandleFunc("/likeComment", middlewares.AuthMiddleware(forumManagementControllers.LikeComment))
	http.HandleFunc("/submitComment", middlewares.AuthMiddleware(forumManagementControllers.SubmitComment)) /*post method*/
	http.HandleFunc("/updateComment", middlewares.AuthMiddleware(forumManagementControllers.UpdateComment)) /*post method*/
	http.HandleFunc("/deleteComment", middlewares.AuthMiddleware(forumManagementControllers.DeleteComment)) /*post method*/
}

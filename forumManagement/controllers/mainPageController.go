package controller

import (
	"fmt"
	"forum/errorManagement"
	"forum/forumManagement/models"
	"net/http"
	"text/template"

	userManagementControllers "forum/userManagement/controllers"

	_ "github.com/mattn/go-sqlite3"
)

const publicUrl = "forumManagement/views/"

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here")
	if r.Method != http.MethodGet {
		errorManagement.HandleErrorPage(w, r, errorManagement.MethodNotAllowedError)
		return
	}

	if r.URL.Path != "/" {
		// If the URL is not exactly "/", respond with 404
		errorManagement.HandleErrorPage(w, r, errorManagement.NotFoundError)
		return
	}

	loginStatus, userId, checkLoginError := userManagementControllers.CheckLogin(r)
	if checkLoginError != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}
	if loginStatus {
		fmt.Println("logged in userid is: ", userId)
		// return
	} else {
		fmt.Println("user is not logged in")
	}

	categories, err := models.ReadAllCategories()
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}

	posts, err := models.ReadAllPosts()
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}

	data_obj_sender := struct {
		Posts      []models.Post
		Categories []models.Category
	}{
		Posts:      posts,
		Categories: categories,
	}

	tmpl, err := template.ParseFiles(
		publicUrl + "index.html",
	)
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}

	err = tmpl.Execute(w, data_obj_sender)
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}
}

package controller

import (
	"fmt"
	errorManagementControllers "forum/errorManagement/controllers"
	"forum/forumManagement/models"
	"net/http"
	"text/template"

	userManagementControllers "forum/userManagement/controllers"

	_ "github.com/mattn/go-sqlite3"
)

const publicUrl = "forumManagement/views/"

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	if r.URL.Path != "/" {
		// If the URL is not exactly "/", respond with 404
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.NotFoundError)
		return
	}

	loginStatus, userId, checkLoginError := userManagementControllers.CheckLogin(r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		fmt.Println("logged in userid is: ", userId)
		userManagementControllers.RedirectToHome(w, r)
		return
	} else {
		fmt.Println("user is not logged in")
	}

	categories, err := models.ReadAllCategories()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	posts, err := models.ReadAllPosts()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
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
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	err = tmpl.Execute(w, data_obj_sender)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	if r.URL.Path != "/home/" {
		// If the URL is not exactly "/", respond with 404
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.NotFoundError)
		return
	}

	loginStatus, userId, checkLoginError := userManagementControllers.CheckLogin(r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		fmt.Println("logged in userid is: ", userId)
	} else {
		fmt.Println("user is not logged in")
		userManagementControllers.RedirectToIndex(w, r)
		return
	}

	categories, err := models.ReadAllCategories()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	posts, err := models.ReadAllPosts()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
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
		publicUrl + "home.html",
	)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	err = tmpl.Execute(w, data_obj_sender)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
}

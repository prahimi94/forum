package controller

import (
	errorManagementControllers "forum/modules/errorManagement/controllers"
	"forum/modules/forumManagement/models"
	"forum/utils"
	"net/http"
	"text/template"

	userManagementControllers "forum/modules/userManagement/controllers"
	userManagementModels "forum/modules/userManagement/models"

	_ "github.com/mattn/go-sqlite3"
)

const publicUrl = "modules/forumManagement/views/"

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
		LoginUser  userManagementModels.User
		Posts      []models.Post
		Categories []models.Category
	}{
		LoginUser:  userManagementModels.User{},
		Posts:      posts,
		Categories: categories,
	}

	loginStatus, loginUser, _, checkLoginError := userManagementControllers.CheckLogin(w, r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		data_obj_sender.LoginUser = loginUser
	}

	// Create a template with a function map
	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"formatDate": utils.FormatDate, // Register function globally
	}).ParseFiles(
		publicUrl+"index.html",
		publicUrl+"templates/header.html",
		publicUrl+"templates/navbar.html",
		publicUrl+"templates/hero.html",
		publicUrl+"templates/posts.html",
		publicUrl+"templates/footer.html",
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

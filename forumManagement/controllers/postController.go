package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	errorManagementControllers "forum/errorManagement/controllers"
	"forum/forumManagement/models"
	"forum/utils"
	"log"
	"net/http"
	"text/template"

	userManagementControllers "forum/userManagement/controllers"

	_ "github.com/mattn/go-sqlite3"
)

func ReadAllPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, userId, checkLoginError := userManagementControllers.CheckLogin(r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		fmt.Println("logged in userid is: ", userId)
		// return
	} else {
		fmt.Println("user is not logged in")
	}

	posts, err := models.ReadAllPosts()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(
		publicUrl + "posts.html",
	)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	err = tmpl.Execute(w, posts)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
}

func ReadPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, userId, checkLoginError := userManagementControllers.CheckLogin(r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		fmt.Println("logged in userid is: ", userId)
		// return
	} else {
		fmt.Println("user is not logged in")
	}

	uuid, errUrl := utils.ExtractUUIDFromUrl(r.URL.Path, "post")
	if errUrl == "not found" {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.NotFoundError)
		return
	}

	post, err := models.ReadPostByUUID(uuid)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(
		publicUrl + "post_details.html",
	)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	err = tmpl.Execute(w, post)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, userId, checkLoginError := userManagementControllers.CheckLogin(r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		fmt.Println("logged in userid is: ", userId)
		// return
	} else {
		fmt.Println("user is not logged in")
	}

	tmpl, err := template.ParseFiles(
		publicUrl + "new_post.html",
	)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
}

func SubmitPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, userId, checkLoginError := userManagementControllers.CheckLogin(r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		fmt.Println("logged in userid is: ", userId)
		// return
	} else {
		fmt.Println("user is not logged in")
	}

	err := r.ParseForm()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}
	title := r.FormValue("title")
	description := r.FormValue("description")
	categories := r.FormValue("categories")
	if len(title) == 0 || len(description) == 0 || len(categories) == 0 {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	// Declare a slice to store the result
	var categoryIds []int
	// Unmarshal the string into the slice
	err = json.Unmarshal([]byte(categories), &categoryIds)
	if err != nil {
		log.Fatal(err)
	}

	post := &models.Post{
		Title:       title,
		Description: description,
	}

	// Insert a record while checking duplicates
	_, insertError := models.InsertPost(post, categoryIds)
	if insertError != nil {
		if errors.Is(insertError, sql.ErrNoRows) {
			// todo show toast
			fmt.Println("Post already exists!")
		} else {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		}
		return
	} else {
		fmt.Println("Post added successfully!")
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

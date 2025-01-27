package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"forum/errorManagement"
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
		errorManagement.HandleErrorPage(w, r, errorManagement.MethodNotAllowedError)
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

	posts, err := models.ReadAllPosts()
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(
		publicUrl + "posts.html",
	)
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}

	err = tmpl.Execute(w, posts)
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}
}

func readPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagement.HandleErrorPage(w, r, errorManagement.MethodNotAllowedError)
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

	uuid, errUrl := utils.ExtractUUIDFromUrl(r.URL.Path, "post")
	if errUrl == "not found" {
		errorManagement.HandleErrorPage(w, r, errorManagement.NotFoundError)
		return
	}

	post, err := models.ReadPostByUUID(uuid)
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(
		publicUrl + "post_details.html",
	)
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}

	err = tmpl.Execute(w, post)
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}
}

func createPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagement.HandleErrorPage(w, r, errorManagement.MethodNotAllowedError)
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

	tmpl, err := template.ParseFiles(
		publicUrl + "new_post.html",
	)
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}
}

func submitPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorManagement.HandleErrorPage(w, r, errorManagement.MethodNotAllowedError)
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

	err := r.ParseForm()
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.BadRequestError)
		return
	}
	title := r.FormValue("title")
	description := r.FormValue("description")
	categories := r.FormValue("categories")
	if len(title) == 0 || len(description) == 0 || len(categories) == 0 {
		errorManagement.HandleErrorPage(w, r, errorManagement.BadRequestError)
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
			errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		}
		return
	} else {
		fmt.Println("Post added successfully!")
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

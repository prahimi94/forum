package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/errorManagement"
	"forum/forumManagement/models"
	"net/http"
	"text/template"

	userManagementControllers "forum/userManagement/controllers"

	_ "github.com/mattn/go-sqlite3"
)

func ReadAllComments(w http.ResponseWriter, r *http.Request) {
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

	comments, err := models.ReadAllComments()
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(
		publicUrl + "comments.html",
	)
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}

	err = tmpl.Execute(w, comments)
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}
}

func readPostComments(w http.ResponseWriter, r *http.Request) {
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

	comments, err := models.ReadCommentsByPostId()
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(
		publicUrl + "post_comments.html",
	)
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}

	err = tmpl.Execute(w, comments)
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}
}

func createComment(w http.ResponseWriter, r *http.Request) {
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
		publicUrl + "new_comment.html",
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

func submitComment(w http.ResponseWriter, r *http.Request) {
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
	post_id := r.FormValue("post_id")
	description := r.FormValue("description")
	if len(post_id) == 0 || len(description) == 0 {
		errorManagement.HandleErrorPage(w, r, errorManagement.BadRequestError)
		return
	}

	comment := &models.Comment{
		PostId:      post_id,
		Description: description,
	}

	// Insert a record while checking duplicates
	_, insertError := models.InsertComment(comment)
	if insertError != nil {
		if errors.Is(insertError, sql.ErrNoRows) {
			// todo show toast
			fmt.Println("Comment already exists!")
		} else {
			errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		}
		return
	} else {
		fmt.Println("Comment added successfully!")
	}
}

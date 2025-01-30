package controller

import (
	"fmt"
	errorManagementControllers "forum/errorManagement/controllers"
	"net/http"

	userManagementControllers "forum/userManagement/controllers"

	_ "github.com/mattn/go-sqlite3"
)

func ReadAllComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, loginUser, _, checkLoginError := userManagementControllers.CheckLogin(r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		fmt.Println("logged in userid is: ", loginUser.ID)
		// return
	} else {
		fmt.Println("user is not logged in")
	}

	// comments, err := models.ReadAllComments()
	// if err != nil {
	// 	errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
	// 	return
	// }

	// tmpl, err := template.ParseFiles(
	// 	publicUrl + "comments.html",
	// )
	// if err != nil {
	// 	errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
	// 	return
	// }

	// err = tmpl.Execute(w, comments)
	// if err != nil {
	// 	errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
	// 	return
	// }
}

func readPostComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, loginUser, _, checkLoginError := userManagementControllers.CheckLogin(r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		fmt.Println("logged in userid is: ", loginUser.ID)
		// return
	} else {
		fmt.Println("user is not logged in")
	}

	// comments, err := models.ReadCommentsByPostId()
	// if err != nil {
	// 	errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
	// 	return
	// }

	// tmpl, err := template.ParseFiles(
	// 	publicUrl + "post_comments.html",
	// )
	// if err != nil {
	// 	errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
	// 	return
	// }

	// err = tmpl.Execute(w, comments)
	// if err != nil {
	// 	errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
	// 	return
	// }
}

func createComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, loginUser, _, checkLoginError := userManagementControllers.CheckLogin(r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		fmt.Println("logged in userid is: ", loginUser.ID)
		// return
	} else {
		fmt.Println("user is not logged in")
	}

	// tmpl, err := template.ParseFiles(
	// 	publicUrl + "new_comment.html",
	// )
	// if err != nil {
	// 	errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
	// 	return
	// }

	// err = tmpl.Execute(w, nil)
	// if err != nil {
	// 	errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
	// 	return
	// }
}

func submitComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, loginUser, _, checkLoginError := userManagementControllers.CheckLogin(r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		fmt.Println("logged in userid is: ", loginUser.ID)
		// return
	} else {
		fmt.Println("user is not logged in")
	}

	err := r.ParseForm()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}
	post_id := r.FormValue("post_id")
	description := r.FormValue("description")
	if len(post_id) == 0 || len(description) == 0 {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	// comment := &models.Comment{
	// 	PostId:      post_id,
	// 	Description: description,
	// }

	// // Insert a record while checking duplicates
	// _, insertError := models.InsertComment(comment)
	// if insertError != nil {
	// 	if errors.Is(insertError, sql.ErrNoRows) {
	// 		// todo show toast
	// 		fmt.Println("Comment already exists!")
	// 	} else {
	// 		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
	// 	}
	// 	return
	// } else {
	// 	fmt.Println("Comment added successfully!")
	// }
}

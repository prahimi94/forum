package controller

import (
	"database/sql"
	"errors"
	"fmt"
	errorManagementControllers "forum/errorManagement/controllers"
	"forum/forumManagement/models"
	"forum/utils"
	"net/http"
	"strconv"
	"text/template"

	userManagementControllers "forum/userManagement/controllers"
	userManagementModels "forum/userManagement/models"

	_ "github.com/mattn/go-sqlite3"
)

func ReadAllPosts(w http.ResponseWriter, r *http.Request) {
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

	posts, err := models.ReadAllPosts()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	// Create a template with a function map
	tmpl, err := template.New("posts.html").Funcs(template.FuncMap{
		"formatDate": utils.FormatDate, // Register function globally
	}).ParseFiles(
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

	data_obj_sender := struct {
		LoginUser userManagementModels.User
		Post      models.Post
		Comments  []models.Comment
	}{
		LoginUser: loginUser,
		Post:      post,
		Comments:  nil,
	}

	if loginStatus {
		comments, err := models.ReadAllCommentsOfUserForPost(post.ID, loginUser.ID)
		if err != nil {
			fmt.Println(err)
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		data_obj_sender.Comments = comments
	} else {
		comments, err := models.ReadAllCommentsForPost(post.ID)
		if err != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		data_obj_sender.Comments = comments
	}

	// Create a template with a function map
	tmpl, err := template.New("post_details.html").Funcs(template.FuncMap{
		"formatDate": utils.FormatDate, // Register function globally
	}).ParseFiles(
		publicUrl+"post_details.html",
		publicUrl+"templates/header.html",
		publicUrl+"templates/navbar.html",
		publicUrl+"templates/footer.html",
	)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	// Execute template with data
	err = tmpl.Execute(w, data_obj_sender)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
	}
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
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

	categories, err := models.ReadAllCategories()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		LoginUser  userManagementModels.User
		Categories []models.Category
	}{
		LoginUser:  loginUser,
		Categories: categories,
	}

	tmpl, err := template.ParseFiles(
		publicUrl+"new_post.html",
		publicUrl+"templates/header.html",
		publicUrl+"templates/navbar.html",
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

func SubmitPost(w http.ResponseWriter, r *http.Request) {
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

	title := r.FormValue("title")
	description := r.FormValue("description")
	categories := r.Form["categories"]
	if len(title) == 0 || len(description) == 0 || len(categories) == 0 {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	post := &models.Post{
		Title:       title,
		Description: description,
		UserId:      loginUser.ID,
	}

	// Convert the string slice to an int slice
	categoryIds := make([]int, 0, len(categories))
	for _, category := range categories {
		if id, err := strconv.Atoi(category); err == nil {
			categoryIds = append(categoryIds, id)
		} else {
			// Handle error if conversion fails (for example, invalid input)
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
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

func EditPost(w http.ResponseWriter, r *http.Request) {
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

	categories, err := models.ReadAllCategories()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		LoginUser  userManagementModels.User
		Post       models.Post
		Categories []models.Category
	}{
		LoginUser:  loginUser,
		Post:       post,
		Categories: categories,
	}

	tmpl, err := template.ParseFiles(
		publicUrl+"edit_post.html",
		publicUrl+"templates/header.html",
		publicUrl+"templates/navbar.html",
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

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
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

	title := r.FormValue("title")
	description := r.FormValue("description")
	categories := r.Form["categories"]
	if len(title) == 0 || len(description) == 0 || len(categories) == 0 {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	post := &models.Post{
		Title:       title,
		Description: description,
		UserId:      loginUser.ID,
	}

	// Convert the string slice to an int slice
	categoryIds := make([]int, 0, len(categories))
	for _, category := range categories {
		if id, err := strconv.Atoi(category); err == nil {
			categoryIds = append(categoryIds, id)
		} else {
			// Handle error if conversion fails (for example, invalid input)
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
	}

	// Update a record while checking duplicates
	updateError := models.UpdatePost(post, categoryIds, loginUser.ID)
	if updateError != nil {
		if errors.Is(updateError, sql.ErrNoRows) {
			// todo show toast
			fmt.Println("Post already exists!")
		} else {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		}
		return
	} else {
		fmt.Println("Post updated successfully!")
	}

	userManagementControllers.RedirectToPrevPage(w, r)
}

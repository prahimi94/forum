package userManagement

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/errorManagement"
	"net/http"
	"text/template"
	"time"

	"github.com/gofrs/uuid/v5"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const publicUrl = "frontend/public/"

var u1 = uuid.Must(uuid.NewV4())

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagement.HandleErrorPage(w, r, errorManagement.MethodNotAllowedError)
		return
	}

	loginStatus, userId, checkLoginError := CheckLogin(r)
	if checkLoginError != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}
	if loginStatus {
		fmt.Println("logged in userid is: ", userId)
		return
	}

	tmpl, err := template.ParseFiles(
		publicUrl + "authPage.html",
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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorManagement.HandleErrorPage(w, r, errorManagement.MethodNotAllowedError)
		return
	}

	loginStatus, userId, checkLoginError := CheckLogin(r)
	if checkLoginError != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}
	if loginStatus {
		fmt.Println("logged in userid is: ", userId)
		return
	}
	err := r.ParseForm()
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.BadRequestError)
		return
	}
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	if len(username) == 0 || len(email) == 0 || len(password) == 0 {
		errorManagement.HandleErrorPage(w, r, errorManagement.BadRequestError)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}

	// Insert a record while checking duplicates
	userId, insertError := insertUser(username, email, string(hashedPassword))
	if insertError != nil {
		if errors.Is(insertError, sql.ErrNoRows) {
			// todo show toast
			fmt.Println("User already exists!")
		} else {
			errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		}
		return
	} else {
		fmt.Println("User added successfully!")
	}

	sessionGenerator(w, r, userId)

	// // Query the records
	// rows, err := db.Query(`SELECT id, name, age FROM users;`)
	// if err != nil {
	// 	fmt.Println("Error querying records:", err)
	// 	return
	// }
	// defer rows.Close()

	// fmt.Println("Users:")
	// for rows.Next() {
	// 	var id int
	// 	var name string
	// 	var age int
	// 	err = rows.Scan(&id, &name, &age)
	// 	if err != nil {
	// 		fmt.Println("Error scanning record:", err)
	// 		return
	// 	}
	// 	fmt.Printf("ID: %d, Name: %s, Age: %d\n", id, name, age)
	// }

	// if err = rows.Err(); err != nil {
	// 	fmt.Println("Row iteration error:", err)
	// }

	tmpl, err := template.ParseFiles(
		publicUrl + "authPage.html",
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorManagement.HandleErrorPage(w, r, errorManagement.MethodNotAllowedError)
		return
	}

	loginStatus, userId, checkLoginError := CheckLogin(r)
	if checkLoginError != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}
	if loginStatus {
		fmt.Println("logged in userid is: ", userId)
		return
	}

	err := r.ParseForm()
	if err != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.BadRequestError)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if len(username) == 0 || len(password) == 0 {
		errorManagement.HandleErrorPage(w, r, errorManagement.BadRequestError)
		return
	}

	// Insert a record while checking duplicates
	authStatus, userId, authError := authenticateUser(username, password)
	if authError != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
	} else if authStatus {
		fmt.Println("User logged in successfully!")
	}

	sessionGenerator(w, r, userId)

	tmpl, err := template.ParseFiles(
		publicUrl + "home.html",
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

func sessionGenerator(w http.ResponseWriter, r *http.Request, userId int) {
	sessionToken, expirationTime, insertError := sessionInsert(userId)
	fmt.Println("expirationTime is: ")
	fmt.Println(expirationTime)
	if insertError != nil {
		errorManagement.HandleErrorPage(w, r, errorManagement.InternalServerError)
		return
	}

	// Set the session token in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expirationTime,
		HttpOnly: true,
		Secure:   true,
	})
}

// Middleware to check for valid user session in cookie
func CheckLogin(r *http.Request) (bool, int, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return false, -1, nil
	}

	sessionToken := cookie.Value
	userId, expirationTime, selectError := sessionSelect(sessionToken)
	if selectError != nil {
		return false, -1, selectError
	}
	fmt.Println("expirationTime")
	fmt.Println(expirationTime)

	// Check if the cookie has expired
	if time.Now().After(expirationTime) {
		// Cookie expired, redirect to login
		return false, -1, nil
	}

	fmt.Println("helloe")

	return true, userId, nil
}

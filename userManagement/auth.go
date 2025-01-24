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
	if r.Method != http.MethodPost {
		errorManagement.HandleErrorPage(w, r, errorManagement.MethodNotAllowedError)
		return
	}

	loginStatus, userId, checkLoginError := CheckLogin(r)
	fmt.Println(loginStatus)
	fmt.Println(userId)
	fmt.Println(checkLoginError)
	if checkLoginError != nil {
		//internal server error
		fmt.Println("error checking the error")
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
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		return
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("error method not allowed")
		return
	}

	loginStatus, userId, checkLoginError := CheckLogin(r)
	if checkLoginError != nil {
		//internal server error
		fmt.Println("error checking the error")
		return
	}
	if loginStatus {
		fmt.Println("logged in userid is: ", userId)
		return
	}
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error Parsing the form")
		//handleErrorPage(w, r, BadRequestError)
		return
	}
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	if len(username) == 0 || len(email) == 0 || len(password) == 0 {
		// handleErrorPage(w, r, BadRequestError)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return
	}

	// // Create a table
	// createTableQuery := `
	//     CREATE TABLE IF NOT EXISTS users (
	//         id INTEGER PRIMARY KEY AUTOINCREMENT,
	//         name TEXT,
	//         age INTEGER
	//     );
	// `
	// _, err = db.Exec(createTableQuery)
	// if err != nil {
	// 	fmt.Println("Error creating table:", err)
	// 	return
	// }

	// Insert a record while checking duplicates
	userId, insertError := insertUser(username, email, string(hashedPassword))
	if insertError != nil {
		if errors.Is(insertError, sql.ErrNoRows) {
			fmt.Println("User already exists!")
		} else {
			fmt.Println("Error inserting user:", insertError)
		}
		return
	} else {
		fmt.Println("User added successfully!")
	}

	sessionGenerator(w, userId)

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
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("error method not allowed")
		return
	}

	loginStatus, userId, checkLoginError := CheckLogin(r)
	if checkLoginError != nil {
		//internal server error
		fmt.Println("error checking the error")
		return
	}
	if loginStatus {
		fmt.Println("logged in userid is: ", userId)
		return
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Println("error Parsing the form")
		//handleErrorPage(w, r, BadRequestError)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if len(username) == 0 || len(password) == 0 {
		fmt.Println("empty username or password")
		// handleErrorPage(w, r, BadRequestError)
		return
	}

	// // Create a table
	// createTableQuery := `
	//     CREATE TABLE IF NOT EXISTS users (
	//         id INTEGER PRIMARY KEY AUTOINCREMENT,
	//         name TEXT,
	//         age INTEGER
	//     );
	// `
	// _, err = db.Exec(createTableQuery)
	// if err != nil {
	// 	fmt.Println("Error creating table:", err)
	// 	return
	// }

	// Insert a record while checking duplicates
	authStatus, userId, authError := authenticateUser(username, password)
	if authError != nil {
		// if errors.Is(err, sql.ErrNoRows) {
		// 	fmt.Println("User already exists!")
		// } else {
		// 	fmt.Println("Error authentication user:", err)
		// }
		fmt.Println(authError.Error())
	} else if authStatus {
		fmt.Println("User logged in successfully!")
	}

	sessionGenerator(w, userId)
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
		publicUrl + "home.html",
	)
	if err != nil {
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		return
	}
}

func sessionGenerator(w http.ResponseWriter, userId int) {
	sessionToken, expirationTime, insertError := sessionInsert(userId)
	fmt.Println("expirationTime is: ")
	fmt.Println(expirationTime)
	if insertError != nil {
		// if errors.Is(err, sql.ErrNoRows) {
		// 	fmt.Println("User already exists!")
		// } else {
		// 	fmt.Println("Error inserting user:", err)
		// }
		fmt.Println(insertError.Error())
		return
	}
	fmt.Println("session created successfully!")

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
	// Get the cookie named "session_token"
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

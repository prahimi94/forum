package userManagement

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const publicUrl = "frontend/public/"

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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
	err = insertUser(username, email, string(hashedPassword))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("User already exists!")
		} else {
			fmt.Println("Error inserting user:", err)
		}
	} else {
		fmt.Println("User added successfully!")
	}

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
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error Parsing the form")
		//handleErrorPage(w, r, BadRequestError)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if len(username) == 0 || len(password) == 0 {
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
	authStatus, authError := authenticateUser(username, password)
	if err != nil {
		// if errors.Is(err, sql.ErrNoRows) {
		// 	fmt.Println("User already exists!")
		// } else {
		// 	fmt.Println("Error inserting user:", err)
		// }
		fmt.Println(authError.Error())
	} else if authStatus {
		fmt.Println("User logged in successfully!")
	}

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

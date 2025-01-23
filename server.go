package main

import (
	"forum/backend"
	"forum/userManagement"
	"log"
	"net/http"
)

// main initializes the HTTP server, registers routes, and starts listening for incoming requests.
func main() {
	http.Handle("/static/", http.FileServer(http.Dir("frontend/public/")))
	http.Handle("/img/", http.FileServer(http.Dir("frontend/public/")))

	// Register route handlers
	http.HandleFunc("/", backend.MainPageHandler)
	http.HandleFunc("/auth/", userManagement.AuthHandler)
	http.HandleFunc("/register", userManagement.RegisterHandler)
	http.HandleFunc("/login", userManagement.LoginHandler)
	//start the server on port 8080
	log.Println("Starting server on: http://localhost:8080")
	log.Println("Status ok: ", http.StatusOK)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

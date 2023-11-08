package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type userData struct {
	Email    string
	Password string
}

func treatError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// Define you MySQL credentials and what database will be used
const dsn = "root:password@tcp(localhost:3306)/loginscreenusers"

// Given an email and a hashed password, checks on the database
// if the user is present, and returns a boolean value
func consultOnDatabase(email string, hashedPassword []byte) bool {
	var storedEmail, storedPassword string

	stringedPassword := hex.EncodeToString(hashedPassword)

	db, err := sql.Open("mysql", dsn)

	treatError(err)
	defer db.Close()

	// Here I used a prepared statement to prevent SQL injection
	query := "SELECT email, password FROM users WHERE email = ? AND password = ?"
	row := db.QueryRow(query, email, stringedPassword)

	err = row.Scan(&storedEmail, &storedPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			// No record found, indicating invalid login.
			return false
		}
		treatError(err)
		return false
	}

	// Matching record found, indicating a valid login.
	return true
}

// Given an email and hashed password, save then on the database
func saveOnDatabase(email string, hashedPassword []byte) {
	stringedPassword := hex.EncodeToString(hashedPassword)

	// Open the database and defer it's closeure.
	database, err := sql.Open("mysql", dsn)
	treatError(err)
	defer database.Close()

	// Create table if it doesn't exist
	_, err = database.Exec("CREATE TABLE IF NOT EXISTS users(id INT AUTO_INCREMENT PRIMARY KEY, email VARCHAR(255), password VARCHAR(64))")
	treatError(err)

	// Insert data using prepared statements to prevent SQL Injection.
	_, err = database.Exec("INSERT INTO users (email, password) VALUES (?, ?)", email, stringedPassword)
	treatError(err)
}

func main() {
	// !! This is a DUMB way to do this, but I don't know a better way !!
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/style.css")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template, err := template.ParseFiles("static/index.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := template.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		// Print the data												!!!!
		fmt.Printf("Email: %s\nPassword: %s\n", email, password)

		// Encrypts the password
		passwordEncrypted := sha256.New()
		passwordEncrypted.Write([]byte(password))
		hashedPassword := passwordEncrypted.Sum(nil)

		// Logs in the server console the hashed password				!!!!
		fmt.Printf("%x\n", hashedPassword)

		isOnDatabase := consultOnDatabase(email, hashedPassword)

		if !isOnDatabase {
			// Call saveOnDatabase to store user data
			// if the user isn't on the database alredy
			saveOnDatabase(email, hashedPassword)

			// and redirect the user to the login screen
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			fmt.Println("USER ALREDY ON DATABASE")
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

	})

	http.HandleFunc("/submitlogin", func(w http.ResponseWriter, r *http.Request) {

	})

	http.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {

	})

	fmt.Println("Server is running on: 8080")
	http.ListenAndServe(":8080", nil)
}

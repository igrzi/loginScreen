package main

import (
	"database/sql"
	"encoding/hex"
	"fmt"

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

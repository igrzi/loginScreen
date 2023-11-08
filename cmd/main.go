package main

import (
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

// Given an email and a hashed password, checks on the database if the user is present, 
// and returns a boolean value
func consultOnDatabase(email string, hashedPassword []byte) bool {
	...
}

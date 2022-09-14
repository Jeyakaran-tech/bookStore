package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
)

// User is a model for the user entity
type User struct {
	gorm.Model

	Name string
}

// createGetUserHandler returns a user handler function that
// returns the first user's name from the DB to the http caller
func createGetUserHandler(db *gorm.DB) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		var user User
		db.First(&user)

		// write user name in the response
		rw.Write([]byte(user.Name))
	}
}

func main() {
	// get our OS variables
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")

	// connect to the DB
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf("host=0.0.0.0 port=%s user=%s password=%s sslmode=disable", os.Getenv("PORT"), user, pass),
	)
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{}) // create tables

	// create dummy user
	db.Create(&User{
		Name: "Peter",
	})

	http.HandleFunc("/", createGetUserHandler(db))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

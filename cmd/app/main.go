package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Jeyakaran-tech/bookStore/cloudsql"
	"github.com/Jeyakaran-tech/bookStore/dbconnect"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)

	users := &User{Database: dbconnect.GetDB()}

	mux := mux.NewRouter()
	mux.Handle("/v1/books/", cloudsql.RootHandler(users.ListOfBooks)).Methods("GET")
	mux.Handle("/v1/books/{book-id}", cloudsql.RootHandler(users.GetOrUpdateBook))
	mux.Handle("/v1/books", cloudsql.RootHandler(users.InsertBook)).Methods("POST")
	mux.Handle("/v1/books", cloudsql.RootHandler(users.GetBookWithWildCard)).Methods("GET")
	mux.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

type User struct {
	Database *sql.DB
}

func (u *User) ListOfBooks(w http.ResponseWriter, r *http.Request) error {
	return cloudsql.GetBooks(w, r, u.Database)
}

func (u *User) GetOrUpdateBook(w http.ResponseWriter, r *http.Request) error {
	return cloudsql.GetOrUpdateBook(w, r, u.Database)
}

func (u *User) InsertBook(w http.ResponseWriter, r *http.Request) error {
	return cloudsql.InsertBook(w, r, u.Database)
}

func (u *User) GetBookWithWildCard(w http.ResponseWriter, r *http.Request) error {
	return cloudsql.GetBookWithWildCard(w, r, u.Database)
}

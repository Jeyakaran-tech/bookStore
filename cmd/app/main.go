package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"

	"github.com/Jeyakaran-tech/bookStore/cloudsql"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	r := chi.NewRouter()

	r.Get("/v1/books/", cloudsql.Books)
	r.Post("/v1/books", cloudsql.Books)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}

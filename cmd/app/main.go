package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Jeyakaran-tech/bookStore/cloudsql"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)

	http.HandleFunc("/v1/books", cloudsql.Books)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

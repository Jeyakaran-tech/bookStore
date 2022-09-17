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

	mux := http.NewServeMux()
	mux.Handle("/v1/books/", RootHandler(cloudsql.ListOfBooks))
	mux.Handle("/v1/books/{book-id}", RootHandler(cloudsql.GetBook))
	mux.Handle("/v1/books", RootHandler(cloudsql.InsertBook))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

type RootHandler func(http.ResponseWriter, *http.Request) error

// rootHandler implements http.Handler interface.
func (fn RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r) // Call handler function
	if err == nil {
		return
	}
	// This is where our error handling logic starts.
	log.Printf("An error occured: %v", err) // Log the error.

	clientError, ok := err.(cloudsql.ClientError) // Check if it is a ClientError.
	if !ok {
		// If the error is not ClientError, assume that it is ServerError.
		w.WriteHeader(500) // return 500 Internal Server Error.
		return
	}

	body, err := clientError.ResponseBody() // Try to get response body of ClientError.
	if err != nil {
		log.Printf("An error accured: %v", err)
		w.WriteHeader(500)
		return
	}
	status, headers := clientError.ResponseHeaders() // Get http status code and headers.
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)
	w.Write(body)
}

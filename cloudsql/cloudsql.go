package cloudsql

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/Jeyakaran-tech/bookStore/types"
)

func ListOfBooks(w http.ResponseWriter, r *http.Request) {

	db := getDB()
	var books []types.Books
	w.Header().Set("content-type", "application/json")

	listOfBooks, err := db.Query("SELECT * FROM bookstore")
	if err != nil {
		log.Fatalf("DB.QueryRow: %v", err)
		http.Error(w, "can't select the table", http.StatusBadRequest)
		return
	}
	defer listOfBooks.Close()

	for listOfBooks.Next() {
		var (
			ID             int
			Author         string
			Published_date string
			Price          float64
			InStock        bool
			Time_added     time.Time
		)
		err := listOfBooks.Scan(&ID, &Author, &Published_date, &Price, &InStock, &Time_added)
		if err != nil {
			log.Fatalf("Rows.Scan: %v", err)
			http.Error(w, "can't scan the rows from Database", http.StatusBadRequest)
			return
		}
		books = append(books, types.Books{
			ID:             ID,
			Author:         Author,
			Published_date: Published_date,
			Price:          Price,
			InStock:        InStock,
			Time_added:     Time_added,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}

func InsertBook(w http.ResponseWriter, r *http.Request) {

	db := getDB()
	var books types.Books
	w.Header().Set("content-type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(badRequest())
		return
	}

	if err := json.Unmarshal(body, &books); err != nil {
		log.Fatalf("Cant unmarshal while reading the request body, %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(badRequest())
		return
	}

	insertVote := "INSERT INTO bookstore(Author,Published_date,Price,In_Stock, time_added) VALUES(?,?,?,?, NOW())"
	date, dateErr := time.Parse("2006-01-02", books.Published_date)
	if dateErr != nil {
		log.Printf("Error parsing date: %v", dateErr)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(badRequest())
		return
	}

	if _, err := db.Exec(insertVote, books.Author, date, books.Price, books.InStock); err != nil {
		log.Fatalf("Cant insert into table, %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(badRequest())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse())
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	db := getDB()
	bookID := path.Base(r.URL.Path)

	var books types.Books
	w.Header().Set("content-type", "application/json")

	book, err := db.Query(fmt.Sprintf("SELECT * FROM bookstore where ID=%s", bookID))
	if err != nil {
		log.Fatalf("DB.QueryRow: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(internalServerError())
		return
	}
	defer book.Close()

	var (
		ID             int
		Author         string
		Published_date string
		Price          float64
		InStock        bool
		Time_added     time.Time
	)
	scanErr := book.Scan(&ID, &Author, &Published_date, &Price, &InStock, &Time_added)
	if scanErr != nil {
		log.Fatalf("Rows.Scan: %v", scanErr)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(badRequest())
		return
	}
	books = types.Books{
		ID:             ID,
		Author:         Author,
		Published_date: Published_date,
		Price:          Price,
		InStock:        InStock,
		Time_added:     Time_added,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)

}

func badRequest() types.Status {
	return types.Status{
		Code:        "8400",
		Description: "Bad Request",
	}
}

func successResponse() types.Status {
	return types.Status{
		Code:        "8200",
		Description: "Success",
	}
}

func internalServerError() types.Status {
	return types.Status{
		Code:        "8500",
		Description: "Internal Server Error",
	}
}

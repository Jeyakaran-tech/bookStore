package cloudsql

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"time"

	"github.com/Jeyakaran-tech/bookStore/types"
)

func ListOfBooks(w http.ResponseWriter, r *http.Request) error {

	db := getDB()
	var books []types.Books
	w.Header().Set("content-type", "application/json")

	listOfBooks, err := db.Query("SELECT * FROM bookstore")
	if err != nil {
		return NewHTTPError(err, 400, "Can't pull the data from table")
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
			return NewHTTPError(err, 400, "Can't pull the data from table")
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
	return nil
}

func InsertBook(w http.ResponseWriter, r *http.Request) error {

	db := getDB()
	var books types.Books
	w.Header().Set("content-type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return NewHTTPError(err, 400, "Error reading the request body")
	}

	if err := json.Unmarshal(body, &books); err != nil {
		return NewHTTPError(err, 400, "Error unmarshaling the request body")
	}

	insertVote := "INSERT INTO bookstore(Author,Published_date,Price,In_Stock, time_added) VALUES(?,?,?,?, NOW())"
	date, dateErr := time.Parse("2006-01-02", books.Published_date)
	if dateErr != nil {
		return NewHTTPError(dateErr, 400, "Error parsing date")
	}

	if _, err := db.Exec(insertVote, books.Author, date, books.Price, books.InStock); err != nil {
		return NewHTTPError(dateErr, 400, "Insertion error")
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.Status{
		Code:        "8200",
		Description: "Success",
	})
	return nil
}

func GetBook(w http.ResponseWriter, r *http.Request) error {
	db := getDB()
	bookID := path.Base(r.URL.Path)
	var (
		ID             int
		Author         string
		Published_date string
		Price          float64
		InStock        bool
		Time_added     time.Time
	)

	w.Header().Set("content-type", "application/json")

	scanErr := db.QueryRow(fmt.Sprintf("SELECT * FROM bookstore where ID=%s", bookID)).Scan(&ID, &Author, &Published_date, &Price, &InStock, &Time_added)
	if scanErr != nil {
		return NewHTTPError(nil, 400, "Error when selecting rows")
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.Books{
		ID:             ID,
		Author:         Author,
		Published_date: Published_date,
		Price:          Price,
		InStock:        InStock,
		Time_added:     Time_added,
	})
	return nil

}

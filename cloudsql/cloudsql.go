package cloudsql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"regexp"
	"time"

	"github.com/Jeyakaran-tech/bookStore/types"
)

func GetBooks(w http.ResponseWriter, r *http.Request, db *sql.DB) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if r.Method != "GET" {
		return NewHTTPError(fmt.Errorf("invalid method passed"), 400, "Bad Request")
	}
	var books []types.Books
	w.Header().Set("content-type", "application/json")

	listOfBooks, err := db.QueryContext(ctx, "SELECT * FROM bookstore")
	if err != nil {
		return NewHTTPError(err, 400, "Can't pull the data from table")
	}
	defer listOfBooks.Close()

	for listOfBooks.Next() {
		var (
			ID             int
			Name           string
			Author         string
			Published_date string
			Price          float64
			InStock        bool
			Time_added     time.Time
		)
		err := listOfBooks.Scan(&ID, &Name, &Author, &Published_date, &Price, &InStock, &Time_added)
		if err != nil {
			return NewHTTPError(err, 400, "Can't pull the data from table")
		}
		books = append(books, types.Books{
			ID:             ID,
			Author:         Author,
			Name:           Name,
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

func InsertBook(w http.ResponseWriter, r *http.Request, db *sql.DB) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if r.Method != "POST" {
		return NewHTTPError(fmt.Errorf("invalid method passed"), 400, "Bad Request")
	}
	var books types.Books
	w.Header().Set("content-type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return NewHTTPError(err, 400, "Error reading the request body")
	}

	if err := json.Unmarshal(body, &books); err != nil {
		return NewHTTPError(err, 400, "Error unmarshaling the request body")
	}

	insertBook := "INSERT INTO bookstore(Author,Name, Published_date,Price,In_Stock, time_added) VALUES(?,?,?,?,?, NOW())"
	date, dateErr := time.Parse("2006-01-02", books.Published_date)
	if dateErr != nil {
		return NewHTTPError(dateErr, 400, "Error parsing date")
	}

	if _, insertErr := db.ExecContext(ctx, insertBook, books.Author, books.Name, date, books.Price, books.InStock); insertErr != nil {
		return NewHTTPError(insertErr, 400, "Insertion error")
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		types.Status{
			Code:        "8200",
			Description: "Success",
		},
	)
	return nil
}

func GetOrUpdateBook(w http.ResponseWriter, r *http.Request, db *sql.DB) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bookID := path.Base(r.URL.Path)
	var (
		ID             int
		Author         string
		Name           string
		Published_date string
		Price          float64
		InStock        bool
		Time_added     time.Time
	)

	w.Header().Set("content-type", "application/json")

	switch r.Method {
	case "GET":

		if r.Method != "GET" {
			return NewHTTPError(fmt.Errorf("invalid method passed"), 400, "Bad Request")
		}

		scanErr := db.QueryRowContext(ctx, fmt.Sprintf("SELECT * FROM bookstore where ID=%s", bookID)).Scan(&ID, &Name, &Author, &Published_date, &Price, &InStock, &Time_added)
		if scanErr != nil {
			return NewHTTPError(scanErr, 400, "Error when selecting rows")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(types.Books{
			ID:             ID,
			Author:         Author,
			Name:           Name,
			Published_date: Published_date,
			Price:          Price,
			InStock:        InStock,
			Time_added:     Time_added,
		})
		return nil

	case "PUT":

		if r.Method != "PUT" {
			return NewHTTPError(fmt.Errorf("invalid method passed"), 400, "Bad Request")
		}
		var books types.Books
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return NewHTTPError(err, 400, "Error reading the request body")
		}

		if err := json.Unmarshal(body, &books); err != nil {
			return NewHTTPError(err, 400, "Error unmarshaling the request body")
		}

		updateBook := fmt.Sprintf("UPDATE bookstore SET Author = ? ,Name = ?, Published_date = ?, Price = ?, In_Stock = ?, time_added = NOW() WHERE ID = %s", bookID)
		date, dateErr := time.Parse("2006-01-02", books.Published_date)
		if dateErr != nil {
			return NewHTTPError(dateErr, 400, "Error parsing date")
		}

		if _, updateErr := db.ExecContext(ctx, updateBook, books.Author, books.Name, date, books.Price, books.InStock); err != nil {
			return NewHTTPError(updateErr, 400, "Updation error")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(types.Status{
			Code:        "8200",
			Description: "Updated successfully",
		})
		return nil

	default:
		return NewHTTPError(fmt.Errorf("bad request"), 400, "Invalid Method")
	}

}

func GetBookWithWildCard(w http.ResponseWriter, r *http.Request, db *sql.DB) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if r.Method != "GET" {
		return NewHTTPError(fmt.Errorf("invalid method passed"), 400, "Bad Request")
	}

	var books []types.Books
	w.Header().Set("content-type", "application/json")

	name := r.URL.Query().Get("name")
	if name == "" {
		return NewHTTPError(fmt.Errorf("missing query params"), 400, "Missing Query Parameter")
	}

	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	bookName := re.ReplaceAllString(name, "%")

	listOfBooks, err := db.QueryContext(ctx, fmt.Sprintf("SELECT * from bookstore where name like '%s'", bookName))
	if err != nil {
		return NewHTTPError(err, 400, "Can't pull the data from table")
	}
	defer listOfBooks.Close()

	for listOfBooks.Next() {
		var (
			ID             int
			Author         string
			Name           string
			Published_date string
			Price          float64
			InStock        bool
			Time_added     time.Time
		)
		err := listOfBooks.Scan(&ID, &Name, &Author, &Published_date, &Price, &InStock, &Time_added)
		if err != nil {
			return NewHTTPError(err, 400, "Can't pull the data from table")
		}
		books = append(books, types.Books{
			ID:             ID,
			Author:         Author,
			Name:           Name,
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

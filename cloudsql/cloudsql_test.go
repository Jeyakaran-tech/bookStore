package cloudsql

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Jeyakaran-tech/bookStore/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gorilla/mux"
)

var readBook = &types.Books{
	ID:             1,
	Name:           "sample name",
	Author:         "Sample author",
	Published_date: "1998-10-12",
	Price:          500.00,
	InStock:        true,
	Time_added:     time.Now(),
}

var writeBook = &types.Books{
	Name:           "sample name",
	Author:         "Sample author",
	Published_date: "1998-10-12",
	Price:          500.00,
	InStock:        true,
	Time_added:     time.Now(),
}

var _ = Describe("Cloudsql", func() {
	Context("Get all the books", func() {

		It("should return success", func() {

			db, mock := NewMock()
			defer db.Close()

			users := &User{Database: db}
			req, _ := http.NewRequest("GET", "/v1/books/", nil)
			rr := httptest.NewRecorder()
			mux := mux.NewRouter()
			rows := sqlmock.NewRows([]string{"id", "name", "author", "publishedDate", "Price", "InStock", "timeAdded"}).
				AddRow(readBook.ID, readBook.Name, readBook.Author, readBook.Published_date, readBook.Price, readBook.InStock, readBook.Time_added)
			mock.ExpectQuery("SELECT * FROM bookstore").WillReturnRows(rows)

			mux.Handle("/v1/books/", RootHandler(users.ListOfBooks))
			mux.ServeHTTP(rr, req)

			Expect(rr.Code).To(BeEquivalentTo(200))
			mock.ExpectClose()
			db.Close()

		})
	})

	Context("Get all the books - passed with invalid method", func() {

		It("should return success", func() {

			db, _ := NewMock()
			defer db.Close()

			users := &User{Database: db}
			req, _ := http.NewRequest("DELETE", "/v1/books/", nil)
			rr := httptest.NewRecorder()
			mux := mux.NewRouter()
			mux.Handle("/v1/books/", RootHandler(users.ListOfBooks))
			mux.ServeHTTP(rr, req)

			Expect(rr.Code).To(BeEquivalentTo(400))

		})
	})

	Context("Insert book", func() {

		It("should return success", func() {

			db, mock := NewMock()
			defer db.Close()

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(writeBook)
			if err != nil {
				log.Fatal(err)
			}

			users := &User{Database: db}
			req, _ := http.NewRequest("POST", "/v1/books", &buf)

			rr := httptest.NewRecorder()
			mux := mux.NewRouter()

			mock.ExpectExec("INSERT INTO bookstore(Author,Name, Published_date,Price,In_Stock, time_added) VALUES(?,?,?,?,?, NOW())").WillReturnResult(driver.ResultNoRows)

			mux.Handle("/v1/books", RootHandler(users.InsertBook))
			mux.ServeHTTP(rr, req)

			Expect(rr.Code).To(BeEquivalentTo(200))
			mock.ExpectClose()
			db.Close()

		})
	})

	Context("Edit book", func() {

		It("should return success", func() {

			db, mock := NewMock()
			defer db.Close()

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(writeBook)
			if err != nil {
				log.Fatal(err)
			}

			users := &User{Database: db}
			req, _ := http.NewRequest("PUT", "/v1/books/101", &buf)

			rr := httptest.NewRecorder()
			mux := mux.NewRouter()

			mock.ExpectExec("UPDATE bookstore SET Author = ? ,Name = ?, Published_date = ?, Price = ?, In_Stock = ?, time_added = NOW() WHERE ID = %s").WillReturnResult(driver.ResultNoRows)

			mux.Handle("/v1/books/{book-id}", RootHandler(users.GetOrUpdateBook))
			mux.ServeHTTP(rr, req)

			Expect(rr.Code).To(BeEquivalentTo(200))
			mock.ExpectClose()
			db.Close()

		})
	})

	Context("Get a book by ID", func() {

		It("should return success", func() {

			db, mock := NewMock()
			defer db.Close()

			row := sqlmock.NewRows([]string{"id", "name", "author", "publishedDate", "Price", "InStock", "timeAdded"}).
				AddRow(readBook.ID, readBook.Name, readBook.Author, readBook.Published_date, readBook.Price, readBook.InStock, readBook.Time_added)

			users := &User{Database: db}
			req, _ := http.NewRequest("GET", "/v1/books/1", nil)

			rr := httptest.NewRecorder()
			mux := mux.NewRouter()
			mock.ExpectQuery("SELECT * FROM bookstore where ID=1").WillReturnRows(row)

			mux.Handle("/v1/books/1", RootHandler(users.GetOrUpdateBook))
			mux.ServeHTTP(rr, req)

			Expect(rr.Code).To(BeEquivalentTo(200))

			mock.ExpectClose()
			db.Close()
		})
	})

	Context("Get a book by wildcard", func() {

		It("should return success", func() {

			db, mock := NewMock()
			defer db.Close()

			row := sqlmock.NewRows([]string{"id", "name", "author", "publishedDate", "Price", "InStock", "timeAdded"}).
				AddRow(readBook.ID, readBook.Name, readBook.Author, readBook.Published_date, readBook.Price, readBook.InStock, readBook.Time_added)

			users := &User{Database: db}
			req, _ := http.NewRequest("GET", "/v1/books?name=*uvi", nil)

			rr := httptest.NewRecorder()
			mux := mux.NewRouter()

			mock.ExpectQuery("SELECT * from bookstore where name like '%uvi'").WillReturnRows(row)
			mux.Handle("/v1/books", RootHandler(users.GetBookWithWildCard))
			mux.ServeHTTP(rr, req)

			Expect(rr.Code).To(BeEquivalentTo(200))

		})
	})

})

type User struct {
	Database *sql.DB
}

func (u *User) ListOfBooks(w http.ResponseWriter, r *http.Request) error {
	return GetBooks(w, r, u.Database)
}

func (u *User) GetOrUpdateBook(w http.ResponseWriter, r *http.Request) error {
	return GetOrUpdateBook(w, r, u.Database)
}

func (u *User) InsertBook(w http.ResponseWriter, r *http.Request) error {
	return InsertBook(w, r, u.Database)
}

func (u *User) GetBookWithWildCard(w http.ResponseWriter, r *http.Request) error {
	return GetBookWithWildCard(w, r, u.Database)
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

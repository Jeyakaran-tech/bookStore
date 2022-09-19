package integrationtests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/Jeyakaran-tech/bookStore/cloudsql"
	"github.com/Jeyakaran-tech/bookStore/types"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gorilla/mux"
)

var writeBook = &types.Books{
	Name:           "sample name",
	Author:         "Sample author",
	Published_date: "1998-10-12",
	Price:          500.00,
	InStock:        true,
	Time_added:     time.Now(),
}

var _ = Describe("IntegrationTests", func() {

	db, connErr := connectTCPSocket()
	if connErr != nil {
		log.Fatal(connErr)
	}

	Context("Insert the book", func() {

		It("Must return the inserted ID", func() {
			var buf bytes.Buffer
			var ID int64
			var status types.Status

			err := json.NewEncoder(&buf).Encode(writeBook)
			if err != nil {
				log.Fatal(err)
			}

			users := &User{Database: db}
			req, _ := http.NewRequest("POST", "/v1/books", &buf)
			rr := httptest.NewRecorder()
			mux := mux.NewRouter()

			mux.Handle("/v1/books", cloudsql.RootHandler(users.InsertBook)).Methods("POST")
			mux.ServeHTTP(rr, req)

			if err := json.Unmarshal(rr.Body.Bytes(), &status); err != nil {
				panic(err)
			}
			Expect(status.Code).To(BeEquivalentTo("8200"))

			deleteBook := fmt.Sprintf("DELETE FROM bookstore  WHERE ID = %d", ID)
			if _, updateErr := db.Exec(deleteBook); err != nil {
				log.Fatal(updateErr)
			}

		})
	})

	Context("Get a book with Wildcard", func() {

		It("Must return correct book", func() {
			var ID int64
			var books []types.Books

			users := &User{Database: db}

			insertBook := "INSERT INTO bookstore(Author,Name, Published_date,Price,In_Stock, time_added) VALUES(?,?,?,?,?, NOW())"
			date, dateErr := time.Parse("2006-01-02", writeBook.Published_date)
			if dateErr != nil {
				log.Fatal(dateErr)
			}

			if _, insertErr := db.Exec(insertBook, writeBook.Author, writeBook.Name, date, writeBook.Price, writeBook.InStock); insertErr != nil {
				log.Fatal(insertErr)
			}

			req, _ := http.NewRequest("GET", "/v1/books?name=sample*", nil)
			rr := httptest.NewRecorder()
			mux := mux.NewRouter()

			mux.Handle("/v1/books", cloudsql.RootHandler(users.GetBookWithWildCard))
			mux.ServeHTTP(rr, req)

			listOfBooks, err := db.Query("SELECT * from bookstore where name like 'sample%'")
			if err != nil {
				log.Fatal(err)
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
					log.Fatal(err)
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

			if err := json.Unmarshal(rr.Body.Bytes(), &books); err != nil {
				panic(err)
			}

			Expect(books[0].Name).To(BeEquivalentTo("sample name"))

			deleteBook := fmt.Sprintf("DELETE FROM bookstore  WHERE ID = %d", ID)
			if _, updateErr := db.Exec(deleteBook); err != nil {
				log.Fatal(updateErr)
			}

		})
	})

})

type User struct {
	Database *sql.DB
}

func (u *User) InsertBook(w http.ResponseWriter, r *http.Request) error {
	return cloudsql.InsertBook(w, r, u.Database)
}

func (u *User) GetBookWithWildCard(w http.ResponseWriter, r *http.Request) error {
	return cloudsql.GetBookWithWildCard(w, r, u.Database)
}

func connectTCPSocket() (*sql.DB, error) {

	var (
		dbUser    = "user"
		dbPwd     = "testbooks"
		dbName    = "books"
		dbPort    = "3306"
		dbTCPHost = "127.0.0.1"
	)

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPwd, dbTCPHost, dbPort, dbName)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	return dbPool, nil
}

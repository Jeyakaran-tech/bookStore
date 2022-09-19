package cloudsql

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	"github.com/gorilla/mux"
)

var _ = Describe("Cloudsql", func() {
	Context("Get all the books", func() {

		It("should return success", func() {

			db := getTestDB()
			defer db.Close()

			users := &User{Database: db}
			req, _ := http.NewRequest("GET", "/v1/books/", nil)
			rr := httptest.NewRecorder()
			mux := mux.NewRouter()
			mux.Handle("/v1/books/", RootHandler(users.ListOfBooks))
			mux.ServeHTTP(rr, req)

			Expect(rr.Code).To(BeEquivalentTo(200))

		})
	})

	Context("Get all the books - passed with invalid method", func() {

		It("should return success", func() {

			db := getTestDB()
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
			jsonBody := []byte(`{
				"author": "love life",
				"publishedDate": "1998-08-12",
				"price": 500,
				"inStock": true
			}`)
			bodyReader := bytes.NewReader(jsonBody)

			db := getTestDB()
			defer db.Close()

			users := &User{Database: db}
			req, _ := http.NewRequest("POST", "/v1/books", bodyReader)

			rr := httptest.NewRecorder()
			mux := mux.NewRouter()
			mux.Handle("/v1/books", RootHandler(users.InsertBook))
			mux.ServeHTTP(rr, req)

			Expect(rr.Code).To(BeEquivalentTo(200))

		})
	})

	Context("Edit book", func() {

		It("should return success", func() {
			jsonBody := []byte(`{
				"author": "love life",
				"publishedDate": "1998-08-12",
				"price": 500,
				"inStock": true
			}`)
			bodyReader := bytes.NewReader(jsonBody)

			db := getTestDB()
			defer db.Close()

			users := &User{Database: db}
			req, _ := http.NewRequest("PUT", "/v1/books/101", bodyReader)

			rr := httptest.NewRecorder()
			mux := mux.NewRouter()
			mux.Handle("/v1/books/{book-id}", RootHandler(users.GetOrUpdateBook))
			mux.ServeHTTP(rr, req)

			Expect(rr.Code).To(BeEquivalentTo(200))

		})
	})

	Context("Get a book by ID", func() {

		It("should return success", func() {

			db := getTestDB()
			defer db.Close()

			users := &User{Database: db}
			req, _ := http.NewRequest("GET", "/v1/books/1", nil)

			rr := httptest.NewRecorder()
			mux := mux.NewRouter()
			mux.Handle("/v1/books/1", RootHandler(users.GetOrUpdateBook))
			mux.ServeHTTP(rr, req)

			Expect(rr.Code).To(BeEquivalentTo(200))

		})
	})

	Context("Get a book by wildcard", func() {

		It("should return success", func() {
			jsonBody := []byte(`{
				"author": "love life",
				"publishedDate": "1998-08-12",
				"price": 500,
				"inStock": true
			}`)
			bodyReader := bytes.NewReader(jsonBody)

			db := getTestDB()
			defer db.Close()

			users := &User{Database: db}
			req, _ := http.NewRequest("GET", "/v1/books?name=*uvi", bodyReader)

			rr := httptest.NewRecorder()
			mux := mux.NewRouter()
			mux.Handle("/v1/books", RootHandler(users.GetBookWithWildCard))
			mux.ServeHTTP(rr, req)

			Expect(rr.Code).To(BeEquivalentTo(200))

		})
	})

})

type User struct {
	Database *sql.DB
}

func getTestDB() *sql.DB {
	cfg := mysql.Cfg("bookstore-362511:australia-southeast2:bookstore", "user", "testbooks")
	cfg.DBName = "books"
	cfg.ParseTime = true

	const timeout = 10 * time.Second
	cfg.Timeout = timeout
	cfg.ReadTimeout = timeout
	cfg.WriteTimeout = timeout

	db, err := mysql.DialCfg(cfg)
	if err != nil {
		panic("couldn't dial: " + err.Error())
	}
	return db
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

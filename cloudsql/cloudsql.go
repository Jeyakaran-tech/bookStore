package cloudsql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/Jeyakaran-tech/bookStore/types"
	"github.com/go-sql-driver/mysql"
)

// func Books(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		listOfBooks(w, r, getDB())
// 	case http.MethodPost:
// 		insertBook(w, r, getDB())
// 	default:
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 	}
// }

func ListOfBooks(w http.ResponseWriter, r *http.Request) {

	db := getDB()
	var books []types.Books
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
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(books)
}

func InsertBook(w http.ResponseWriter, r *http.Request) {

	db := getDB()
	var books types.Books
	w.Header().Set("Content-Type", "application/json")

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

	insertVote := "INSERT INTO bookstore(Author,Published_date,Price,In_Stock, created_at) VALUES(?,?,?,?, NOW())"
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

func connectWithConnector() (*sql.DB, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Warning: %s environment variable not set.", k)
		}
		return v
	}

	var (
		dbUser                 = mustGetenv("DB_USER")
		dbPwd                  = mustGetenv("DB_PASS")
		dbName                 = mustGetenv("DB_NAME")
		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME")
		usePrivate             = os.Getenv("PRIVATE_IP")
	)

	d, err := cloudsqlconn.NewDialer(context.Background())
	if err != nil {
		return nil, fmt.Errorf("cloudsqlconn.NewDialer: %v", err)
	}
	mysql.RegisterDialContext("cloudsqlconn",
		func(ctx context.Context, addr string) (net.Conn, error) {
			if usePrivate != "" {
				return d.Dial(ctx, instanceConnectionName, cloudsqlconn.WithPrivateIP())
			}
			return d.Dial(ctx, instanceConnectionName)
		})

	dbURI := fmt.Sprintf("%s:%s@cloudsqlconn(localhost:3306)/%s?parseTime=true",
		dbUser, dbPwd, dbName)

	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}
	return dbPool, nil
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

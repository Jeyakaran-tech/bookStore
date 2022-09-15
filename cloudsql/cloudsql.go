package cloudsql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/Jeyakaran-tech/bookStore/types"
	"github.com/go-sql-driver/mysql"
)

func Books(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listOfBooks(w, r, getDB())
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func listOfBooks(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var books []types.Books
	listOfBooks, err := db.Query("SELECT * FROM bookstore")
	if err != nil {
		log.Fatalf("DB.QueryRow: %v", err)
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
	booksData, err := json.Marshal(books)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprint(w, string(booksData))
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
		dbUser                 = mustGetenv("DB_USER")                  // e.g. 'my-db-user'
		dbPwd                  = mustGetenv("DB_PASS")                  // e.g. 'my-db-password'
		dbName                 = mustGetenv("DB_NAME")                  // e.g. 'my-database'
		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
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

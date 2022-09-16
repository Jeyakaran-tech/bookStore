package cloudsql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/go-sql-driver/mysql"
)

func getDB() *sql.DB {
	once.Do(func() {
		db = mustConnect()
	})
	return db
}

func mustConnect() *sql.DB {
	var (
		db  *sql.DB
		err error
	)

	// Use the connector when INSTANCE_CONNECTION_NAME (proj:region:instance) is defined.
	if os.Getenv("INSTANCE_CONNECTION_NAME") != "" {
		db, err = connectWithConnector()
		if err != nil {
			log.Fatalf("connectConnector: unable to connect: %s", err)
		}
	}

	if db == nil {
		log.Fatal("Missing database connection type - INSTANCE_CONNECTION_NAME")
	}

	if err := createTable(db); err != nil {
		log.Fatalf("unable to create table: %s", err)
	}

	return db
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

func createTable(db *sql.DB) error {
	createBook := `CREATE TABLE IF NOT EXISTS bookstore (
		ID SERIAL NOT NULL, 
		Author VARCHAR(50) NOT NULL, 
		Published_date DATE NOT NULL, 
		Price DOUBLE NOT NULL, 
		In_Stock BOOL NOT NULL, 
		time_added DATETIME NOT NULL
	);`
	_, err := db.Exec(createBook)
	return err
}

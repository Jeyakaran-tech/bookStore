package cloudsql

import (
	"database/sql"
	"log"
	"os"
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
		log.Fatal("Missing database connection type. Please define one of INSTANCE_HOST, INSTANCE_UNIX_SOCKET, or INSTANCE_CONNECTION_NAME")
	}

	if err := createTable(db); err != nil {
		log.Fatalf("unable to create table: %s", err)
	}

	return db
}

func createTable(db *sql.DB) error {
	createBook := `CREATE TABLE IF NOT EXISTS bookstore (
		ID SERIAL NOT NULL, 
		Author VARCHAR(50) NOT NULL, 
		Published_date DATE NOT NULL, 
		Price DOUBLE NOT NULL, 
		In_Stock BOOL NOT NULL, 
		time_added DATETIME NOT NULL)
	);`
	_, err := db.Exec(createBook)
	return err
}

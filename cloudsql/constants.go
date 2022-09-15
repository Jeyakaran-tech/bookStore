package cloudsql

import (
	"database/sql"
	"sync"
)

var (
	db   *sql.DB
	once sync.Once
)

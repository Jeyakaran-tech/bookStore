package types

import (
	"time"
)

type Books struct {
	ID             int       `json:"id" url:"id"`
	Author         string    `json:"author" url:"author"`
	Published_date string    `json:"publishedDate" url:"publishedDate"`
	Price          float64   `json:"price" url:"price"`
	InStock        bool      `json:"inStock" url:"inStock"`
	Time_added     time.Time `json:"timeAdded" url:"timeAdded"`
}

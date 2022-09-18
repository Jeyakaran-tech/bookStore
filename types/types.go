package types

import (
	"time"
)

type Books struct {
	ID             int       `json:"id" url:"id"`
	Name           string    `json:"name" url:"name"`
	Author         string    `json:"author" url:"author"`
	Published_date string    `json:"publishedDate" url:"publishedDate"`
	Price          float64   `json:"price" url:"price"`
	InStock        bool      `json:"inStock" url:"inStock"`
	Time_added     time.Time `json:"timeAdded" url:"timeAdded"`
}

type Status struct {
	Code        string `json:"code" url:"code"`
	Description string `json:"description" url:"description"`
}

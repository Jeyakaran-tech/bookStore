package types

import (
	"time"
)

type Books struct {
	ID             int
	Author         string
	Published_date string
	Price          float64
	InStock        bool
	Time_added     time.Time
}

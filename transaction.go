package mf

import "time"

type TransactionHistory struct {
	Content string
	Amount  int
	Date    time.Time
}

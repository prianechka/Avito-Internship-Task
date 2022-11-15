package transaction

import "time"

type Transaction struct {
	TransactionID   int
	UserID          int
	TransactionType int
	Sum             float64
	Time            time.Time
	ActionComments  string
	AddComments     string
}

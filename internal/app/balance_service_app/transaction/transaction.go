package transaction

import "time"

type Transaction struct {
	TransactionID   int64
	UserID          int64
	TransactionType int64
	Sum             float64
	Time            time.Time
	ActionComments  string
	AddComments     string
}

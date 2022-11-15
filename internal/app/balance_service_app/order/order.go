package order

import "time"

type Order struct {
	OrderID      int
	UserID       int
	ServiceID    int
	OrderCost    float64
	CreatingTime time.Time
	Comment      string
	OrderState   int
}

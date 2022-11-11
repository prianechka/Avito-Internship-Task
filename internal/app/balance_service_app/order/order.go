package order

import "time"

type Order struct {
	OrderID      int64
	UserID       int64
	ServiceType  int64
	OrderCost    float64
	CreatingTime time.Time
	Comment      string
	OrderState   int64
}

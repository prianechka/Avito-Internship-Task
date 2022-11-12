package messages

type AcceptServiceMessage struct {
	UserID    int64 `json:"user_id"`
	OrderID   int64 `json:"order_id"`
	ServiceID int64 `json:"service_id"`
}

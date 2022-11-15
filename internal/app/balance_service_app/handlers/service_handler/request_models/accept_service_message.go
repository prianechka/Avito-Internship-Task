package request_models

type AcceptServiceMessage struct {
	UserID    int `json:"user_id"`
	OrderID   int `json:"order_id"`
	ServiceID int `json:"service_id"`
}

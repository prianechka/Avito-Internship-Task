package request_models

type BuyServiceMessage struct {
	UserID    int     `json:"user_id"`
	OrderID   int     `json:"order_id"`
	ServiceID int     `json:"service_id"`
	Sum       float64 `json:"sum"`
	Comment   string  `json:"comment"`
}

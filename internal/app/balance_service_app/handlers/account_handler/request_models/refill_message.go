package request_models

type RefillMessage struct {
	UserID  int64   `json:"user_id"`
	Sum     float64 `json:"sum"`
	Comment string  `json:"comment"`
}

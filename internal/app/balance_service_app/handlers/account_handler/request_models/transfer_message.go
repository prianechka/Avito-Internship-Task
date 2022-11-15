package request_models

type TransferMessage struct {
	SrcUserID int     `json:"src_user_id"`
	DstUserID int     `json:"dst_user_id"`
	Sum       float64 `json:"sum"`
	Comment   string  `json:"comment"`
}

package messages

type TransferMessage struct {
	SrcUserID int64   `json:"src_user_id"`
	DstUserID int64   `json:"dst_user_id"`
	Sum       float64 `json:"sum"`
	Comment   string  `json:"comment"`
}

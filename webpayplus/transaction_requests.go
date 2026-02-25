package webpayplus

type TransactionCreateRequest struct {
	BuyOrder  string  `json:"buy_order"`
	SessionId string  `json:"session_id"`
	Amount    float64 `json:"amount"`
	ReturnUrl string  `json:"return_url"`
}

type TransactionRefundRequest struct {
	Amount float64 `json:"amount"`
}

type TransactionCaptureRequest struct {
	BuyOrder          string  `json:"buy_order"`
	AuthorizationCode string  `json:"authorization_code"`
	CaptureAmount     float64 `json:"capture_amount"`
}

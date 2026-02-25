package webpayplus

type MallTransactionCreateRequest struct {
	BuyOrder  string        `json:"buy_order"`
	SessionId string        `json:"session_id"`
	ReturnUrl string        `json:"return_url"`
	Details   []MallDetails `json:"details"`
}

type MallTransactionRefundRequest struct {
	BuyOrder     string  `json:"buy_order"`
	CommerceCode string  `json:"commerce_code"`
	Amout        float64 `json:"amount"`
}

type MallTransactionCaptureRequest struct {
	BuyOrder          string  `json:"buy_order"`
	CommerceCode      string  `json:"commerce_code"`
	AuthorizationCode string  `json:"authorization_code"`
	CaptureAmount     float64 `json:"capture_amount"`
}

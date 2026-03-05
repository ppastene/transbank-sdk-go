package webpayplus

type TransactionCreateResponse struct {
	Token string `json:"token"`
	Url   string `json:"url"`
}

type TransactionStatusResponse struct {
	Vci                string                 `json:"vci"`
	Amount             float64                `json:"amount"`
	Status             string                 `json:"status"`
	BuyOrder           string                 `json:"buy_order"`
	SessionId          string                 `json:"session_id"`
	CardDetail         TransactionCardDetails `json:"card_detail"`
	AccountingDate     string                 `json:"accounting_date"`
	TransactionDate    string                 `json:"transaction_date"`
	AuthorizationCode  string                 `json:"authorization_code"`
	PaymentTypeCode    string                 `json:"payment_type_code"`
	ResponseCode       int                    `json:"response_code"`
	InstallmentsAmount int                    `json:"installments_amount"`
	InstallmentsNumber int                    `json:"installments_number"`
	Balance            float64                `json:"balance"`
}

type TransactionCardDetails struct {
	CardNumber string `json:"card_number"`
}

type TransactionCommitResponse struct {
	TransactionStatusResponse
}

type TransactionRefundResponse struct {
	Type              string  `json:"type"`
	AuthorizationCode string  `json:"authorization_code"`
	AuthorizationDate string  `json:"authorization_date"`
	NullifiedAmount   float64 `json:"nullified_amount"`
	Balance           float64 `json:"balance"`
	ResponseCode      int     `json:"response_code"`
}

type TransactionCaptureResponse struct {
	AuthorizationCode string  `json:"authorization_code"`
	AuthorizationDate string  `json:"authorization_date"`
	CapturedAmount    float64 `json:"captured_amount"`
	ResponseCode      int     `json:"response_code"`
}

func (t *TransactionStatusResponse) IsApproved() bool {
	if t.ResponseCode != 0 {
		return false
	}
	switch t.Status {
	case "CAPTURED", "REVERSED", "NULLIFIED", "AUTHORIZED", "PARTIALLY_NULLIFIED":
		return true
	default:
		return false
	}
}

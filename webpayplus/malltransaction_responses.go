package webpayplus

type MallTransactionCreateResponse struct {
	Token string `json:"token"`
	Url   string `json:"url"`
}

type MallTransactionStatusResponse struct {
	BuyOrder        string                           `json:"buy_order"`
	SessionId       string                           `json:"session_id"`
	CardDetail      MallTransactionCardDetails       `json:"card_detail"`
	AccountingDate  string                           `json:"accounting_date"`
	TransactionDate string                           `json:"transaction_date"`
	Vci             string                           `json:"vci"`
	Details         []MallTransactionDetailsResponse `json:"details"`
}

type MallTransactionCardDetails struct {
	CardNumber string `json:"card_number"`
}

type MallTransactionDetailsResponse struct {
	AuthorizationCode  string  `json:"authorization_code"`
	PaymentTypeCode    string  `json:"payment_type_code"`
	ResponseCode       int     `json:"response_code"`
	Amount             float64 `json:"amount"`
	InstallmentsAmount int     `json:"installments_amount"`
	InstallmentsNumber int     `json:"installments_number"`
	CommerceCode       string  `json:"commerce_code"`
	BuyOrder           string  `json:"buy_order"`
	Status             string  `json:"status"`
	Balance            float64 `json:"balance"`
}

type MallTransactionCommitResponse struct {
	MallTransactionStatusResponse
}

type MallTransactionRefundResponse struct {
	Type              string
	AuthorizationCode string  `json:"authorization_code"`
	AuthorizationDate string  `json:"authorization_date"`
	Balance           float64 `json:"balance"`
	NullifiedAmount   float64 `json:"nullified_amount"`
	ResponseCode      int     `json:"response_code"`
}

type MallTransactionCaptureResponse struct {
	Token             string  `json:"token"`
	AuthorizationCode string  `json:"authorization_code"`
	AuthorizationDate string  `json:"authorization_date"`
	CapturedAmount    float64 `json:"captured_amount"`
	ResponseCode      int     `json:"response_code"`
}

func (m *MallTransactionStatusResponse) IsApproved() bool {
	if len(m.Details) == 0 {
		return false
	}

	for _, detail := range m.Details {
		if detail.IsApproved() {
			return true
		}
	}
	return false
}

func (d *MallTransactionDetailsResponse) IsApproved() bool {
	if d.ResponseCode != 0 {
		return false
	}
	switch d.Status {
	case "CAPTURED", "REVERSED", "NULLIFIED", "AUTHORIZED", "PARTIALLY_NULLIFIED":
		return true
	default:
		return false
	}
}

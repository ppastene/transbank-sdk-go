package webpayplus

// MallTransactionCreateResponse is the result of a Create call.
type MallTransactionCreateResponse struct {
	Token string `json:"token"`
	Url   string `json:"url"`
}

// MallTransactionStatusResponse is the current state of a mall transaction,
// including one detail per store.
type MallTransactionStatusResponse struct {
	BuyOrder        string                           `json:"buy_order"`
	SessionId       string                           `json:"session_id"`
	CardDetail      MallTransactionCardDetails       `json:"card_detail"`
	AccountingDate  string                           `json:"accounting_date"`
	TransactionDate string                           `json:"transaction_date"`
	Vci             string                           `json:"vci"`
	Details         []MallTransactionDetailsResponse `json:"details"`
}

// MallTransactionCardDetails contains the masked card number used in the
// payment.
type MallTransactionCardDetails struct {
	CardNumber string `json:"card_number"`
}

// MallTransactionDetailsResponse is the state of a single store transaction
// within a mall transaction.
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

// MallTransactionCommitResponse is the result of a Commit call.
type MallTransactionCommitResponse struct {
	MallTransactionStatusResponse
}

// MallTransactionRefundResponse is the result of a Refund call.
type MallTransactionRefundResponse struct {
	Type              string  `json:"type"`
	AuthorizationCode string  `json:"authorization_code"`
	AuthorizationDate string  `json:"authorization_date"`
	Balance           float64 `json:"balance"`
	NullifiedAmount   float64 `json:"nullified_amount"`
	ResponseCode      int     `json:"response_code"`
}

// MallTransactionCaptureResponse is the result of a Capture call.
type MallTransactionCaptureResponse struct {
	AuthorizationCode string  `json:"authorization_code"`
	AuthorizationDate string  `json:"authorization_date"`
	CapturedAmount    float64 `json:"captured_amount"`
	ResponseCode      int     `json:"response_code"`
}

// IsApproved returns true when all transaction details are approved and there
// is at least one detail.
func (m *MallTransactionStatusResponse) IsApproved() bool {
	if len(m.Details) == 0 {
		return false
	}
	for _, detail := range m.Details {
		if !detail.IsApproved() {
			return false
		}
	}
	return true
}

// IsApproved returns true when the detail is AUTHORIZED or CAPTURED and its
// response code is 0.
func (d *MallTransactionDetailsResponse) IsApproved() bool {
	if d.ResponseCode != 0 {
		return false
	}
	switch d.Status {
	case "AUTHORIZED", "CAPTURED":
		return true
	default:
		return false
	}
}

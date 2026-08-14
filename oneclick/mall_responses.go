package oneclick

// OneclickMallInscriptionStartResponse is the result of a Start call.
type OneclickMallInscriptionStartResponse struct {
	Token     string `json:"token"`
	UrlWebpay string `json:"url_webpay"`
}

// OneclickMallInscriptionFinishResponse is the result of a Finish call.
type OneclickMallInscriptionFinishResponse struct {
	ResponseCode      int    `json:"response_code"`
	TbkUser           string `json:"tbk_user"`
	AuthorizationCode string `json:"authorization_code"`
	CardType          string `json:"card_type"`
	CardNumber        string `json:"card_number"`
}

// OneclickMallTransactionStatusResponse is the current state of a transaction,
// including one detail per store.
type OneclickMallTransactionStatusResponse struct {
	BuyOrder        string                             `json:"buy_order"`
	CardDetail      OneclickMallTransactionCardDetails `json:"card_detail"`
	AccountingDate  string                             `json:"accounting_date"`
	TransactionDate string                             `json:"transaction_date"`
	Details         []OneclickMallTransactionDetails   `json:"details"`
}

// OneclickMallTransactionCardDetails contains the masked card number used in
// the payment.
type OneclickMallTransactionCardDetails struct {
	CardNumber string `json:"card_number"`
}

// OneclickMallTransactionDetails is the state of a single store transaction
// within an OneClick Mall transaction.
type OneclickMallTransactionDetails struct {
	Amount             float64 `json:"amount"`
	Status             string  `json:"status"`
	AuthorizationCode  string  `json:"authorization_code"`
	PaymentTypeCode    string  `json:"payment_type_code"`
	ResponseCode       int     `json:"response_code"`
	InstallmentsAmount int     `json:"installments_amount"`
	InstallmentsNumber int     `json:"installments_number"`
	CommerceCode       string  `json:"commerce_code"`
	BuyOrder           string  `json:"buy_order"`
}

// OneclickMallTransactionAuthorizeResponse is the result of an Authorize call.
type OneclickMallTransactionAuthorizeResponse struct {
	OneclickMallTransactionStatusResponse
}

// OneclickMallTransactionRefundResponse is the result of a Refund call.
type OneclickMallTransactionRefundResponse struct {
	Type              string  `json:"type"`
	AuthorizationCode string  `json:"authorization_code"`
	AuthorizationDate string  `json:"authorization_date"`
	NullifiedAmount   float64 `json:"nullified_amount"`
	Balance           float64 `json:"balance"`
	ResponseCode      int     `json:"response_code"`
}

// OneclickMallTransactionCaptureResponse is the result of a Capture call.
type OneclickMallTransactionCaptureResponse struct {
	AuthorizationCode string  `json:"authorization_code"`
	AuthorizationDate string  `json:"authorization_date"`
	CapturedAmount    float64 `json:"captured_amount"`
	ResponseCode      int     `json:"response_code"`
}

// IsApproved returns true when the detail is AUTHORIZED or CAPTURED and its
// response code is 0.
func (d *OneclickMallTransactionDetails) IsApproved() bool {
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

// IsApproved returns true when all transaction details are approved and there
// is at least one detail.
func (m *OneclickMallTransactionStatusResponse) IsApproved() bool {
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

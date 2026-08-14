package transaccioncompleta

// MallTransactionCreateResponse is the result of a Create call.
type MallTransactionCreateResponse struct {
	Token string `json:"token"`
}

// MallTransactionInstallmentsResponse is the result of an Installments call,
// with one detail per store.
type MallTransactionInstallmentsResponse []MallTransactionInstallmentsDetail

// MallTransactionInstallmentsDetail is the installment data of a single store.
type MallTransactionInstallmentsDetail struct {
	InstallmentsAmount  int               `json:"installments_amount"`
	IdQueryInstallments int               `json:"id_query_installments"`
	DeferredPeriods     []DeferredPeriods `json:"deferred_periods"`
}

// MallTransactionStatusResponse is the current state of a mall transaction,
// including one detail per store.
type MallTransactionStatusResponse struct {
	BuyOrder        string                   `json:"buy_order"`
	CardDetail      TransactionCardDetails   `json:"card_detail"`
	AccountingDate  string                   `json:"accounting_date"`
	TransactionDate string                   `json:"transaction_date"`
	Details         []MallTransactionDetails `json:"details"`
	PrepaidBalance  int                      `json:"prepaid_balance"`
}

// MallTransactionDetails is the state of a single store transaction within a
// mall transaction.
type MallTransactionDetails struct {
	Amount             float64 `json:"amount"`
	Status             string  `json:"status"`
	AuthorizationCode  string  `json:"authorization_code"`
	PaymentTypeCode    string  `json:"payment_type_code"`
	ResponseCode       int     `json:"response_code"`
	InstallmentsAmount int     `json:"installments_amount"`
	InstallmentsNumber int     `json:"installments_number"`
	CommerceCode       string  `json:"commerce_code"`
	BuyOrder           string  `json:"buy_order"`
	Balance            float64 `json:"balance"`
}

// IsApproved returns true when the detail is AUTHORIZED or CAPTURED and its
// response code is 0.
func (d *MallTransactionDetails) IsApproved() bool {
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

// MallTransactionCommitResponse is the result of a Commit call.
type MallTransactionCommitResponse struct {
	MallTransactionStatusResponse
}

// MallTransactionRefundResponse is the result of a Refund call.
type MallTransactionRefundResponse struct {
	Type              string  `json:"type"`
	AuthorizationCode string  `json:"authorization_code"`
	AuthorizationDate string  `json:"authorization_date"`
	NullifiedAmount   float64 `json:"nullified_amount"`
	Balance           float64 `json:"balance"`
	ResponseCode      int     `json:"response_code"`
}

// MallTransactionCaptureResponse is the result of a Capture call.
type MallTransactionCaptureResponse struct {
	AuthorizationCode string  `json:"authorization_code"`
	AuthorizationDate string  `json:"authorization_date"`
	CapturedAmount    float64 `json:"captured_amount"`
	ResponseCode      int     `json:"response_code"`
}

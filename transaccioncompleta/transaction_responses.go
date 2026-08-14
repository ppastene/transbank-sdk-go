package transaccioncompleta

// TransactionCreateResponse is the result of a Create call.
type TransactionCreateResponse struct {
	Token string `json:"token"`
}

// TransactionInstallmentsResponse is the result of an Installments call.
type TransactionInstallmentsResponse struct {
	InstallmentsAmount  int               `json:"installments_amount"`
	IdQueryInstallments int               `json:"id_query_installments"`
	DeferredPeriods     []DeferredPeriods `json:"deferred_periods"`
}

// DeferredPeriods describes a deferred period of an installment payment.
type DeferredPeriods struct {
	Amount int `json:"amount"`
	Period int `json:"period"`
}

// TransactionStatusResponse is the current state of a transaction.
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
	PrepaidBalance     int                    `json:"prepaid_balance"`
}

// IsApproved returns true when the transaction is AUTHORIZED or CAPTURED and
// its response code is 0.
func (t *TransactionStatusResponse) IsApproved() bool {
	if t.ResponseCode != 0 {
		return false
	}
	switch t.Status {
	case "AUTHORIZED", "CAPTURED":
		return true
	default:
		return false
	}
}

// TransactionCardDetails contains the masked card number used in the payment.
type TransactionCardDetails struct {
	CardNumber string `json:"card_number"`
}

// TransactionCommitResponse is the result of a Commit call.
type TransactionCommitResponse struct {
	TransactionStatusResponse
}

// TransactionRefundResponse is the result of a Refund call.
type TransactionRefundResponse struct {
	Type              string  `json:"type"`
	AuthorizationCode string  `json:"authorization_code"`
	AuthorizationDate string  `json:"authorization_date"`
	NullifiedAmount   float64 `json:"nullified_amount"`
	Balance           float64 `json:"balance"`
	ResponseCode      int     `json:"response_code"`
}

// TransactionCaptureResponse is the result of a Capture call.
type TransactionCaptureResponse struct {
	AuthorizationCode string  `json:"authorization_code"`
	AuthorizationDate string  `json:"authorization_date"`
	CapturedAmount    float64 `json:"captured_amount"`
	ResponseCode      int     `json:"response_code"`
}

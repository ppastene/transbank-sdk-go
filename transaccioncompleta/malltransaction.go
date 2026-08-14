package transaccioncompleta

import (
	"fmt"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/internal"
)

// MallDetails describes a single store transaction inside a Full Transaction
// Mall transaction.
type MallDetails struct {
	Amount       float64 `json:"amount"`        // amount to be charged to the customer
	CommerceCode string  `json:"commerce_code"` // commerce code of the store
	BuyOrder     string  `json:"buy_order"`     // buy order of the store
}

// MallInstallmentsDetails describes the installment query for a single store
// inside a Full Transaction Mall transaction.
type MallInstallmentsDetails struct {
	CommerceCode       string `json:"commerce_code"`       // commerce code of the store
	BuyOrder           string `json:"buy_order"`           // buy order of the store
	InstallmentsNumber int    `json:"installments_number"` // number of installments
}

// MallCommitDetails describes the commit parameters for a single store inside a
// Full Transaction Mall transaction. IdQueryInstallments, DeferredPeriodIndex
// and GracePeriod are optional and only sent when not nil.
type MallCommitDetails struct {
	CommerceCode        string `json:"commerce_code"`                   // commerce code of the store
	BuyOrder            string `json:"buy_order"`                       // buy order of the store
	IdQueryInstallments *int   `json:"id_query_installments,omitempty"` // installment query id
	DeferredPeriodIndex *int   `json:"deferred_period_index,omitempty"` // deferred period index
	GracePeriod         *bool  `json:"grace_period,omitempty"`          // apply grace period
}

// MallTransaction provides access to the Full Transaction Mall API. Create an
// instance with NewMallTransaction.
type MallTransaction struct {
	config internal.Config
}

// NewMallTransaction returns a MallTransaction for the given options. It
// validates the options and the commerce code and returns a
// *transbank.ValidationError on failure.
func NewMallTransaction(opts transbank.Options) (*MallTransaction, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	if err := internal.ValidateCommerceCode(opts.CommerceCode); err != nil {
		return nil, err
	}
	baseURL := internal.INTEGRATION_URL
	if opts.Environment == transbank.Production {
		baseURL = internal.PRODUCTION_URL
	}
	cfg := internal.NewConfig(opts.CommerceCode, opts.ApiKey, baseURL)
	cfg.Headers = map[string]string{
		"Tbk-Api-Key-Id":     opts.CommerceCode,
		"Tbk-Api-Key-Secret": opts.ApiKey,
	}
	if opts.HTTPClient != nil {
		cfg.HTTP = opts.HTTPClient
	}
	return &MallTransaction{config: cfg}, nil
}

// Create starts a new mall transaction for the given order and session, the
// customer card data and one detail per store. The cvv is optional: pass "" to
// omit it when the merchant has the "without cvv" option enabled. It returns
// the transaction token.
func (m *MallTransaction) Create(buyOrder, sessionId, cardNumber, cardExpirationDate string, details []MallDetails, cvv string) (*MallTransactionCreateResponse, error) {
	if err := internal.ValidateBuyOrder(buyOrder); err != nil {
		return nil, err
	}
	if err := internal.ValidateSessionID(sessionId); err != nil {
		return nil, err
	}
	if err := internal.ValidateCardNumber(cardNumber); err != nil {
		return nil, err
	}
	if err := internal.ValidateCardExpirationDate(cardExpirationDate); err != nil {
		return nil, err
	}
	if err := validateMallDetails(details); err != nil {
		return nil, err
	}
	if err := internal.ValidateCVV(cvv); err != nil {
		return nil, err
	}

	payload := map[string]any{
		"buy_order":            buyOrder,
		"session_id":           sessionId,
		"card_number":          cardNumber,
		"card_expiration_date": cardExpirationDate,
		"details":              details,
	}
	if cvv != "" {
		payload["cvv"] = cvv
	}
	var response MallTransactionCreateResponse
	if err := internal.NewRequestor(&m.config).Post(transactionsPath, payload, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Installments queries the installments available for a transaction identified
// by its token, with one query detail per store. It returns the deferred
// periods when available.
func (m *MallTransaction) Installments(token string, details []MallInstallmentsDetails) (*MallTransactionInstallmentsResponse, error) {
	if err := internal.ValidateToken(token); err != nil {
		return nil, err
	}
	if err := validateMallInstallmentsDetails(details); err != nil {
		return nil, err
	}
	var response MallTransactionInstallmentsResponse
	if err := internal.NewRequestor(&m.config).Post(fmt.Sprintf("%s/%s/installments", transactionsPath, token), details, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Commit confirms a previously created mall transaction identified by its
// token, authorizing the payment of all its details. Each detail's optional
// fields are only sent when not nil.
func (m *MallTransaction) Commit(token string, details []MallCommitDetails) (*MallTransactionCommitResponse, error) {
	if err := internal.ValidateToken(token); err != nil {
		return nil, err
	}
	if err := validateMallCommitDetails(details); err != nil {
		return nil, err
	}
	payload := map[string]any{
		"details": details,
	}
	var response MallTransactionCommitResponse
	if err := internal.NewRequestor(&m.config).Put(fmt.Sprintf("%s/%s", transactionsPath, token), payload, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Status returns the current state of a mall transaction identified by its
// token, including one detail per store.
func (m *MallTransaction) Status(token string) (*MallTransactionStatusResponse, error) {
	if err := internal.ValidateToken(token); err != nil {
		return nil, err
	}
	var response MallTransactionStatusResponse
	if err := internal.NewRequestor(&m.config).Get(fmt.Sprintf("%s/%s", transactionsPath, token), &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Refund refunds a paid detail of a mall transaction, identified by the
// transaction token, the store buy order and its commerce code, for the
// specified amount. The refund type is either "NULLIFY" or "REVERSED".
func (m *MallTransaction) Refund(token, buyOrder, commerceCode string, amount float64) (*MallTransactionRefundResponse, error) {
	if err := internal.ValidateToken(token); err != nil {
		return nil, err
	}
	if err := internal.ValidateBuyOrder(buyOrder); err != nil {
		return nil, err
	}
	if err := internal.ValidateCommerceCode(commerceCode); err != nil {
		return nil, err
	}
	if err := internal.ValidateAmount(amount); err != nil {
		return nil, err
	}
	payload := map[string]any{
		"buy_order":     buyOrder,
		"commerce_code": commerceCode,
		"amount":        amount,
	}
	var response MallTransactionRefundResponse
	if err := internal.NewRequestor(&m.config).Post(fmt.Sprintf("%s/%s/refunds", transactionsPath, token), payload, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Capture captures a detail of a mall transaction with deferred capture,
// identified by the transaction token, the store commerce code and buy order,
// and the authorization code obtained after Commit. Only available in
// environments with deferred capture enabled.
func (m *MallTransaction) Capture(token, commerceCode, buyOrder, authorizationCode string, captureAmount float64) (*MallTransactionCaptureResponse, error) {
	if err := internal.ValidateToken(token); err != nil {
		return nil, err
	}
	if err := internal.ValidateCommerceCode(commerceCode); err != nil {
		return nil, err
	}
	if err := internal.ValidateBuyOrder(buyOrder); err != nil {
		return nil, err
	}
	if err := internal.ValidateAmount(captureAmount); err != nil {
		return nil, err
	}
	if authorizationCode == "" {
		return nil, &transbank.ValidationError{Message: "authorization_code must not be empty"}
	}
	payload := map[string]any{
		"commerce_code":      commerceCode,
		"buy_order":          buyOrder,
		"authorization_code": authorizationCode,
		"capture_amount":     captureAmount,
	}
	var response MallTransactionCaptureResponse
	if err := internal.NewRequestor(&m.config).Put(fmt.Sprintf("%s/%s/capture", transactionsPath, token), payload, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func validateMallDetails(details []MallDetails) error {
	if len(details) == 0 {
		return &transbank.ValidationError{Message: "details must not be empty"}
	}
	for _, d := range details {
		if err := internal.ValidateAmount(d.Amount); err != nil {
			return err
		}
		if err := internal.ValidateCommerceCode(d.CommerceCode); err != nil {
			return err
		}
		if err := internal.ValidateBuyOrder(d.BuyOrder); err != nil {
			return err
		}
	}
	return nil
}

func validateMallInstallmentsDetails(details []MallInstallmentsDetails) error {
	if len(details) == 0 {
		return &transbank.ValidationError{Message: "details must not be empty"}
	}
	for _, d := range details {
		if err := internal.ValidateCommerceCode(d.CommerceCode); err != nil {
			return err
		}
		if err := internal.ValidateBuyOrder(d.BuyOrder); err != nil {
			return err
		}
		if err := internal.ValidateInstallmentsNumber(d.InstallmentsNumber); err != nil {
			return err
		}
	}
	return nil
}

func validateMallCommitDetails(details []MallCommitDetails) error {
	if len(details) == 0 {
		return &transbank.ValidationError{Message: "details must not be empty"}
	}
	for _, d := range details {
		if err := internal.ValidateCommerceCode(d.CommerceCode); err != nil {
			return err
		}
		if err := internal.ValidateBuyOrder(d.BuyOrder); err != nil {
			return err
		}
	}
	return nil
}

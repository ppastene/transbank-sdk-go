// Package webpayplus provides an implementation of the Transbank Webpay Plus
// and Webpay Plus Mall REST APIs.
package webpayplus

import (
	"fmt"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/internal"
)

const pkg = "webpayplus"

const transactionsPath = "/rswebpaytransaction/api/webpay/v1.2/transactions"

// Transaction provides access to the Webpay Plus API for a single store.
// Create an instance with NewTransaction.
type Transaction struct {
	config internal.Config
}

// NewTransaction returns a Transaction for the given options. It validates the
// options and returns a *transbank.ValidationError on failure.
func NewTransaction(opts transbank.Options) (*Transaction, error) {
	if err := opts.Validate(); err != nil {
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
	return &Transaction{config: cfg}, nil
}

// Create starts a new transaction for the given order and session and the
// specified amount. It returns the transaction token and the URL of the payment
// form to which the customer must be redirected. The return URL is where
// Transbank redirects the customer after payment.
func (t *Transaction) Create(buyOrder, sessionId string, amount float64, returnUrl string) (*TransactionCreateResponse, error) {
	payload := map[string]any{
		"buy_order":  buyOrder,
		"session_id": sessionId,
		"amount":     amount,
		"return_url": returnUrl,
	}
	var response TransactionCreateResponse
	if err := internal.NewRequestor(&t.config).Post(transactionsPath, payload, &response); err != nil {
		return nil, fmt.Errorf("%s create: %w", pkg, err)
	}
	return &response, nil
}

// Commit confirms a previously created transaction identified by its token,
// authorizing the payment. It returns the authorization details.
func (t *Transaction) Commit(token string) (*TransactionCommitResponse, error) {
	if err := internal.ValidateURLParam("token", token, 64); err != nil {
		return nil, err
	}
	var response TransactionCommitResponse
	if err := internal.NewRequestor(&t.config).Put(fmt.Sprintf("%s/%s", transactionsPath, token), nil, &response); err != nil {
		return nil, fmt.Errorf("%s commit: %w", pkg, err)
	}
	return &response, nil
}

// Status returns the current state of a transaction identified by its token.
func (t *Transaction) Status(token string) (*TransactionStatusResponse, error) {
	if err := internal.ValidateURLParam("token", token, 64); err != nil {
		return nil, err
	}
	var response TransactionStatusResponse
	if err := internal.NewRequestor(&t.config).Get(fmt.Sprintf("%s/%s", transactionsPath, token), &response); err != nil {
		return nil, fmt.Errorf("%s status: %w", pkg, err)
	}
	return &response, nil
}

// Refund refunds a paid transaction identified by its token, totally or
// partially, for the specified amount. The refund type is either "NULLIFY" or
// "REVERSED".
func (t *Transaction) Refund(token string, amount float64) (*TransactionRefundResponse, error) {
	if err := internal.ValidateURLParam("token", token, 64); err != nil {
		return nil, err
	}
	payload := map[string]float64{
		"amount": amount,
	}
	var response TransactionRefundResponse
	if err := internal.NewRequestor(&t.config).Post(fmt.Sprintf("%s/%s/refunds", transactionsPath, token), payload, &response); err != nil {
		return nil, fmt.Errorf("%s refund: %w", pkg, err)
	}
	return &response, nil
}

// Capture captures a transaction with deferred capture identified by its token,
// using the buy order and authorization code obtained after Commit. Only
// available in environments with deferred capture enabled.
func (t *Transaction) Capture(token, buyOrder, authorizationCode string, captureAmount float64) (*TransactionCaptureResponse, error) {
	if err := internal.ValidateURLParam("token", token, 64); err != nil {
		return nil, err
	}
	payload := map[string]any{
		"buy_order":          buyOrder,
		"authorization_code": authorizationCode,
		"capture_amount":     captureAmount,
	}
	var response TransactionCaptureResponse
	if err := internal.NewRequestor(&t.config).Put(fmt.Sprintf("%s/%s/capture", transactionsPath, token), payload, &response); err != nil {
		return nil, fmt.Errorf("%s capture: %w", pkg, err)
	}
	return &response, nil
}

// Package transaccioncompleta provides an implementation of the Transbank Full
// Transaction (Transacción Completa) and Full Transaction Mall REST APIs.
package transaccioncompleta

import (
	"fmt"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/internal"
)

const transactionsPath = "/rswebpaytransaction/api/webpay/v1.2/transactions"

// Transaction provides access to the Full Transaction API for a single store.
// Create an instance with NewTransaction.
type Transaction struct {
	config internal.Config
}

// NewTransaction returns a Transaction for the given options. It validates the
// options and the commerce code and returns a *transbank.ValidationError on
// failure.
func NewTransaction(opts transbank.Options) (*Transaction, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	if opts.ValidateInputs {
		if err := internal.ValidateCommerceCode(opts.CommerceCode); err != nil {
			return nil, err
		}
	}
	baseURL := internal.INTEGRATION_URL
	if opts.Environment == transbank.Production {
		baseURL = internal.PRODUCTION_URL
	}
	cfg := internal.NewConfig(opts.CommerceCode, opts.ApiKey, baseURL, opts.ValidateInputs)
	cfg.Headers = map[string]string{
		"Tbk-Api-Key-Id":     opts.CommerceCode,
		"Tbk-Api-Key-Secret": opts.ApiKey,
	}
	if opts.HTTPClient != nil {
		cfg.HTTP = opts.HTTPClient
	}
	return &Transaction{config: cfg}, nil
}

// Create starts a new transaction for the given order and session, the
// specified amount and the customer card data. The cvv is optional: pass "" to
// omit it when the merchant has the "without cvv" option enabled. It returns
// the transaction token.
func (t *Transaction) Create(buyOrder, sessionId string, amount float64, cardNumber, cardExpirationDate, cvv string) (*TransactionCreateResponse, error) {
	if t.config.ValidateInputs {
		if err := internal.ValidateBuyOrder(buyOrder); err != nil {
			return nil, err
		}
		if err := internal.ValidateSessionID(sessionId); err != nil {
			return nil, err
		}
		if err := internal.ValidateAmount(amount); err != nil {
			return nil, err
		}
		if err := internal.ValidateCardNumber(cardNumber); err != nil {
			return nil, err
		}
		if err := internal.ValidateCardExpirationDate(cardExpirationDate); err != nil {
			return nil, err
		}
		if err := internal.ValidateCVV(cvv); err != nil {
			return nil, err
		}
	}

	payload := map[string]any{
		"buy_order":            buyOrder,
		"session_id":           sessionId,
		"amount":               amount,
		"card_number":          cardNumber,
		"card_expiration_date": cardExpirationDate,
	}
	if cvv != "" {
		payload["cvv"] = cvv
	}
	var response TransactionCreateResponse
	if err := internal.NewRequestor(&t.config).Post(transactionsPath, payload, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Installments queries the installments available for a transaction identified
// by its token, for the given number of installments. It returns the deferred
// periods when available.
func (t *Transaction) Installments(token string, installmentsNumber int) (*TransactionInstallmentsResponse, error) {
	if err := internal.ValidateToken(token); err != nil {
		return nil, err
	}
	if t.config.ValidateInputs {
		if err := internal.ValidateInstallmentsNumber(installmentsNumber); err != nil {
			return nil, err
		}
	}
	payload := map[string]any{
		"installments_number": installmentsNumber,
	}
	var response TransactionInstallmentsResponse
	if err := internal.NewRequestor(&t.config).Post(fmt.Sprintf("%s/%s/installments", transactionsPath, token), payload, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Commit confirms a previously created transaction identified by its token,
// authorizing the payment. The idQueryInstallments, deferredPeriodIndex and
// gracePeriod parameters are optional and only sent when not nil; pass nil for
// a single-installment payment.
func (t *Transaction) Commit(token string, idQueryInstallments *int, deferredPeriodIndex *int, gracePeriod *bool) (*TransactionCommitResponse, error) {
	if err := internal.ValidateToken(token); err != nil {
		return nil, err
	}
	payload := map[string]any{}
	if idQueryInstallments != nil {
		payload["id_query_installments"] = *idQueryInstallments
	}
	if deferredPeriodIndex != nil {
		payload["deferred_period_index"] = *deferredPeriodIndex
	}
	if gracePeriod != nil {
		payload["grace_period"] = *gracePeriod
	}
	var response TransactionCommitResponse
	if err := internal.NewRequestor(&t.config).Put(fmt.Sprintf("%s/%s", transactionsPath, token), payload, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Status returns the current state of a transaction identified by its token.
func (t *Transaction) Status(token string) (*TransactionStatusResponse, error) {
	if err := internal.ValidateToken(token); err != nil {
		return nil, err
	}
	var response TransactionStatusResponse
	if err := internal.NewRequestor(&t.config).Get(fmt.Sprintf("%s/%s", transactionsPath, token), &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Refund refunds a paid transaction identified by its token, totally or
// partially, for the specified amount. The refund type is either "NULLIFY" or
// "REVERSED".
func (t *Transaction) Refund(token string, amount float64) (*TransactionRefundResponse, error) {
	if err := internal.ValidateToken(token); err != nil {
		return nil, err
	}
	if t.config.ValidateInputs {
		if err := internal.ValidateAmount(amount); err != nil {
			return nil, err
		}
	}
	payload := map[string]float64{
		"amount": amount,
	}
	var response TransactionRefundResponse
	if err := internal.NewRequestor(&t.config).Post(fmt.Sprintf("%s/%s/refunds", transactionsPath, token), payload, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Capture captures a transaction with deferred capture identified by its token,
// using the buy order and authorization code obtained after Commit. Only
// available in environments with deferred capture enabled.
func (t *Transaction) Capture(token, buyOrder, authorizationCode string, captureAmount float64) (*TransactionCaptureResponse, error) {
	if err := internal.ValidateToken(token); err != nil {
		return nil, err
	}
	if t.config.ValidateInputs {
		if err := internal.ValidateBuyOrder(buyOrder); err != nil {
			return nil, err
		}
		if err := internal.ValidateAmount(captureAmount); err != nil {
			return nil, err
		}
		if authorizationCode == "" {
			return nil, &transbank.ValidationError{Message: "authorization_code must not be empty"}
		}
	}
	payload := map[string]any{
		"buy_order":          buyOrder,
		"authorization_code": authorizationCode,
		"capture_amount":     captureAmount,
	}
	var response TransactionCaptureResponse
	if err := internal.NewRequestor(&t.config).Put(fmt.Sprintf("%s/%s/capture", transactionsPath, token), payload, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

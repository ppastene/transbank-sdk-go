package webpayplus

import (
	"fmt"

	"github.com/ppastene/transbank-sdk-go/internal/httpclient"
	"github.com/ppastene/transbank-sdk-go/internal/shared"
)

type Transaction struct {
	requestor *shared.Requestor
}

func NewTransaction(client shared.HTTPClientInterface, options *shared.Options) *Transaction {
	if client == nil {
		client = httpclient.NewDefaultClient()
	}
	return &Transaction{
		&shared.Requestor{
			Client:  client,
			Options: options,
		},
	}
}

func (t *Transaction) Create(buyOrder, sessionId string, amount float64, returnUrl string) (*TransactionCreateResponse, error) {
	if err := shared.HasTextWithMaxLength(buyOrder, 26, "buyOrder"); err != nil {
		return nil, &shared.WebpayError{Code: -1, ServiceMessage: "SDK Validation Error", Cause: err}
	}
	if err := shared.HasTextWithMaxLength(sessionId, 61, "sessionId"); err != nil {
		return nil, &shared.WebpayError{Code: -1, ServiceMessage: "SDK Validation Error", Cause: err}
	}
	if err := shared.HasTextWithMaxLength(returnUrl, 256, "returnUrl"); err != nil {
		return nil, &shared.WebpayError{Code: -1, ServiceMessage: "SDK Validation Error", Cause: err}
	}
	payload := map[string]any{
		"buy_order":  buyOrder,
		"session_id": sessionId,
		"amount":     amount,
		"return_url": returnUrl,
	}
	var response TransactionCreateResponse
	_, err := t.requestor.Do("POST", "/rswebpaytransaction/api/webpay/v1.2/transactions", payload, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (t *Transaction) Commit(token string) (*TransactionCommitResponse, error) {
	if err := shared.HasTextWithMaxLength(token, 64, "token"); err != nil {
		return nil, &shared.WebpayError{Code: -1, ServiceMessage: "SDK Validation Error", Cause: err}
	}
	var response TransactionCommitResponse

	_, err := t.requestor.Do("PUT", fmt.Sprintf("/rswebpaytransaction/api/webpay/v1.2/transactions/%s", token), nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (t *Transaction) Status(token string) (*TransactionStatusResponse, error) {
	if err := shared.HasTextWithMaxLength(token, 64, "token"); err != nil {
		return nil, &shared.WebpayError{Code: -1, ServiceMessage: "SDK Validation Error", Cause: err}
	}
	var response TransactionStatusResponse

	_, err := t.requestor.Do("GET", fmt.Sprintf("/rswebpaytransaction/api/webpay/v1.2/transactions/%s", token), nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (t *Transaction) Refund(token string, amount float64) (*TransactionRefundResponse, error) {
	if err := shared.HasTextWithMaxLength(token, 64, "token"); err != nil {
		return nil, &shared.WebpayError{Code: -1, ServiceMessage: "SDK Validation Error", Cause: err}
	}

	payload := map[string]float64{
		"amount": amount,
	}

	var response TransactionRefundResponse

	_, err := t.requestor.Do("POST", fmt.Sprintf("/rswebpaytransaction/api/webpay/v1.2/transactions/%s/refunds", token), payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (t *Transaction) Capture(token, buyOrder, authorizationCode string, captureAmount float64) (*TransactionCaptureResponse, error) {
	if err := shared.HasTextWithMaxLength(token, 64, "token"); err != nil {
		return nil, &shared.WebpayError{Code: -1, ServiceMessage: "SDK Validation Error", Cause: err}
	}
	if err := shared.HasTextWithMaxLength(buyOrder, 26, "buyOrder"); err != nil {
		return nil, &shared.WebpayError{Code: -1, ServiceMessage: "SDK Validation Error", Cause: err}
	}
	if err := shared.HasTextWithMaxLength(authorizationCode, 6, "authorizationCode"); err != nil {
		return nil, &shared.WebpayError{Code: -1, ServiceMessage: "SDK Validation Error", Cause: err}
	}

	payload := map[string]any{
		"buy_order":          buyOrder,
		"authorization_code": authorizationCode,
		"capture_amount":     captureAmount,
	}

	var response TransactionCaptureResponse

	_, err := t.requestor.Do("PUT", fmt.Sprintf("/rswebpaytransaction/api/webpay/v1.2/transactions/%s/capture", token), payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

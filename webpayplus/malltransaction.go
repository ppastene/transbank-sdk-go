package webpayplus

import (
	"fmt"

	httpclient "github.com/ppastene/transbank-sdk-go/internal/httpclient"
	"github.com/ppastene/transbank-sdk-go/internal/shared"
)

type MallDetails struct {
	Amount       float64 `json:"amount"`
	CommerceCode string  `json:"commerce_code"`
	BuyOrder     string  `json:"buy_order"`
}

type MallTransaction struct {
	requestor *shared.Requestor
}

func NewMallTransaction(client shared.HTTPClientInterface, options *shared.Options) *MallTransaction {
	if client == nil {
		client = httpclient.NewDefaultClient()
	}
	return &MallTransaction{
		&shared.Requestor{
			Client:  client,
			Options: options,
		},
	}
}

func (m *MallTransaction) Create(buyOrder, sessionId, returnUrl string, details []MallDetails) (*MallTransactionCreateResponse, error) {

	payload := map[string]any{
		"buy_order":  buyOrder,
		"session_id": sessionId,
		"return_url": returnUrl,
		"details":    details,
	}

	var response MallTransactionCreateResponse

	_, err := m.requestor.Do("POST", "/rswebpaytransaction/api/webpay/v1.2/transactions", payload, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (m *MallTransaction) Commit(token string) (*MallTransactionCommitResponse, error) {
	if err := shared.HasTextWithMaxLength(token, 64, "token"); err != nil {
		return nil, &shared.WebpayError{Code: -1, ServiceMessage: "SDK Error", Cause: err}
	}

	var response MallTransactionCommitResponse
	_, err := m.requestor.Do("PUT", fmt.Sprintf("/rswebpaytransaction/api/webpay/v1.2/transactions/%s", token), nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (m *MallTransaction) Status(token string) (*MallTransactionStatusResponse, error) {
	if err := shared.HasTextWithMaxLength(token, 64, "token"); err != nil {
		return nil, &shared.WebpayError{Code: -1, ServiceMessage: "SDK Error", Cause: err}
	}

	var response MallTransactionStatusResponse

	_, err := m.requestor.Do("GET", fmt.Sprintf("/rswebpaytransaction/api/webpay/v1.2/transactions/%s", token), nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (m *MallTransaction) Refund(token, buyOrder, childCommerceCode string, amount float64) (*MallTransactionRefundResponse, error) {
	if err := shared.HasTextWithMaxLength(token, 64, "token"); err != nil {
		return nil, &shared.WebpayError{Code: -1, ServiceMessage: "SDK Error", Cause: err}
	}

	payload := map[string]any{
		"buy_order":     buyOrder,
		"commerce_code": childCommerceCode,
		"amount":        amount,
	}

	var response MallTransactionRefundResponse

	_, err := m.requestor.Do("POST", fmt.Sprintf("/rswebpaytransaction/api/webpay/v1.2/transactions/%s/refunds", token), payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (m *MallTransaction) Capture(token, childCommerceCode, buyOrder, authorizationCode string, captureAmount float64) (*MallTransactionCaptureResponse, error) {
	if err := shared.HasTextWithMaxLength(token, 64, "token"); err != nil {
		return nil, &shared.WebpayError{Code: -1, ServiceMessage: "SDK Error", Cause: err}
	}

	payload := map[string]any{
		"buy_order":          buyOrder,
		"commerce_code":      childCommerceCode,
		"authorization_code": authorizationCode,
		"capture_amount":     captureAmount,
	}

	var response MallTransactionCaptureResponse

	_, err := m.requestor.Do("PUT", fmt.Sprintf("/rswebpaytransaction/api/webpay/v1.2/transactions/%s/capture", token), payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

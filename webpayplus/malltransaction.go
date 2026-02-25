package webpayplus

import (
	"errors"
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
	payload := MallTransactionCreateRequest{
		BuyOrder:  buyOrder,
		SessionId: sessionId,
		ReturnUrl: returnUrl,
		Details:   details,
	}

	var response MallTransactionCreateResponse

	code, err := m.requestor.Do("POST", "/rswebpaytransaction/api/webpay/v1.2/transactions", payload, &response)
	if err != nil {
		if code >= 0 {
			return nil, NewMallTransactionCreateException(err.Error(), code)
		}
		return nil, err
	}
	return &response, nil
}

func (m *MallTransaction) Commit(token string) (*MallTransactionCommitResponse, error) {
	if len(token) == 0 {
		return nil, errors.New("Token parameter given is empty.")
	}

	var response MallTransactionCommitResponse
	code, err := m.requestor.Do("PUT", fmt.Sprintf("/rswebpaytransaction/api/webpay/v1.2/transactions/%s", token), nil, &response)
	if err != nil {
		if code >= 0 {
			return nil, NewMallTransactionCommitException(err.Error(), code)
		}
		return nil, err
	}

	return &response, nil
}

func (m *MallTransaction) Status(token string) (*MallTransactionStatusResponse, error) {
	var response MallTransactionStatusResponse

	code, err := m.requestor.Do("GET", fmt.Sprintf("/rswebpaytransaction/api/webpay/v1.2/transactions/%s", token), nil, &response)
	if err != nil {
		if code >= 0 {
			return nil, NewMallTransactionStatusException(err.Error(), code)
		}
		return nil, err
	}

	return &response, nil
}

func (m *MallTransaction) Refund(token, buyOrder, childCommerceCode string, amount float64) (*MallTransactionRefundResponse, error) {
	payload := MallTransactionRefundRequest{
		BuyOrder:     buyOrder,
		CommerceCode: childCommerceCode,
		Amout:        amount,
	}

	var response MallTransactionRefundResponse

	code, err := m.requestor.Do("POST", fmt.Sprintf("/rswebpaytransaction/api/webpay/v1.2/transactions/%s/refunds", token), payload, &response)
	if err != nil {
		if code >= 0 {
			return nil, NewMallTransactionRefundException(err.Error(), code)
		}
		return nil, err
	}

	return &response, nil
}

func (m *MallTransaction) Capture(childCommerceCode, token, buyOrder, authorizationCode string, captureAmount float64) (*MallTransactionCaptureResponse, error) {
	payload := MallTransactionCaptureRequest{
		BuyOrder:          buyOrder,
		CommerceCode:      childCommerceCode,
		AuthorizationCode: authorizationCode,
		CaptureAmount:     captureAmount,
	}

	var response MallTransactionCaptureResponse

	code, err := m.requestor.Do("PUT", fmt.Sprintf("/rswebpaytransaction/api/webpay/v1.2/transactions/%s/capture", token), payload, &response)
	if err != nil {
		if code >= 0 {
			return nil, NewMallTransactionCaptureException(err.Error(), code)
		}
		return nil, err
	}

	return &response, nil
}

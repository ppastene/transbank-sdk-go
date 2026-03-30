package oneclick

import (
	"fmt"

	"github.com/ppastene/transbank-sdk-go/internal/httpclient"
	"github.com/ppastene/transbank-sdk-go/internal/shared"
)

type MallDetails struct {
	CommerceCode       string  `json:"commerce_code"`
	BuyOrder           string  `json:"buy_order"`
	Amount             float64 `json:"amount"`
	InstallmentsNumber int     `json:"installments_number"`
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

func (m *MallTransaction) Authorize(username, tbkUser, buyOrder string, details []MallDetails) (*OneclickMallTransactionAuthorizeResponse, error) {
	payload := map[string]any{
		"username":  username,
		"tbk_user":  tbkUser,
		"buy_order": buyOrder,
		"details":   details,
	}
	var response OneclickMallTransactionAuthorizeResponse
	_, err := m.requestor.Do("POST", "/rswebpaytransaction/api/oneclick/v1.2/transactions", payload, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (m *MallTransaction) Status(buyOrder string) (*OneclickMallTransactionStatusResponse, error) {
	var response OneclickMallTransactionStatusResponse
	_, err := m.requestor.Do("GET", fmt.Sprintf("/rswebpaytransaction/api/oneclick/v1.2/transactions/%s", buyOrder), nil, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (m *MallTransaction) Refund(buyOrder, childCommerceCode, childBuyOrder string, amount float64) (*OneclickMallTransactionRefundResponse, error) {
	payload := map[string]any{
		"commerce_code":    childCommerceCode,
		"detail_buy_order": childBuyOrder,
		"amount":           amount,
	}
	var response OneclickMallTransactionRefundResponse
	_, err := m.requestor.Do("POST", fmt.Sprintf("/rswebpaytransaction/api/oneclick/v1.2/transactions/%s/refunds", buyOrder), payload, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (m *MallTransaction) Capture(buyOrder, commerceCode, authorizationCode string, amount float64) (*OneclickMallTransactionCaptureResponse, error) {
	payload := map[string]any{
		"commerce_code":      commerceCode,
		"buy_order":          buyOrder,
		"capture_amount":     amount,
		"authorization_code": authorizationCode,
	}
	var response OneclickMallTransactionCaptureResponse
	_, err := m.requestor.Do("POST", "/rswebpaytransaction/api/oneclick/v1.2/transactions/capture", payload, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

package oneclick

import (
	"fmt"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/internal"
)

// MallDetails describes a single store payment inside an OneClick Mall
// transaction.
type MallDetails struct {
	CommerceCode       string  `json:"commerce_code"`       // commerce code of the store
	BuyOrder           string  `json:"buy_order"`           // buy order of the store
	Amount             float64 `json:"amount"`              // amount to be charged to the customer
	InstallmentsNumber int     `json:"installments_number"` // number of installments, between 0 and 99
}

// MallTransaction provides access to the OneClick Mall payment API using
// enrolled cards. Create an instance with NewMallTransaction.
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
	return &MallTransaction{config: cfg}, nil
}

// Authorize charges a payment using the card enrolled for the given user and
// tbk_user, with one detail per store. It returns the authorization details.
func (m *MallTransaction) Authorize(username, tbkUser, buyOrder string, details []MallDetails) (*OneclickMallTransactionAuthorizeResponse, error) {
	if m.config.ValidateInputs {
		if err := internal.ValidateUsername(username); err != nil {
			return nil, err
		}
		if err := internal.ValidateTbkUser(tbkUser); err != nil {
			return nil, err
		}
		if err := internal.ValidateBuyOrder(buyOrder); err != nil {
			return nil, err
		}
		if err := validateMallDetails(details, m.config.ValidateInputs); err != nil {
			return nil, err
		}
	}

	payload := map[string]any{
		"username":  username,
		"tbk_user":  tbkUser,
		"buy_order": buyOrder,
		"details":   details,
	}
	var response OneclickMallTransactionAuthorizeResponse
	if err := internal.NewRequestor(&m.config).Post(oneClickPath+"/transactions", payload, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func validateMallDetails(details []MallDetails, validateInputs bool) error {
	if !validateInputs {
		return nil
	}
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
		if err := internal.ValidateAmount(d.Amount); err != nil {
			return err
		}
		if d.InstallmentsNumber < 0 || d.InstallmentsNumber > 99 {
			return &transbank.ValidationError{Message: "installments_number must be between 0 and 99"}
		}
	}
	return nil
}

// Status returns the current state of a transaction identified by its buy
// order, including one detail per store.
func (m *MallTransaction) Status(buyOrder string) (*OneclickMallTransactionStatusResponse, error) {
	if err := internal.ValidateBuyOrder(buyOrder); err != nil {
		return nil, err
	}
	var response OneclickMallTransactionStatusResponse
	if err := internal.NewRequestor(&m.config).Get(fmt.Sprintf("%s/transactions/%s", oneClickPath, buyOrder), &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Refund refunds a paid detail of a transaction, identified by the parent buy
// order, the store commerce code and its buy order, for the specified amount.
// The refund type is either "NULLIFY" or "REVERSED".
func (m *MallTransaction) Refund(buyOrder, childCommerceCode, childBuyOrder string, amount float64) (*OneclickMallTransactionRefundResponse, error) {
	if err := internal.ValidateBuyOrder(buyOrder); err != nil {
		return nil, err
	}
	if m.config.ValidateInputs {
		if err := internal.ValidateCommerceCode(childCommerceCode); err != nil {
			return nil, err
		}
		if err := internal.ValidateBuyOrder(childBuyOrder); err != nil {
			return nil, err
		}
		if err := internal.ValidateAmount(amount); err != nil {
			return nil, err
		}
	}

	payload := map[string]any{
		"commerce_code":    childCommerceCode,
		"detail_buy_order": childBuyOrder,
		"amount":           amount,
	}
	var response OneclickMallTransactionRefundResponse
	if err := internal.NewRequestor(&m.config).Post(fmt.Sprintf("%s/transactions/%s/refunds", oneClickPath, buyOrder), payload, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Capture captures a transaction with deferred capture, identified by the buy
// order, the store commerce code and the authorization code returned by
// Authorize, for the specified amount. Only available in environments with
// deferred capture enabled.
func (m *MallTransaction) Capture(buyOrder, commerceCode, authorizationCode string, amount float64) (*OneclickMallTransactionCaptureResponse, error) {
	if m.config.ValidateInputs {
		if err := internal.ValidateBuyOrder(buyOrder); err != nil {
			return nil, err
		}
		if err := internal.ValidateCommerceCode(commerceCode); err != nil {
			return nil, err
		}
		if err := internal.ValidateAmount(amount); err != nil {
			return nil, err
		}
		if authorizationCode == "" {
			return nil, &transbank.ValidationError{Message: "authorization_code must not be empty"}
		}
	}

	payload := map[string]any{
		"commerce_code":      commerceCode,
		"buy_order":          buyOrder,
		"capture_amount":     amount,
		"authorization_code": authorizationCode,
	}
	var response OneclickMallTransactionCaptureResponse
	if err := internal.NewRequestor(&m.config).Put(oneClickPath+"/transactions/capture", payload, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

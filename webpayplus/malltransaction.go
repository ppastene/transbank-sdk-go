package webpayplus

import (
	"fmt"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/internal"
)

// MallDetails describes a single store transaction inside a Webpay Plus Mall
// transaction.
type MallDetails struct {
	Amount       float64 `json:"amount"`        // amount to be charged to the customer
	CommerceCode string  `json:"commerce_code"` // commerce code of the store
	BuyOrder     string  `json:"buy_order"`     // buy order of the store
}

// MallTransaction provides access to the Webpay Plus Mall API. Create an
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

// Create starts a new mall transaction for the given order and session, with
// one detail per store, and the return URL where Transbank redirects the
// customer after payment. It returns the transaction token and the URL of the
// payment form.
func (m *MallTransaction) Create(buyOrder, sessionId, returnUrl string, details []MallDetails) (*MallTransactionCreateResponse, error) {
	if m.config.ValidateInputs {
		if err := internal.ValidateBuyOrder(buyOrder); err != nil {
			return nil, err
		}
		if err := internal.ValidateSessionID(sessionId); err != nil {
			return nil, err
		}
		if err := internal.ValidateReturnURL(returnUrl); err != nil {
			return nil, err
		}
		if err := validateMallDetails(details, m.config.ValidateInputs); err != nil {
			return nil, err
		}
	}

	payload := map[string]any{
		"buy_order":  buyOrder,
		"session_id": sessionId,
		"return_url": returnUrl,
		"details":    details,
	}
	var response MallTransactionCreateResponse
	if err := internal.NewRequestor(&m.config).Post(transactionsPath, payload, &response); err != nil {
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

// Commit confirms a previously created mall transaction identified by its
// token, authorizing the payment of all its details. It returns the
// authorization details.
func (m *MallTransaction) Commit(token string) (*MallTransactionCommitResponse, error) {
	if m.config.ValidateInputs {
		if err := internal.ValidateToken(token); err != nil {
			return nil, err
		}
	}
	var response MallTransactionCommitResponse
	if err := internal.NewRequestor(&m.config).Put(fmt.Sprintf("%s/%s", transactionsPath, token), nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Status returns the current state of a mall transaction identified by its
// token, including one detail per store.
func (m *MallTransaction) Status(token string) (*MallTransactionStatusResponse, error) {
	if m.config.ValidateInputs {
		if err := internal.ValidateToken(token); err != nil {
			return nil, err
		}
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
func (m *MallTransaction) Refund(token, buyOrder, childCommerceCode string, amount float64) (*MallTransactionRefundResponse, error) {
	if m.config.ValidateInputs {
		if err := internal.ValidateToken(token); err != nil {
			return nil, err
		}
		if err := internal.ValidateBuyOrder(buyOrder); err != nil {
			return nil, err
		}
		if err := internal.ValidateCommerceCode(childCommerceCode); err != nil {
			return nil, err
		}
		if err := internal.ValidateAmount(amount); err != nil {
			return nil, err
		}
	}

	payload := map[string]any{
		"buy_order":     buyOrder,
		"commerce_code": childCommerceCode,
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
func (m *MallTransaction) Capture(token, childCommerceCode, buyOrder, authorizationCode string, captureAmount float64) (*MallTransactionCaptureResponse, error) {
	if m.config.ValidateInputs {
		if err := internal.ValidateToken(token); err != nil {
			return nil, err
		}
		if err := internal.ValidateCommerceCode(childCommerceCode); err != nil {
			return nil, err
		}
		if err := internal.ValidateBuyOrder(buyOrder); err != nil {
			return nil, err
		}
		if err := internal.ValidateAmount(captureAmount); err != nil {
			return nil, err
		}
	}

	payload := map[string]any{
		"buy_order":          buyOrder,
		"commerce_code":      childCommerceCode,
		"authorization_code": authorizationCode,
		"capture_amount":     captureAmount,
	}
	var response MallTransactionCaptureResponse
	if err := internal.NewRequestor(&m.config).Put(fmt.Sprintf("%s/%s/capture", transactionsPath, token), payload, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

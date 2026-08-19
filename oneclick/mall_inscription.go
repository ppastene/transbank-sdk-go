// Package oneclick provides an implementation of the Transbank OneClick Mall
// REST API.
package oneclick

import (
	"fmt"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/internal"
)

const oneClickPath = "/rswebpaytransaction/api/oneclick/v1.2"

// MallInscription handles the inscription and deletion of customer cards for
// OneClick Mall. Create an instance with NewMallInscription.
type MallInscription struct {
	config internal.Config
}

// NewMallInscription returns a MallInscription for the given options. It
// validates the options and the commerce code and returns a
// *transbank.ValidationError on failure.
func NewMallInscription(opts transbank.Options) (*MallInscription, error) {
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
	return &MallInscription{config: cfg}, nil
}

// Start begins the inscription of the customer's card. It returns the
// inscription token and the URL of the inscription form to which the customer
// must be redirected. The response URL is where Transbank redirects the
// customer after the inscription flow.
func (m *MallInscription) Start(username, email, responseUrl string) (*OneclickMallInscriptionStartResponse, error) {
	if m.config.ValidateInputs {
		if err := internal.ValidateUsername(username); err != nil {
			return nil, err
		}
		if err := internal.ValidateEmail(email); err != nil {
			return nil, err
		}
		if err := internal.ValidateReturnURL(responseUrl); err != nil {
			return nil, err
		}
	}

	payload := map[string]string{
		"username":     username,
		"email":        email,
		"response_url": responseUrl,
	}
	var response OneclickMallInscriptionStartResponse
	if err := internal.NewRequestor(&m.config).Post(oneClickPath+"/inscriptions", payload, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Finish confirms a previously started inscription using its token, returning
// the tbk_user identifier of the enrolled card. The token is the one received
// by Transbank in the return URL (TBK_TOKEN).
func (m *MallInscription) Finish(token string) (*OneclickMallInscriptionFinishResponse, error) {
	if err := internal.ValidateToken(token); err != nil {
		return nil, err
	}
	var response OneclickMallInscriptionFinishResponse
	if err := internal.NewRequestor(&m.config).Put(fmt.Sprintf("%s/inscriptions/%s", oneClickPath, token), map[string]any{}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Delete removes a previously enrolled card for the given user. It returns an
// error if the removal fails; the API returns no response body on success.
func (m *MallInscription) Delete(tbkUser, username string) error {
	if m.config.ValidateInputs {
		if err := internal.ValidateTbkUser(tbkUser); err != nil {
			return err
		}
		if err := internal.ValidateUsername(username); err != nil {
			return err
		}
	}
	payload := map[string]string{
		"username": username,
		"tbk_user": tbkUser,
	}
	return internal.NewRequestor(&m.config).Delete(oneClickPath+"/inscriptions", payload, nil)
}

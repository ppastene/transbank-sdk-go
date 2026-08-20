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
// validates the options and returns a *transbank.ValidationError on failure.
func NewMallInscription(opts transbank.Options) (*MallInscription, error) {
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
	return &MallInscription{config: cfg}, nil
}

// Start begins the inscription of the customer's card. It returns the
// inscription token and the URL of the inscription form to which the customer
// must be redirected. The response URL is where Transbank redirects the
// customer after the inscription flow.
func (m *MallInscription) Start(username, email, responseUrl string) (*OneclickMallInscriptionStartResponse, error) {
	payload := map[string]string{
		"username":     username,
		"email":        email,
		"response_url": responseUrl,
	}
	var response OneclickMallInscriptionStartResponse
	if err := internal.NewRequestor(&m.config).Post(oneClickPath+"/inscriptions", payload, &response); err != nil {
		return nil, fmt.Errorf("oneclick inscription start: %w", err)
	}
	return &response, nil
}

// Finish confirms a previously started inscription using its token, returning
// the tbk_user identifier of the enrolled card. The token is the one received
// by Transbank in the return URL (TBK_TOKEN).
func (m *MallInscription) Finish(token string) (*OneclickMallInscriptionFinishResponse, error) {
	if err := internal.ValidateURLParam("token", token, 64); err != nil {
		return nil, err
	}
	var response OneclickMallInscriptionFinishResponse
	if err := internal.NewRequestor(&m.config).Put(fmt.Sprintf("%s/inscriptions/%s", oneClickPath, token), map[string]any{}, &response); err != nil {
		return nil, fmt.Errorf("oneclick inscription finish: %w", err)
	}
	return &response, nil
}

// Delete removes a previously enrolled card for the given user. It returns an
// error if the removal fails; the API returns no response body on success.
func (m *MallInscription) Delete(tbkUser, username string) error {
	payload := map[string]string{
		"username": username,
		"tbk_user": tbkUser,
	}
	if err := internal.NewRequestor(&m.config).Delete(oneClickPath+"/inscriptions", payload, nil); err != nil {
		return fmt.Errorf("oneclick inscription delete: %w", err)
	}
	return nil
}

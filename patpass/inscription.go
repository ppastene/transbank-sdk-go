// Package patpass provides an implementation of the Transbank PatPass
// (Patpass Comercio) REST API.
package patpass

import (
	"fmt"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/internal"
)

const inscriptionPath = "restpatpass/v1/services"

// Inscription handles the enrollment of customers in PatPass. Create an
// instance with NewInscription.
type Inscription struct {
	config internal.Config
}

// NewInscription returns an Inscription for the given options. It validates
// the options and returns a *transbank.ValidationError on failure.
func NewInscription(opts Options) (*Inscription, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	baseURL := integration_url
	if opts.Environment == Production {
		baseURL = production_url
	}
	cfg := internal.NewConfig(opts.CommerceCode, opts.Authorization, baseURL)
	cfg.Headers = map[string]string{
		"Commercecode":  opts.CommerceCode,
		"Authorization": opts.Authorization,
	}
	if opts.HTTPClient != nil {
		cfg.HTTP = opts.HTTPClient
	}
	return &Inscription{config: cfg}, nil
}

// Start enrolls a customer in PatPass with the given personal data, the return
// URL where the customer is redirected after enrolling and the final URL of the
// subscription voucher. maxAmount may be empty. It returns the enrollment token
// and the URL of the enrollment form.
func (i *Inscription) Start(url, firstName, fLastname, sLastname, rut, serviceId, finalUrl, maxAmount, phoneNumber, mobileNumber, patPassName, userEmail, commerceEmail, userAddress, userCity string) (*InscriptionStartResponse, error) {
	if err := internal.ValidateReturnURL(url); err != nil {
		return nil, err
	}
	if err := internal.ValidateReturnURL(finalUrl); err != nil {
		return nil, err
	}
	if err := internal.ValidateEmail(userEmail); err != nil {
		return nil, err
	}
	if err := internal.ValidateEmail(commerceEmail); err != nil {
		return nil, err
	}
	for _, v := range []struct {
		value string
		name  string
	}{
		{firstName, "name"},
		{fLastname, "fLastname"},
		{sLastname, "sLastname"},
		{rut, "rut"},
		{serviceId, "serviceId"},
		{phoneNumber, "phoneNumber"},
		{mobileNumber, "mobileNumber"},
		{patPassName, "patPassName"},
		{userAddress, "userAddress"},
		{userCity, "userCity"},
	} {
		if v.value == "" {
			return nil, &transbank.ValidationError{Message: v.name + " must not be empty"}
		}
	}

	payload := map[string]string{
		"url":             url,
		"nombre":          firstName,
		"pApellido":       fLastname,
		"sApellido":       sLastname,
		"rut":             rut,
		"serviceId":       serviceId,
		"finalUrl":        finalUrl,
		"commerceCode":    i.config.Credentials.CommerceCode,
		"montoMaximo":     maxAmount,
		"telefonoFijo":    phoneNumber,
		"telefonoCelular": mobileNumber,
		"nombrePatPass":   patPassName,
		"correoPersona":   userEmail,
		"correoComercio":  commerceEmail,
		"direccion":       userAddress,
		"ciudad":          userCity,
	}

	var response InscriptionStartResponse
	if err := internal.NewRequestor(&i.config).Post(fmt.Sprintf("%s/patInscription", inscriptionPath), payload, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Status returns the state of a previously started enrollment identified by its
// token, indicating whether it was authorized and the URL of the voucher.
func (i *Inscription) Status(token string) (*InscriptionStatusResponse, error) {
	if err := internal.ValidateToken(token); err != nil {
		return nil, err
	}
	payload := map[string]string{
		"token": token,
	}
	var response InscriptionStatusResponse
	if err := internal.NewRequestor(&i.config).Post(fmt.Sprintf("%s/status", inscriptionPath), payload, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

package oneclick

import (
	"fmt"

	"github.com/ppastene/transbank-sdk-go/internal/httpclient"
	"github.com/ppastene/transbank-sdk-go/internal/shared"
)

type MallInscription struct {
	requestor *shared.Requestor
}

func NewMallInscription(client shared.HTTPClientInterface, options *shared.Options) *MallInscription {
	if client == nil {
		client = httpclient.NewDefaultClient()
	}
	return &MallInscription{
		&shared.Requestor{
			Client:  client,
			Options: options,
		},
	}
}

func (m *MallInscription) Start(username, email, responseUrl string) (*OneclickMallInscriptionStartResponse, error) {
	payload := map[string]string{
		"username":     username,
		"email":        email,
		"response_url": responseUrl,
	}
	var response OneclickMallInscriptionStartResponse
	_, err := m.requestor.Do("POST", "/rswebpaytransaction/api/oneclick/v1.2/inscriptions", payload, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (m *MallInscription) Finish(token string) (*OneclickMallInscriptionFinishResponse, error) {
	var response OneclickMallInscriptionFinishResponse
	_, err := m.requestor.Do("PUT", fmt.Sprintf("/rswebpaytransaction/api/oneclick/v1.2/inscriptions/%s", token), nil, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (m *MallInscription) Delete(tbkUser, username string) error {
	payload := map[string]string{
		"username": username,
		"tbk_user": tbkUser,
	}
	_, err := m.requestor.Do("DELETE", "/rswebpaytransaction/api/oneclick/v1.2/inscriptions", payload, nil)
	if err != nil {
		return err
	}
	return nil
}

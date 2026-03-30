package webpay

import (
	"github.com/ppastene/transbank-sdk-go/internal/shared"
	"github.com/ppastene/transbank-sdk-go/oneclick"
	"github.com/ppastene/transbank-sdk-go/webpayplus"
)

const (
	ProductionURL  = shared.Live
	IntegrationURL = shared.Test
)

type WebpayPlusTransaction = webpayplus.Transaction
type WebpayPlusMallTransaction = webpayplus.MallTransaction
type WebpayPlusMallDetails = webpayplus.MallDetails
type OneclickMallInscription = oneclick.MallInscription
type OneclickMallTransaction = oneclick.MallTransaction
type OneclickMallDetails = oneclick.MallDetails
type Options = shared.Options

func NewTransaction(opt *shared.Options) *WebpayPlusTransaction {
	return webpayplus.NewTransaction(nil, opt)
}

func NewTransactionWithClient(client shared.HTTPClientInterface, opt *shared.Options) *WebpayPlusTransaction {
	return webpayplus.NewTransaction(client, opt)
}

func NewMallTransaction(opt *shared.Options) *WebpayPlusMallTransaction {
	return webpayplus.NewMallTransaction(nil, opt)
}

func NewMallTransactionWithClient(client shared.HTTPClientInterface, opt *shared.Options) *WebpayPlusMallTransaction {
	return webpayplus.NewMallTransaction(client, opt)
}

func NewOneclickMallInscription(opt *shared.Options) *OneclickMallInscription {
	return oneclick.NewMallInscription(nil, opt)
}

func NewOneclickMallInscriptionWithClient(client shared.HTTPClientInterface, opt *shared.Options) *OneclickMallInscription {
	return oneclick.NewMallInscription(client, opt)
}

func NewOneclickMallTransaction(opt *shared.Options) *OneclickMallTransaction {
	return oneclick.NewMallTransaction(nil, opt)
}

func NewOneclickMallTransactionWithClient(client shared.HTTPClientInterface, opt *shared.Options) *OneclickMallTransaction {
	return oneclick.NewMallTransaction(client, opt)
}

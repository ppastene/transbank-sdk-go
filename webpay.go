package webpay

import (
	"github.com/ppastene/transbank-sdk-go/internal/shared"
	"github.com/ppastene/transbank-sdk-go/webpayplus"
)

const (
	ProductionURL  = shared.Live
	IntegrationURL = shared.Test
)

type WebpayPlusTransaction = webpayplus.Transaction
type WebpayPlusMallTransaction = webpayplus.MallTransaction
type WebpayPlusMallDetails = webpayplus.MallDetails
type Options = shared.Options

func NewTransaction(opt *shared.Options) *Transaction {
	return webpayplus.NewTransaction(nil, opt)
}

func NewTransactionWithClient(client shared.HTTPClientInterface, opt *shared.Options) *Transaction {
	return webpayplus.NewTransaction(client, opt)
}

func NewMallTransaction(opt *shared.Options) *MallTransaction {
	return webpayplus.NewMallTransaction(nil, opt)
}

func NewMallTransactionWithClient(client shared.HTTPClientInterface, opt *shared.Options) *MallTransaction {
	return webpayplus.NewMallTransaction(client, opt)
}

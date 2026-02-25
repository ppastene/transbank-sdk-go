package webpayplus

import errors "github.com/ppastene/transbank-sdk-go/internal/shared"

func NewMallTransactionCreateException(tbkMessage string, httpCode int) *errors.WebpayError {
	return &errors.WebpayError{
		ServiceMessage: tbkMessage,
		Code:           httpCode,
	}
}

func NewMallTransactionStatusException(tbkMessage string, httpCode int) *errors.WebpayError {
	return &errors.WebpayError{
		ServiceMessage: tbkMessage,
		Code:           httpCode,
	}
}

func NewMallTransactionCommitException(tbkMessage string, httpCode int) *errors.WebpayError {
	return &errors.WebpayError{
		ServiceMessage: tbkMessage,
		Code:           httpCode,
	}
}

func NewMallTransactionRefundException(tbkMessage string, httpCode int) *errors.WebpayError {
	return &errors.WebpayError{
		ServiceMessage: tbkMessage,
		Code:           httpCode,
	}
}

func NewMallTransactionCaptureException(tbkMessage string, httpCode int) *errors.WebpayError {
	return &errors.WebpayError{
		ServiceMessage: tbkMessage,
		Code:           httpCode,
	}
}

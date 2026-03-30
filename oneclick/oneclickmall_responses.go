package oneclick

type OneclickMallInscriptionStartResponse struct {
	Token     string `json:"token"`
	UrlWebpay string `json:"url_webpay"`
}

type OneclickMallInscriptionFinishResponse struct {
	ResponseCode      int    `json:"response_code"`
	TbkUser           string `json:"tbk_user"`
	AuthorizationCode string `json:"authorization_code"`
	CardType          string `json:"card_type"`
	CardNumber        string `json:"card_number"`
}

type OneclickMallTransactionStatusResponse struct {
	BuyOrder        string                             `json:"buy_order"`
	CardDetail      OneclickMallTransactionCardDetails `json:"card_detail"`
	AccountingDate  string                             `json:"accounting_date"`
	TransactionDate string                             `json:"transaction_date"`
	Details         OneclickMallTransactionDetails     `json:"details"`
}

type OneclickMallTransactionCardDetails struct {
	CardNumber string `json:"card_number"`
}

type OneclickMallTransactionDetails struct {
	Amount             float64 `json:"amount"`
	Status             string  `json:"status"`
	AuthorizationCode  string  `json:"authorization_code"`
	PaymentTypeCode    string  `json:"payment_type_code"`
	ResponseCode       int     `json:"response_code"`
	InstallmentsNumber int     `json:"installments_number"`
	CommerceCode       string  `json:"commerce_code"`
	BuyOrder           string  `json:"buy_order"`
}

type OneclickMallTransactionAuthorizeResponse struct {
	OneclickMallTransactionStatusResponse
}

type OneclickMallTransactionRefundResponse struct {
	Type              string  `json:"type"`
	AuthorizationCode string  `json:"authorization_code"`
	AuthorizationDate string  `json:"authorization_date"`
	NullifiedAmount   float64 `json:"nullified_amount"`
	Balance           float64 `json:"balance"`
	ResponseCode      int     `json:"response_code"`
}

type OneclickMallTransactionCaptureResponse struct {
	AuthorizationCode string  `json:"authorization_code"`
	AuthorizationDate string  `json:"authorization_date"`
	CapturedAmount    float64 `json:"captured_amount"`
	ResponseCode      int     `json:"response_code"`
}

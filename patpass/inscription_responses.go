package patpass

// InscriptionStartResponse is the result of a Start call.
type InscriptionStartResponse struct {
	Token string `json:"token"`
	Url   string `json:"url"`
}

// InscriptionStatusResponse is the state of a previously started enrollment.
type InscriptionStatusResponse struct {
	Authorized bool   `json:"authorized"`
	VoucherUrl string `json:"voucherUrl"`
}

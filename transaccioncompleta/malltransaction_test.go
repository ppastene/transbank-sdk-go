package transaccioncompleta_test

import (
	"net/http"
	"testing"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/transaccioncompleta"
)

func testMallDetails() []transaccioncompleta.MallDetails {
	return []transaccioncompleta.MallDetails{
		{Amount: 10000, CommerceCode: testChildCode1, BuyOrder: testChildBuyOrder1},
		{Amount: 20000, CommerceCode: testChildCode2, BuyOrder: testChildBuyOrder2},
	}
}

func testMallInstallmentsDetails() []transaccioncompleta.MallInstallmentsDetails {
	return []transaccioncompleta.MallInstallmentsDetails{
		{CommerceCode: testChildCode1, BuyOrder: testChildBuyOrder1, InstallmentsNumber: 10},
		{CommerceCode: testChildCode2, BuyOrder: testChildBuyOrder2, InstallmentsNumber: 10},
	}
}

func testMallCommitDetails() []transaccioncompleta.MallCommitDetails {
	return []transaccioncompleta.MallCommitDetails{
		{CommerceCode: testChildCode1, BuyOrder: testChildBuyOrder1},
		{CommerceCode: testChildCode2, BuyOrder: testChildBuyOrder2},
	}
}

func testMallOptions(m *mockServer) transbank.Options {
	return transbank.Options{
		CommerceCode:   testMallCode,
		ApiKey:         testAPIKey,
		Environment:    transbank.Integration,
		HTTPClient:     m.Client(),
		ValidateInputs: true,
	}
}

func TestNewMallTransaction(t *testing.T) {
	tests := []struct {
		name string
		opts transbank.Options
		want *transbank.ValidationError
	}{
		{
			name: "valid integration",
			opts: transbank.Options{
				CommerceCode: testMallCode,
				ApiKey:       testAPIKey,
				Environment:  transbank.Integration,
			},
		},
		{
			name: "valid production",
			opts: transbank.Options{
				CommerceCode: testMallCode,
				ApiKey:       testAPIKey,
				Environment:  transbank.Production,
			},
		},
		{
			name: "invalid environment",
			opts: transbank.Options{
				CommerceCode: testMallCode,
				ApiKey:       testAPIKey,
				Environment:  transbank.Environment(99),
			},
			want: &transbank.ValidationError{},
		},
		{
			name: "non numeric commerce code",
			opts: transbank.Options{
				CommerceCode: "59705555555a",
				ApiKey:       testAPIKey,
				Environment:  transbank.Integration,
			},
			want: &transbank.ValidationError{},
		},
		{
			name: "empty api key",
			opts: transbank.Options{
				CommerceCode: testMallCode,
				Environment:  transbank.Integration,
			},
			want: &transbank.ValidationError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := transaccioncompleta.NewMallTransaction(tt.opts)
			if tt.want != nil {
				wantValidationError(t, err)
				if tx != nil {
					t.Error("expected nil MallTransaction on error")
				}
				return
			}
			if err != nil {
				t.Fatalf("NewMallTransaction returned error: %v", err)
			}
			if tx == nil {
				t.Error("expected non-nil MallTransaction")
			}
		})
	}
}

func TestMallTransactionCreate(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"token":"`+testToken+`"}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	details := testMallDetails()
	resp, err := tx.Create(testBuyOrder, testSessionID, testCardNumber, testCardExpiry, details, testCVV)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if resp.Token != testToken {
		t.Errorf("Token = %q, want %q", resp.Token, testToken)
	}

	req := server.LastRequest()
	assertMallRequest(t, req, http.MethodPost, testTransactionsPath)
	assertBody(t, req, `{"buy_order":"orden-compra-123","card_expiration_date":"22/10","card_number":"4051885600446623","cvv":"123","details":[{"amount":10000,"commerce_code":"597055555552","buy_order":"orden-hija-1"},{"amount":20000,"commerce_code":"597055555553","buy_order":"orden-hija-2"}],"session_id":"sesion-456"}`)
}

func TestMallTransactionCreateWithoutCVV(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"token":"`+testToken+`"}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	details := testMallDetails()
	resp, err := tx.Create(testBuyOrder, testSessionID, testCardNumber, testCardExpiry, details, "")
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if resp.Token != testToken {
		t.Errorf("Token = %q, want %q", resp.Token, testToken)
	}

	req := server.LastRequest()
	assertMallRequest(t, req, http.MethodPost, testTransactionsPath)
	assertBody(t, req, `{"buy_order":"orden-compra-123","card_expiration_date":"22/10","card_number":"4051885600446623","details":[{"amount":10000,"commerce_code":"597055555552","buy_order":"orden-hija-1"},{"amount":20000,"commerce_code":"597055555553","buy_order":"orden-hija-2"}],"session_id":"sesion-456"}`)
	assertBodyNotContains(t, req, "cvv")
}

func TestMallTransactionCreateValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"token":"`+testToken+`"}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	emptyDetails := []transaccioncompleta.MallDetails{}
	badDetail := []transaccioncompleta.MallDetails{
		{Amount: 10000, CommerceCode: "59705555555a", BuyOrder: testChildBuyOrder1},
	}

	tests := []struct {
		name           string
		buyOrder       string
		sessionID      string
		cardNumber     string
		cardExpiration string
		details        []transaccioncompleta.MallDetails
		cvv            string
	}{
		{name: "empty buy order", sessionID: testSessionID, cardNumber: testCardNumber, cardExpiration: testCardExpiry, details: testMallDetails(), cvv: testCVV},
		{name: "empty session id", buyOrder: testBuyOrder, cardNumber: testCardNumber, cardExpiration: testCardExpiry, details: testMallDetails(), cvv: testCVV},
		{name: "empty card number", buyOrder: testBuyOrder, sessionID: testSessionID, cardExpiration: testCardExpiry, details: testMallDetails(), cvv: testCVV},
		{name: "bad card expiration", buyOrder: testBuyOrder, sessionID: testSessionID, cardNumber: testCardNumber, cardExpiration: "2210", details: testMallDetails(), cvv: testCVV},
		{name: "empty details", buyOrder: testBuyOrder, sessionID: testSessionID, cardNumber: testCardNumber, cardExpiration: testCardExpiry, details: emptyDetails, cvv: testCVV},
		{name: "bad child commerce code", buyOrder: testBuyOrder, sessionID: testSessionID, cardNumber: testCardNumber, cardExpiration: testCardExpiry, details: badDetail, cvv: testCVV},
		{name: "cvv too long", buyOrder: testBuyOrder, sessionID: testSessionID, cardNumber: testCardNumber, cardExpiration: testCardExpiry, details: testMallDetails(), cvv: "12345"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tx.Create(tt.buyOrder, tt.sessionID, tt.cardNumber, tt.cardExpiration, tt.details, tt.cvv)
			wantValidationError(t, err)
			if server.RequestCount() != 0 {
				t.Errorf("request count = %d, want 0", server.RequestCount())
			}
		})
	}
}

func TestMallTransactionInstallments(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `[{"installments_amount":3334,"id_query_installments":11,"deferred_periods":[{"amount":1000,"period":1}]},{"installments_amount":6667,"id_query_installments":12,"deferred_periods":[]}]`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	details := testMallInstallmentsDetails()
	resp, err := tx.Installments(testToken, details)
	if err != nil {
		t.Fatalf("Installments returned error: %v", err)
	}
	if len(*resp) != 2 {
		t.Fatalf("len(resp) = %d, want 2", len(*resp))
	}
	if (*resp)[0].InstallmentsAmount != 3334 {
		t.Errorf("resp[0].InstallmentsAmount = %d, want 3334", (*resp)[0].InstallmentsAmount)
	}
	if (*resp)[0].IdQueryInstallments != 11 {
		t.Errorf("resp[0].IdQueryInstallments = %d, want 11", (*resp)[0].IdQueryInstallments)
	}
	if len((*resp)[1].DeferredPeriods) != 0 {
		t.Errorf("len(resp[1].DeferredPeriods) = %d, want 0", len((*resp)[1].DeferredPeriods))
	}

	req := server.LastRequest()
	assertMallRequest(t, req, http.MethodPost, testTransactionsPath+"/"+testToken+"/installments")
	assertBody(t, req, `[{"commerce_code":"597055555552","buy_order":"orden-hija-1","installments_number":10},{"commerce_code":"597055555553","buy_order":"orden-hija-2","installments_number":10}]`)
}

func TestMallTransactionInstallmentsValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `[]`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	badDetail := []transaccioncompleta.MallInstallmentsDetails{
		{CommerceCode: testChildCode1, BuyOrder: testChildBuyOrder1, InstallmentsNumber: 0},
	}

	tests := []struct {
		name    string
		token   string
		details []transaccioncompleta.MallInstallmentsDetails
	}{
		{name: "empty token", details: testMallInstallmentsDetails()},
		{name: "empty details", token: testToken, details: []transaccioncompleta.MallInstallmentsDetails{}},
		{name: "zero installments", token: testToken, details: badDetail},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tx.Installments(tt.token, tt.details)
			wantValidationError(t, err)
			if server.RequestCount() != 0 {
				t.Errorf("request count = %d, want 0", server.RequestCount())
			}
		})
	}
}

func TestMallTransactionCommitNoInstallments(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"buy_order":"orden-compra-123","card_detail":{"card_number":"4239"},"accounting_date":"0321","transaction_date":"2019-03-21T15:43:48.523Z","details":[{"amount":10000,"status":"AUTHORIZED","authorization_code":"123456","payment_type_code":"VN","response_code":0,"installments_number":0,"commerce_code":"597055555552","buy_order":"orden-hija-1"}]}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	resp, err := tx.Commit(testToken, testMallCommitDetails())
	if err != nil {
		t.Fatalf("Commit returned error: %v", err)
	}
	if resp.BuyOrder != testBuyOrder {
		t.Errorf("BuyOrder = %q, want %q", resp.BuyOrder, testBuyOrder)
	}
	if !resp.IsApproved() {
		t.Error("IsApproved = false, want true")
	}

	req := server.LastRequest()
	assertMallRequest(t, req, http.MethodPut, testTransactionsPath+"/"+testToken)
	assertBody(t, req, `{"details":[{"commerce_code":"597055555552","buy_order":"orden-hija-1"},{"commerce_code":"597055555553","buy_order":"orden-hija-2"}]}`)
	assertBodyNotContains(t, req, "null")
}

func TestMallTransactionCommitWithInstallments(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"buy_order":"orden-compra-123","card_detail":{"card_number":"4239"},"accounting_date":"0321","transaction_date":"2019-03-21T15:43:48.523Z","details":[{"amount":10000,"status":"AUTHORIZED","authorization_code":"123456","payment_type_code":"VN","response_code":0,"installments_number":3,"commerce_code":"597055555552","buy_order":"orden-hija-1"}]}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	details := []transaccioncompleta.MallCommitDetails{
		{CommerceCode: testChildCode1, BuyOrder: testChildBuyOrder1, IdQueryInstallments: intPtr(15), DeferredPeriodIndex: intPtr(1), GracePeriod: boolPtr(false)},
	}
	resp, err := tx.Commit(testToken, details)
	if err != nil {
		t.Fatalf("Commit returned error: %v", err)
	}
	if len(resp.Details) != 1 {
		t.Fatalf("len(Details) = %d, want 1", len(resp.Details))
	}
	if resp.Details[0].InstallmentsNumber != 3 {
		t.Errorf("Details[0].InstallmentsNumber = %d, want 3", resp.Details[0].InstallmentsNumber)
	}

	req := server.LastRequest()
	assertMallRequest(t, req, http.MethodPut, testTransactionsPath+"/"+testToken)
	assertBody(t, req, `{"details":[{"commerce_code":"597055555552","buy_order":"orden-hija-1","id_query_installments":15,"deferred_period_index":1,"grace_period":false}]}`)
	assertBodyNotContains(t, req, "null")
}

func TestMallTransactionCommitValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	badDetail := []transaccioncompleta.MallCommitDetails{
		{CommerceCode: testChildCode1, BuyOrder: ""},
	}

	tests := []struct {
		name    string
		token   string
		details []transaccioncompleta.MallCommitDetails
	}{
		{name: "empty token", details: testMallCommitDetails()},
		{name: "empty details", token: testToken, details: []transaccioncompleta.MallCommitDetails{}},
		{name: "empty child buy order", token: testToken, details: badDetail},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tx.Commit(tt.token, tt.details)
			wantValidationError(t, err)
			if server.RequestCount() != 0 {
				t.Errorf("request count = %d, want 0", server.RequestCount())
			}
		})
	}
}

func TestMallTransactionStatus(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"buy_order":"orden-compra-123","card_detail":{"card_number":"4239"},"accounting_date":"0321","transaction_date":"2019-03-21T15:43:48.523Z","details":[{"amount":10000,"status":"AUTHORIZED","authorization_code":"123456","payment_type_code":"VN","response_code":0,"installments_number":0,"commerce_code":"597055555552","buy_order":"orden-hija-1"}]}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	resp, err := tx.Status(testToken)
	if err != nil {
		t.Fatalf("Status returned error: %v", err)
	}
	if resp.CardDetail.CardNumber != "4239" {
		t.Errorf("CardDetail.CardNumber = %q, want 4239", resp.CardDetail.CardNumber)
	}
	if len(resp.Details) != 1 {
		t.Fatalf("len(Details) = %d, want 1", len(resp.Details))
	}
	if resp.Details[0].CommerceCode != testChildCode1 {
		t.Errorf("Details[0].CommerceCode = %q, want %q", resp.Details[0].CommerceCode, testChildCode1)
	}

	req := server.LastRequest()
	assertMallRequest(t, req, http.MethodGet, testTransactionsPath+"/"+testToken)
}

func TestMallTransactionStatusValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	_, err = tx.Status("")
	wantValidationError(t, err)
	if server.RequestCount() != 0 {
		t.Errorf("request count = %d, want 0", server.RequestCount())
	}
}

func TestMallTransactionRefund(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"type":"NULLIFY","authorization_code":"123456","authorization_date":"2019-03-20T20:18:20Z","nullified_amount":1000.00,"balance":0.00,"response_code":0}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	resp, err := tx.Refund(testToken, testChildBuyOrder1, testChildCode1, 1000)
	if err != nil {
		t.Fatalf("Refund returned error: %v", err)
	}
	if resp.Type != "NULLIFY" {
		t.Errorf("Type = %q, want NULLIFY", resp.Type)
	}
	if resp.NullifiedAmount != 1000 {
		t.Errorf("NullifiedAmount = %f, want 1000", resp.NullifiedAmount)
	}

	req := server.LastRequest()
	assertMallRequest(t, req, http.MethodPost, testTransactionsPath+"/"+testToken+"/refunds")
	assertBody(t, req, `{"amount":1000,"buy_order":"orden-hija-1","commerce_code":"597055555552"}`)
}

func TestMallTransactionRefundValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	tests := []struct {
		name         string
		token        string
		buyOrder     string
		commerceCode string
		amount       float64
	}{
		{name: "empty token", buyOrder: testChildBuyOrder1, commerceCode: testChildCode1, amount: 1000},
		{name: "empty buy order", token: testToken, commerceCode: testChildCode1, amount: 1000},
		{name: "bad commerce code", token: testToken, buyOrder: testChildBuyOrder1, commerceCode: "59705555555a", amount: 1000},
		{name: "zero amount", token: testToken, buyOrder: testChildBuyOrder1, commerceCode: testChildCode1, amount: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tx.Refund(tt.token, tt.buyOrder, tt.commerceCode, tt.amount)
			wantValidationError(t, err)
			if server.RequestCount() != 0 {
				t.Errorf("request count = %d, want 0", server.RequestCount())
			}
		})
	}
}

func TestMallTransactionCapture(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{"authorization_code":"123456","authorization_date":"2019-03-20T20:18:20Z","captured_amount":1000,"response_code":0}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	resp, err := tx.Capture(testToken, testChildCode1, testBuyOrder, testAuthCode, 1000)
	if err != nil {
		t.Fatalf("Capture returned error: %v", err)
	}
	if resp.CapturedAmount != 1000 {
		t.Errorf("CapturedAmount = %f, want 1000", resp.CapturedAmount)
	}
	if resp.AuthorizationCode != testAuthCode {
		t.Errorf("AuthorizationCode = %q, want %q", resp.AuthorizationCode, testAuthCode)
	}

	req := server.LastRequest()
	assertMallRequest(t, req, http.MethodPut, testTransactionsPath+"/"+testToken+"/capture")
	assertBody(t, req, `{"authorization_code":"123456","buy_order":"orden-compra-123","capture_amount":1000,"commerce_code":"597055555552"}`)
}

func TestMallTransactionCaptureValidation(t *testing.T) {
	server := newMockServer(t, http.StatusOK, `{}`)
	tx, err := transaccioncompleta.NewMallTransaction(testMallOptions(server))
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	tests := []struct {
		name          string
		token         string
		commerceCode  string
		buyOrder      string
		authorization string
		captureAmount float64
	}{
		{name: "empty token", commerceCode: testChildCode1, buyOrder: testBuyOrder, authorization: testAuthCode, captureAmount: 1000},
		{name: "bad commerce code", token: testToken, commerceCode: "59705555555a", buyOrder: testBuyOrder, authorization: testAuthCode, captureAmount: 1000},
		{name: "empty buy order", token: testToken, commerceCode: testChildCode1, authorization: testAuthCode, captureAmount: 1000},
		{name: "empty authorization code", token: testToken, commerceCode: testChildCode1, buyOrder: testBuyOrder, captureAmount: 1000},
		{name: "zero capture amount", token: testToken, commerceCode: testChildCode1, buyOrder: testBuyOrder, authorization: testAuthCode, captureAmount: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tx.Capture(tt.token, tt.commerceCode, tt.buyOrder, tt.authorization, tt.captureAmount)
			wantValidationError(t, err)
			if server.RequestCount() != 0 {
				t.Errorf("request count = %d, want 0", server.RequestCount())
			}
		})
	}
}

package oneclick_test

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/oneclick"
)

func requireIntegration(t *testing.T) {
	t.Helper()
	if os.Getenv("TBK_INTEGRATION") == "" {
		t.Skip("integration tests disabled: set TBK_INTEGRATION=1 to hit the real Transbank API")
	}
}

func integrationOptions() transbank.Options {
	return transbank.Options{
		CommerceCode: testCommerceCode,
		ApiKey:       testAPIKey,
		Environment:  transbank.Integration,
	}
}

func integrationBuyOrder() string {
	return fmt.Sprintf("integ-%d", time.Now().UnixNano()%100000000000000000)
}

func integrationUsername() string {
	return fmt.Sprintf("integ_%d", time.Now().UnixNano()%1000000000000)
}

// assertServed passes when the API served the call: either a response was
// decoded (its properties may come empty) or a *transbank.HTTPError was
// returned.
// It fails only if the SDK panics or returns an unexpected error type.
func assertServed(t *testing.T, name string, err error) {
	t.Helper()
	if err == nil {
		return
	}
	var httpErr *transbank.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("%s: unexpected error type %T: %v", name, err, err)
	}
	t.Logf("%s: served by the API as %s", name, httpErr)
}

func TestIntegrationOneClickInscriptionFlow(t *testing.T) {
	requireIntegration(t)

	ins, err := oneclick.NewMallInscription(integrationOptions())
	if err != nil {
		t.Fatalf("NewMallInscription: %v", err)
	}

	start, err := ins.Start(integrationUsername(), testEmail, testResponseURL)
	if err != nil {
		t.Fatalf("Start must always succeed: %v", err)
	}
	if start.Token == "" {
		t.Error("Start: token must not be empty")
	}
	if start.UrlWebpay == "" {
		t.Error("Start: url_webpay must not be empty")
	}

	// Finish consumes the inscription token. That token depends on the customer
	// completing the inscription form in the portal, so without that step the
	// API either returns a response with empty properties or an error message.
	finish, err := ins.Finish(start.Token)
	assertServed(t, "Finish", err)

	tbkUser := testTbkUser
	if err == nil && finish.TbkUser != "" {
		tbkUser = finish.TbkUser
	}
	err = ins.Delete(tbkUser, integrationUsername())
	assertServed(t, "Delete", err)
}

func TestIntegrationOneClickTransactionFlow(t *testing.T) {
	requireIntegration(t)

	tx, err := oneclick.NewMallTransaction(integrationOptions())
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	username := integrationUsername()
	details := []oneclick.MallDetails{
		{CommerceCode: testChildCode1, BuyOrder: integrationBuyOrder(), Amount: 1000, InstallmentsNumber: 1},
		{CommerceCode: testChildCode2, BuyOrder: integrationBuyOrder(), Amount: 2000, InstallmentsNumber: 1},
	}

	_, err = tx.Authorize(username, testTbkUser, integrationBuyOrder(), details)
	assertServed(t, "Authorize", err)

	_, err = tx.Status(integrationBuyOrder())
	assertServed(t, "Status", err)

	_, err = tx.Refund(integrationBuyOrder(), testChildCode1, integrationBuyOrder(), 500)
	assertServed(t, "Refund", err)

	_, err = tx.Capture(integrationBuyOrder(), testChildCode1, testAuthCode, 1000)
	assertServed(t, "Capture", err)
}

package webpayplus_test

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/webpayplus"
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

func integrationMallOptions() transbank.Options {
	return transbank.Options{
		CommerceCode: "597055555535", // comercio mall oficial del sandbox
		ApiKey:       testAPIKey,
		Environment:  transbank.Integration,
	}
}

func integrationBuyOrder() string {
	return fmt.Sprintf("integ-%d", time.Now().UnixNano()%100000000000000000)
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

func TestIntegrationWebpayPlusTransaction(t *testing.T) {
	requireIntegration(t)

	tx, err := webpayplus.NewTransaction(integrationOptions())
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	create, err := tx.Create(integrationBuyOrder(), testSessionID, 1000, testReturnURL)
	if err != nil {
		t.Fatalf("Create must always succeed: %v", err)
	}
	if create.Token == "" {
		t.Error("Create: token must not be empty")
	}
	if create.Url == "" {
		t.Error("Create: url must not be empty")
	}

	_, err = tx.Commit(create.Token)
	assertServed(t, "Commit", err)

	_, err = tx.Status(create.Token)
	assertServed(t, "Status", err)

	_, err = tx.Refund(create.Token, 500)
	assertServed(t, "Refund", err)

	_, err = tx.Capture(create.Token, testBuyOrder, testAuthCode, 1000)
	assertServed(t, "Capture", err)
}

func TestIntegrationWebpayPlusMallTransaction(t *testing.T) {
	requireIntegration(t)

	tx, err := webpayplus.NewMallTransaction(integrationMallOptions())
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	details := []webpayplus.MallDetails{
		{Amount: 1000, CommerceCode: testChildCode1, BuyOrder: integrationBuyOrder()},
		{Amount: 2000, CommerceCode: testChildCode2, BuyOrder: integrationBuyOrder()},
	}

	create, err := tx.Create(integrationBuyOrder(), testSessionID, testReturnURL, details)
	if err != nil {
		t.Fatalf("Create must always succeed: %v", err)
	}
	if create.Token == "" {
		t.Error("Create: token must not be empty")
	}
	if create.Url == "" {
		t.Error("Create: url must not be empty")
	}

	_, err = tx.Commit(create.Token)
	assertServed(t, "Commit", err)

	_, err = tx.Status(create.Token)
	assertServed(t, "Status", err)

	_, err = tx.Refund(create.Token, testBuyOrder, testChildCode1, 500)
	assertServed(t, "Refund", err)

	_, err = tx.Capture(create.Token, testChildCode1, testBuyOrder, testAuthCode, 1000)
	assertServed(t, "Capture", err)
}

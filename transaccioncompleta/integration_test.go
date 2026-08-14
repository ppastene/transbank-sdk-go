package transaccioncompleta_test

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/transaccioncompleta"
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
		CommerceCode: testMallCode,
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

func TestIntegrationTransaccionCompletaTransaction(t *testing.T) {
	requireIntegration(t)

	tx, err := transaccioncompleta.NewTransaction(integrationOptions())
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}

	create, err := tx.Create(integrationBuyOrder(), testSessionID, 1000, testCardNumber, testCardExpiry, testCVV)
	if err != nil {
		t.Fatalf("Create must always succeed: %v", err)
	}
	if create.Token == "" {
		t.Error("Create: token must not be empty")
	}

	_, err = tx.Installments(create.Token, 3)
	assertServed(t, "Installments", err)

	_, err = tx.Commit(create.Token, nil, nil, nil)
	assertServed(t, "Commit", err)

	_, err = tx.Status(create.Token)
	assertServed(t, "Status", err)

	_, err = tx.Refund(create.Token, 500)
	assertServed(t, "Refund", err)

	_, err = tx.Capture(create.Token, testBuyOrder, testAuthCode, 1000)
	assertServed(t, "Capture", err)
}

func TestIntegrationTransaccionCompletaMallTransaction(t *testing.T) {
	requireIntegration(t)

	tx, err := transaccioncompleta.NewMallTransaction(integrationMallOptions())
	if err != nil {
		t.Fatalf("NewMallTransaction: %v", err)
	}

	details := []transaccioncompleta.MallDetails{
		{Amount: 1000, CommerceCode: testChildCode1, BuyOrder: integrationBuyOrder()},
		{Amount: 2000, CommerceCode: testChildCode2, BuyOrder: integrationBuyOrder()},
	}

	create, err := tx.Create(integrationBuyOrder(), testSessionID, testCardNumber, testCardExpiry, details, "")
	if err != nil {
		t.Fatalf("Create must always succeed: %v", err)
	}
	if create.Token == "" {
		t.Error("Create: token must not be empty")
	}

	installments := []transaccioncompleta.MallInstallmentsDetails{
		{CommerceCode: testChildCode1, BuyOrder: integrationBuyOrder(), InstallmentsNumber: 3},
		{CommerceCode: testChildCode2, BuyOrder: integrationBuyOrder(), InstallmentsNumber: 3},
	}
	_, err = tx.Installments(create.Token, installments)
	assertServed(t, "Installments", err)

	commit := []transaccioncompleta.MallCommitDetails{
		{CommerceCode: testChildCode1, BuyOrder: integrationBuyOrder()},
		{CommerceCode: testChildCode2, BuyOrder: integrationBuyOrder()},
	}
	_, err = tx.Commit(create.Token, commit)
	assertServed(t, "Commit", err)

	_, err = tx.Status(create.Token)
	assertServed(t, "Status", err)

	_, err = tx.Refund(create.Token, integrationBuyOrder(), testChildCode1, 500)
	assertServed(t, "Refund", err)

	_, err = tx.Capture(create.Token, testChildCode1, integrationBuyOrder(), testAuthCode, 1000)
	assertServed(t, "Capture", err)
}

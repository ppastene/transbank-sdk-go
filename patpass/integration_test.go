package patpass_test

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/patpass"
)

func requireIntegration(t *testing.T) {
	t.Helper()
	if os.Getenv("TBK_INTEGRATION") == "" {
		t.Skip("integration tests disabled: set TBK_INTEGRATION=1 to hit the real Transbank API")
	}
}

func integrationOptions() patpass.Options {
	return patpass.Options{
		CommerceCode:  testPatpassCommerceCode,
		Authorization: testPatpassAuth,
		Environment:   patpass.Integration,
	}
}

func integrationServiceID() string {
	return fmt.Sprintf("integ%d", time.Now().UnixNano()%1000000000)
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

func TestIntegrationPatPassInscription(t *testing.T) {
	requireIntegration(t)

	ins, err := patpass.NewInscription(integrationOptions())
	if err != nil {
		t.Fatalf("NewInscription: %v", err)
	}

	start, err := ins.Start(
		testURL, testName, testFLastname, testSLastname, testRut,
		integrationServiceID(), testFinalURL, "", testPhone, testMobile,
		testPatpassName, testPersonEmail, testCommerceEmail, testAddress, testCity,
	)
	if err != nil {
		t.Fatalf("Start must always succeed: %v", err)
	}
	if start.Token == "" {
		t.Error("Start: token must not be empty")
	}
	if start.Url == "" {
		t.Error("Start: url must not be empty")
	}

	// Status consumes the enrollment token; without the portal flow it either
	// returns a response with empty properties or an error message.
	_, err = ins.Status(start.Token)
	assertServed(t, "Status", err)
}

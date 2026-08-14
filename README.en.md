# Transbank SDK Go

[Español](README.md) · **English**

Integration library for the Transbank API written in the Go language.

## Table of Contents

- [Implementation](#implementation)
- [Requirements](#requirements)
- [Installation](#installation)
- [Getting started](#getting-started)
- [Usage](#usage)
  - [Webpay Plus](#webpay-plus)
  - [Webpay Plus Mall](#webpay-plus-mall)
  - [OneClick](#oneclick)
  - [Transaccion Completa](#transaccion-completa)
  - [Transaccion Completa Mall](#transaccion-completa-mall)
  - [PatPass](#patpass)
- [Error handling](#error-handling)
- [HTTP client injection](#http-client-injection)

## Implementation

| Service                   | Implemented | Documentation | Tests |
|:--------------------------|:-----------:|:-------------:|:-----:|
| Webpay Plus               |      ✅     |       ✅      |   ✅  |
| Webpay Plus Mall          |      ✅     |       ✅      |   ✅  |
| OneClick Mall             |      ✅     |       ✅      |   ✅  |
| Transaccion Completa      |      ✅     |       ✅      |   ✅  |
| Transaccion Completa Mall |      ✅     |       ✅      |   ✅  |
| PatPass                   |      ✅     |       ✅      |   ✅  |

## Requirements

- Go 1.24 or higher.

## Installation

Make sure your project uses Go Modules (a `go.mod` file must exist):

```go
go mod init
```

Then fetch the SDK:

```bash
go get github.com/ppastene/transbank-sdk-go
```

## Getting started

Declare your Transbank credentials and environment in a `transbank.Options{}` and pass it as an argument to the constructor of the service you want to use.

```go
import (
	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/webpayplus"
)

opts := transbank.Options{
	CommerceCode: "597055555532",         // your commerce code
	ApiKey:       "579B532A7440BB0C9...", // your secret key
	Environment:  transbank.Integration,  // or transbank.Production
}

tx, err := webpayplus.NewTransaction(opts)
if err != nil {
	// handle the construction error (credential validation)
}
```

With that you can use the methods of the service indicated by the [Transbank documentation](https://www.transbankdevelopers.cl/documentacion/como_empezar).

## Usage

### Webpay Plus

```go
opts := transbank.Options{
	CommerceCode: "597055555532",
	ApiKey:       "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
	Environment:  transbank.Integration,
}

tx, err := webpayplus.NewTransaction(opts)
if err != nil {
	// handle the error
}

// Create a transaction
create, err := tx.Create("orden-de-compra", "id-de-sesion", 10000.50, "https://www.mi-tienda.cl/retorno")
if err != nil {
	// handle the error
}
_ = create.Token // transaction token
_ = create.Url   // payment form URL

// Confirm a transaction
commit, err := tx.Commit("token")
if err != nil {
	// handle the error
}

// Check the status of a transaction
status, err := tx.Status("token")
if err != nil {
	// handle the error
}

// Refund a transaction
refund, err := tx.Refund("token", 10000.50)
if err != nil {
	// handle the error
}

// Capture a transaction (only in environments with deferred capture)
capture, err := tx.Capture("token", "orden-de-compra", "codigo-de-autorizacion", 10000.50)
if err != nil {
	// handle the error
}
```

### Webpay Plus Mall

```go
opts := transbank.Options{
	CommerceCode: "597055555535", // Mall commerce code
	ApiKey:       "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
	Environment:  transbank.Integration,
}

mallTx, err := webpayplus.NewMallTransaction(opts)
if err != nil {
	// handle the error
}

details := []webpayplus.MallDetails{
	{
		Amount:       10000,
		CommerceCode: "597055555536", // commerce code of store 1
		BuyOrder:     "orden-detalle-1234",
	},
	{
		Amount:       10000,
		CommerceCode: "597055555537", // commerce code of store 2
		BuyOrder:     "orden-detalle-4321",
	},
}

create, err := mallTx.Create("orden-de-compra", "id-de-sesion", "https://www.mi-tienda.cl/retorno", details)
if err != nil {
	// handle the error
}
_ = create.Token

commit, err := mallTx.Commit("token")
if err != nil {
	// handle the error
}

status, err := mallTx.Status("token")
if err != nil {
	// handle the error
}

refund, err := mallTx.Refund("token", "orden-detalle-1234", "597055555536", 10000)
if err != nil {
	// handle the error
}

// Only in environments with deferred capture
capture, err := mallTx.Capture("token", "597055555536", "orden-detalle-1234", "codigo-de-autorizacion", 10000)
if err != nil {
	// handle the error
}
```

### OneClick

OneClick uses two services: `oneclick.NewMallInscription` to register/delete cards and `oneclick.NewMallTransaction` to authorize payments, check statuses, refund and capture.

```go
import (
	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/oneclick"
)

opts := transbank.Options{
	CommerceCode: "597055555541", // OneClick commerce code
	ApiKey:       "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
	Environment:  transbank.Integration,
}
```

#### Inscription

```go
ins, err := oneclick.NewMallInscription(opts)
if err != nil {
	// handle the error
}

// Start the customer's card inscription
start, err := ins.Start("nombre-de-usuario", "cliente@mi-tienda.cl", "https://www.mi-tienda.cl/retorno")
if err != nil {
	// handle the error
}
_ = start.Token     // inscription token
_ = start.UrlWebpay // inscription form URL

// Confirm the inscription with the token received on the return (TBK_TOKEN)
finish, err := ins.Finish("token")
if err != nil {
	// handle the error
}
_ = finish.TbkUser // identifier of the registered card

// Delete an inscription
err = ins.Delete("tbkUser", "nombre-de-usuario")
if err != nil {
	// handle the error
}
```

#### Transactions

```go
tx, err := oneclick.NewMallTransaction(opts)
if err != nil {
	// handle the error
}

details := []oneclick.MallDetails{
	{
		Amount:             10000,
		CommerceCode:       "597055555542", // commerce code of store 1
		BuyOrder:           "orden-detalle-1234",
		InstallmentsNumber: 3,
	},
	{
		Amount:             50000,
		CommerceCode:       "597055555543", // commerce code of store 2
		BuyOrder:           "orden-detalle-4321",
		InstallmentsNumber: 3,
	},
}

// Authorize a payment with the registered card
authorize, err := tx.Authorize("nombre-de-usuario", "tbkUser", "orden-de-compra", details)
if err != nil {
	// handle the error
}
_ = authorize.IsApproved() // true if all details were approved

// Check the status of a transaction
status, err := tx.Status("orden-de-compra")
if err != nil {
	// handle the error
}

// Refund a transaction
refund, err := tx.Refund("orden-de-compra", "597055555542", "orden-detalle-1234", 10000)
if err != nil {
	// handle the error
}

// Only in environments with deferred capture
capture, err := tx.Capture("orden-de-compra", "597055555548", "codigo-de-autorizacion", 10000)
if err != nil {
	// handle the error
}
```

### Transaccion Completa

```go
import (
	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/transaccioncompleta"
)

opts := transbank.Options{
	CommerceCode: "597055555530", // Transaccion Completa commerce code
	ApiKey:       "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
	Environment:  transbank.Integration,
}

tx, err := transaccioncompleta.NewTransaction(opts)
if err != nil {
	// handle the error
}

// Create a transaction with the card details.
// The cvv is optional (send "" if your commerce has the "no cvv" option enabled).
create, err := tx.Create("orden-de-compra", "id-de-sesion", 10000, "4051885600446623", "22/10", "123")
if err != nil {
	// handle the error
}
_ = create.Token // transaction token

// Check the available installments (optional, only if the payment is in installments)
installments, err := tx.Installments("token", 3)
if err != nil {
	// handle the error
}
_ = installments.DeferredPeriods // []DeferredPeriods, empty if there are no deferred periods

// Confirm the transaction.
// The last three parameters are optional and only sent if the payment is in installments.
idQuery := 15
periodIndex := 1
gracePeriod := false
commit, err := tx.Commit("token", &idQuery, &periodIndex, &gracePeriod)
if err != nil {
	// handle the error
}
// For a single-installment payment: tx.Commit("token", nil, nil, nil)

// Check the status of a transaction
status, err := tx.Status("token")
if err != nil {
	// handle the error
}
_ = status.IsApproved()

// Refund a transaction
refund, err := tx.Refund("token", 10000)
if err != nil {
	// handle the error
}
_ = refund.Type // "NULLIFY" or "REVERSED"

// Capture a transaction (only in environments with deferred capture)
capture, err := tx.Capture("token", "orden-de-compra", "codigo-de-autorizacion", 10000)
if err != nil {
	// handle the error
}
```

### Transaccion Completa Mall

```go
import (
	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/transaccioncompleta"
)

opts := transbank.Options{
	CommerceCode: "597055555551", // Transaccion Completa Mall commerce code
	ApiKey:       "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
	Environment:  transbank.Integration,
}

mallTx, err := transaccioncompleta.NewMallTransaction(opts)
if err != nil {
	// handle the error
}

// Create a transaction with the card details.
// The cvv is optional (send "" if your commerce has the "no cvv" option enabled).
details := []transaccioncompleta.MallDetails{
	{
		Amount:       10000,
		CommerceCode: "597055555552", // commerce code of store 1
		BuyOrder:     "orden-detalle-1234",
	},
	{
		Amount:       20000,
		CommerceCode: "597055555553", // commerce code of store 2
		BuyOrder:     "orden-detalle-4321",
	},
}

create, err := mallTx.Create("orden-de-compra", "id-de-sesion", "4051885600446623", "22/10", details, "")
if err != nil {
	// handle the error
}
_ = create.Token // transaction token

// Check the available installments per store (optional, only if the payment is in installments)
installmentsDetails := []transaccioncompleta.MallInstallmentsDetails{
	{
		CommerceCode:       "597055555552",
		BuyOrder:           "orden-detalle-1234",
		InstallmentsNumber: 3,
	},
}
installments, err := mallTx.Installments("token", installmentsDetails)
if err != nil {
	// handle the error
}
_ = installments // []MallTransactionInstallmentsDetail, one entry per store

// Confirm the transaction.
// The IdQueryInstallments, DeferredPeriodIndex and GracePeriod fields are optional
// and only sent if the payment is in installments (omit them for a single-installment payment).
idQuery := 15
periodIndex := 1
commitDetails := []transaccioncompleta.MallCommitDetails{
	{
		CommerceCode:        "597055555552",
		BuyOrder:            "orden-detalle-1234",
		IdQueryInstallments: &idQuery,
		DeferredPeriodIndex: &periodIndex,
	},
}
commit, err := mallTx.Commit("token", commitDetails)
if err != nil {
	// handle the error
}
_ = commit.IsApproved() // true if all details were approved

// Check the status of a transaction
status, err := mallTx.Status("token")
if err != nil {
	// handle the error
}
_ = status.IsApproved()

// Refund a transaction (specify the store and its order)
refund, err := mallTx.Refund("token", "orden-detalle-1234", "597055555552", 10000)
if err != nil {
	// handle the error
}
_ = refund.Type // "NULLIFY" or "REVERSED"

// Capture a transaction (only in environments with deferred capture)
capture, err := mallTx.Capture("token", "597055555552", "orden-detalle-1234", "codigo-de-autorizacion", 10000)
if err != nil {
	// handle the error
}
```

### PatPass

PatPass uses its own package (`patpass`) and credentials: the commerce code and an authorization key, not the `ApiKey` of the other services. The environment is configured with `patpass.Integration` or `patpass.Production`.

```go
import (
	"github.com/ppastene/transbank-sdk-go/patpass"
)

opts := patpass.Options{
	CommerceCode:  "28299257",      // your PatPass commerce code
	Authorization: "cxxXQgGD9vrVe4M41FIt", // your authorization key
	Environment:   patpass.Integration,
}

ins, err := patpass.NewInscription(opts)
if err != nil {
	// handle the error
}

// Start the customer's inscription
start, err := ins.Start(
	"https://www.mi-tienda.cl/finalizar-suscripcion", // url
	"Diego",           // firstName
	"Sanchez",         // fLastname
	"Valdovinos",      // sLastname
	"12345678-9",      // rut
	"323123",          // serviceId
	"https://www.mi-tienda.cl/voucher", // finalUrl
	"",                // maxAmount (can be empty)
	"57508624",        // landlinePhone
	"57508624",        // mobilePhone
	"Help - 8050014",  // patpassName
	"persona@test.cl", // personEmail
	"comercio@test.cl", // commerceEmail
	"Merced 156, Santiago, Chile", // address
	"Santiago",        // city
)
if err != nil {
	// handle the error
}
_ = start.Token // inscription token
_ = start.Url   // inscription form URL

// Check the status of the inscription with the token received on the return
status, err := ins.Status(start.Token)
if err != nil {
	// handle the error
}
_ = status.Authorized // true if the inscription was approved
_ = status.VoucherUrl // voucher URL
```

## Error handling

The SDK returns typed errors that reflect where the operation failed, all
distinguishable with `errors.As`:

- `*transbank.ValidationError`: invalid parameters or credentials; the API was not called.
- `*transbank.TransportError`: the request or response could not be completed (network, encoding, parsing). `Err` is the root cause.
- `*transbank.HTTPError`: the Transbank API responded with a non-2xx status code. `StatusCode` and `Body` (raw response) help diagnose it.

```go
import (
	"errors"
	"fmt"

	"github.com/ppastene/transbank-sdk-go"
)

_, err := tx.Status("token")
if err != nil {
	var valErr *transbank.ValidationError
	var tpErr *transbank.TransportError
	var httpErr *transbank.HTTPError
	switch {
	case errors.As(err, &valErr):
		fmt.Println("SDK validation error:", valErr.Message)
	case errors.As(err, &tpErr):
		fmt.Println("Communication error:", tpErr.Err)
	case errors.As(err, &httpErr):
		fmt.Printf("HTTP error %d: %s\n", httpErr.StatusCode, httpErr.Body)
	}
}
```

## HTTP client injection

The SDK communicates with the Transbank API through the `transbank.HTTPClient` interface (a single `Do(req *http.Request) (*http.Response, error)` method), which the standard library's `http.Client` satisfies. By default an internal client with a 30 second timeout is used. If you need different behavior (proxy, TLS, custom timeouts or a mock for tests), inject it in `Options.HTTPClient`.

With a standard library `http.Client`:

```go
import (
	"net/http"
	"time"
)

client := &http.Client{
	Timeout: 15 * time.Second,
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	},
}

opts := transbank.Options{
	CommerceCode: "597055555532",
	ApiKey:       "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
	Environment:  transbank.Integration,
	HTTPClient:   client,
}

tx, err := webpayplus.NewTransaction(opts)
if err != nil {
	// handle the error
}
```

`*http.Client` already implements `Do`, so it is injected as is.

### Adapting any HTTP client

The `transbank.HTTPClient` interface only requires one method: `Do(req *http.Request) (*http.Response, error)`. If you have your own client or use an external one like [Resty](https://github.com/go-resty/resty), there are two ways to adapt it.

**1. Use the library's underlying `*http.Client`.** Many libraries are built on top of `net/http` and expose their internal client, which already implements `Do`. Resty exposes it with `GetClient()`:

```go
import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/ppastene/transbank-sdk-go"
)

restyClient := resty.New().SetTimeout(15 * time.Second)

opts := transbank.Options{
	CommerceCode: "597055555532",
	ApiKey:       "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
	Environment:  transbank.Integration,
	HTTPClient:   restyClient.GetClient(), // Resty's underlying *http.Client
}
```

**2. Write an adapter that implements `Do`.** If the library uses its own types, or you want requests to go through its pipeline (Resty's retries and middlewares), write a wrapper that translates `*http.Request` into the library's types and returns a valid `*http.Response` (with `StatusCode` and `Body`), which is all the SDK consumes:

```go
import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/ppastene/transbank-sdk-go"
)

// RestyAdapter wraps a resty.Client so it satisfies transbank.HTTPClient.
type RestyAdapter struct {
	client *resty.Client
}

func (a *RestyAdapter) Do(req *http.Request) (*http.Response, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	r := a.client.R().SetBody(body)
	for name, values := range req.Header {
		r.SetHeader(name, values[0])
	}

	rr, err := r.Execute(req.Method, req.URL.String())
	if err != nil {
		return nil, err
	}

	return &http.Response{
		StatusCode: rr.StatusCode(),
		Header:     rr.Header(),
		Body:       io.NopCloser(bytes.NewReader(rr.Body())),
	}, nil
}

restyClient := resty.New().
	SetRetryCount(3).
	SetRetryWaitTime(200 * time.Millisecond)

opts := transbank.Options{
	CommerceCode: "597055555532",
	ApiKey:       "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
	Environment:  transbank.Integration,
	HTTPClient:   &RestyAdapter{client: restyClient},
}
```

If not specified, the SDK uses its internal client with a 30 second timeout.

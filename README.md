# Transbank SDK Go

**Español** · [English](README.en.md)

Librería de integración con la API de Transbank escrita en el lenguaje Go.

## Índice

- [Implementación](#implementacion)
- [Requisitos](#requisitos)
- [Instalación](#instalación)
- [Primeros pasos](#primeros-pasos)
- [Uso](#uso)
  - [Webpay Plus](#webpay-plus)
  - [Webpay Plus Mall](#webpay-plus-mall)
  - [OneClick](#oneclick)
  - [Transacción Completa](#transacción-completa)
  - [Transacción Completa Mall](#transacción-completa-mall)
  - [PatPass](#patpass)
- [Manejo de errores](#manejo-de-errores)
- [Inyección de cliente HTTP](#inyección-de-cliente-http)

## Implementación

| Servicio                  | Implementado | Documentación | Tests |
|:--------------------------|:------------:|:-------------:|:-----:|
| Webpay Plus               |      ✅      |       ✅      |   ✅  |
| Webpay Plus Mall          |      ✅      |       ✅      |   ✅  |
| OneClick Mall             |      ✅      |       ✅      |   ✅  |
| Transacción Completa      |      ✅      |       ✅      |   ✅  |
| Transacción Completa Mall |      ✅      |       ✅      |   ✅  |
| PatPass                   |      ✅      |       ✅      |   ✅  |

## Requisitos

- Go 1.24 o superior.

## Instalación

Asegúrese de que su proyecto use Go Modules (debe existir un archivo `go.mod`):

```go
go mod init
```

Luego obtenga el SDK:

```bash
go get github.com/ppastene/transbank-sdk-go
```

## Primeros pasos

Declare sus credenciales y ambiente de Transbank en un `transbank.Options{}` y páselo como argumento al constructor del servicio que quiera utilizar.

```go
import (
	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/webpayplus"
)

opts := transbank.Options{
	CommerceCode: "597055555532",         // su código de comercio
	ApiKey:       "579B532A7440BB0C9...", // su llave secreta
	Environment:  transbank.Integration,  // o transbank.Production
}

tx, err := webpayplus.NewTransaction(opts)
if err != nil {
	// manejar el error de construcción (validación de credenciales)
}
```

Con eso ya puede usar los métodos del servicio indicado por la [documentación de Transbank](https://www.transbankdevelopers.cl/documentacion/como_empezar).

## Uso

### Webpay Plus

```go
opts := transbank.Options{
	CommerceCode: "597055555532",
	ApiKey:       "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
	Environment:  transbank.Integration,
}

tx, err := webpayplus.NewTransaction(opts)
if err != nil {
	// manejar el error
}

// Crear una transacción
create, err := tx.Create("orden-de-compra", "id-de-sesion", 10000.50, "https://www.mi-tienda.cl/retorno")
if err != nil {
	// manejar el error
}
_ = create.Token // token de la transacción
_ = create.Url   // URL del formulario de pago

// Confirmar una transacción
commit, err := tx.Commit("token")
if err != nil {
	// manejar el error
}

// Consultar el estado de una transacción
status, err := tx.Status("token")
if err != nil {
	// manejar el error
}

// Reembolsar una transacción
refund, err := tx.Refund("token", 10000.50)
if err != nil {
	// manejar el error
}

// Capturar una transacción (solo en ambientes con captura diferida)
capture, err := tx.Capture("token", "orden-de-compra", "codigo-de-autorizacion", 10000.50)
if err != nil {
	// manejar el error
}
```

### Webpay Plus Mall

```go
opts := transbank.Options{
	CommerceCode: "597055555535", // código de comercio Mall
	ApiKey:       "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
	Environment:  transbank.Integration,
}

mallTx, err := webpayplus.NewMallTransaction(opts)
if err != nil {
	// manejar el error
}

details := []webpayplus.MallDetails{
	{
		Amount:       10000,
		CommerceCode: "597055555536", // código de comercio de la tienda 1
		BuyOrder:     "orden-detalle-1234",
	},
	{
		Amount:       10000,
		CommerceCode: "597055555537", // código de comercio de la tienda 2
		BuyOrder:     "orden-detalle-4321",
	},
}

create, err := mallTx.Create("orden-de-compra", "id-de-sesion", "https://www.mi-tienda.cl/retorno", details)
if err != nil {
	// manejar el error
}
_ = create.Token

commit, err := mallTx.Commit("token")
if err != nil {
	// manejar el error
}

status, err := mallTx.Status("token")
if err != nil {
	// manejar el error
}

refund, err := mallTx.Refund("token", "orden-detalle-1234", "597055555536", 10000)
if err != nil {
	// manejar el error
}

// Solo en ambientes con captura diferida
capture, err := mallTx.Capture("token", "597055555536", "orden-detalle-1234", "codigo-de-autorizacion", 10000)
if err != nil {
	// manejar el error
}
```

### OneClick

OneClick usa dos servicios: `oneclick.NewMallInscription` para inscribir/eliminar tarjetas y `oneclick.NewMallTransaction` para autorizar pagos, consultar estados, reembolsar y capturar.

```go
import (
	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/oneclick"
)

opts := transbank.Options{
	CommerceCode: "597055555541", // código de comercio OneClick
	ApiKey:       "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
	Environment:  transbank.Integration,
}
```

#### Inscripción

```go
ins, err := oneclick.NewMallInscription(opts)
if err != nil {
	// manejar el error
}

// Iniciar la inscripción de la tarjeta del cliente
start, err := ins.Start("nombre-de-usuario", "cliente@mi-tienda.cl", "https://www.mi-tienda.cl/retorno")
if err != nil {
	// manejar el error
}
_ = start.Token     // token de la inscripción
_ = start.UrlWebpay // URL del formulario de inscripción

// Confirmar la inscripción con el token recibido en el retorno (TBK_TOKEN)
finish, err := ins.Finish("token")
if err != nil {
	// manejar el error
}
_ = finish.TbkUser // identificador de la tarjeta inscrita

// Eliminar una inscripción
err = ins.Delete("tbkUser", "nombre-de-usuario")
if err != nil {
	// manejar el error
}
```

#### Transacciones

```go
tx, err := oneclick.NewMallTransaction(opts)
if err != nil {
	// manejar el error
}

details := []oneclick.MallDetails{
	{
		Amount:             10000,
		CommerceCode:       "597055555542", // código de comercio de la tienda 1
		BuyOrder:           "orden-detalle-1234",
		InstallmentsNumber: 3,
	},
	{
		Amount:             50000,
		CommerceCode:       "597055555543", // código de comercio de la tienda 2
		BuyOrder:           "orden-detalle-4321",
		InstallmentsNumber: 3,
	},
}

// Autorizar un pago con la tarjeta inscrita
authorize, err := tx.Authorize("nombre-de-usuario", "tbkUser", "orden-de-compra", details)
if err != nil {
	// manejar el error
}
_ = authorize.IsApproved() // true si todos los detalles fueron aprobados

// Consultar el estado de una transacción
status, err := tx.Status("orden-de-compra")
if err != nil {
	// manejar el error
}

// Reembolsar una transacción
refund, err := tx.Refund("orden-de-compra", "597055555542", "orden-detalle-1234", 10000)
if err != nil {
	// manejar el error
}

// Solo en ambientes con captura diferida
capture, err := tx.Capture("orden-de-compra", "597055555548", "codigo-de-autorizacion", 10000)
if err != nil {
	// manejar el error
}
```

### Transacción Completa

```go
import (
	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/transaccioncompleta"
)

opts := transbank.Options{
	CommerceCode: "597055555530", // código de comercio Transacción Completa
	ApiKey:       "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
	Environment:  transbank.Integration,
}

tx, err := transaccioncompleta.NewTransaction(opts)
if err != nil {
	// manejar el error
}

// Crear una transacción con los datos de la tarjeta.
// El cvv es opcional (envíe "" si su comercio tiene la opción "sin cvv" habilitada).
create, err := tx.Create("orden-de-compra", "id-de-sesion", 10000, "4051885600446623", "22/10", "123")
if err != nil {
	// manejar el error
}
_ = create.Token // token de la transacción

// Consultar las cuotas disponibles (opcional, solo si el pago es en cuotas)
installments, err := tx.Installments("token", 3)
if err != nil {
	// manejar el error
}
_ = installments.DeferredPeriods // []DeferredPeriods, vacío si no hay periodos diferidos

// Confirmar la transacción.
// Los tres últimos parámetros son opcionales y solo se envían si el pago es en cuotas.
idQuery := 15
periodIndex := 1
gracePeriod := false
commit, err := tx.Commit("token", &idQuery, &periodIndex, &gracePeriod)
if err != nil {
	// manejar el error
}
// Para pago en una sola cuota: tx.Commit("token", nil, nil, nil)

// Consultar el estado de una transacción
status, err := tx.Status("token")
if err != nil {
	// manejar el error
}
_ = status.IsApproved()

// Reembolsar una transacción
refund, err := tx.Refund("token", 10000)
if err != nil {
	// manejar el error
}
_ = refund.Type // "NULLIFY" o "REVERSED"

// Capturar una transacción (solo en ambientes con captura diferida)
capture, err := tx.Capture("token", "orden-de-compra", "codigo-de-autorizacion", 10000)
if err != nil {
	// manejar el error
}
```

### Transacción Completa Mall

```go
import (
	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/transaccioncompleta"
)

opts := transbank.Options{
	CommerceCode: "597055555551", // código de comercio Transacción Completa Mall
	ApiKey:       "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
	Environment:  transbank.Integration,
}

mallTx, err := transaccioncompleta.NewMallTransaction(opts)
if err != nil {
	// manejar el error
}

// Crear una transacción con los datos de la tarjeta.
// El cvv es opcional (envíe "" si su comercio tiene la opción "sin cvv" habilitada).
details := []transaccioncompleta.MallDetails{
	{
		Amount:       10000,
		CommerceCode: "597055555552", // código de comercio de la tienda 1
		BuyOrder:     "orden-detalle-1234",
	},
	{
		Amount:       20000,
		CommerceCode: "597055555553", // código de comercio de la tienda 2
		BuyOrder:     "orden-detalle-4321",
	},
}

create, err := mallTx.Create("orden-de-compra", "id-de-sesion", "4051885600446623", "22/10", details, "")
if err != nil {
	// manejar el error
}
_ = create.Token // token de la transacción

// Consultar las cuotas disponibles por tienda (opcional, solo si el pago es en cuotas)
installmentsDetails := []transaccioncompleta.MallInstallmentsDetails{
	{
		CommerceCode:       "597055555552",
		BuyOrder:           "orden-detalle-1234",
		InstallmentsNumber: 3,
	},
}
installments, err := mallTx.Installments("token", installmentsDetails)
if err != nil {
	// manejar el error
}
_ = installments // []MallTransactionInstallmentsDetail, una entrada por tienda

// Confirmar la transacción.
// Los campos IdQueryInstallments, DeferredPeriodIndex y GracePeriod son opcionales
// y solo se envían si el pago es en cuotas (omítalos para pago en una sola cuota).
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
	// manejar el error
}
_ = commit.IsApproved() // true si todos los detalles fueron aprobados

// Consultar el estado de una transacción
status, err := mallTx.Status("token")
if err != nil {
	// manejar el error
}
_ = status.IsApproved()

// Reembolsar una transacción (se indica la tienda y su orden)
refund, err := mallTx.Refund("token", "orden-detalle-1234", "597055555552", 10000)
if err != nil {
	// manejar el error
}
_ = refund.Type // "NULLIFY" o "REVERSED"

// Capturar una transacción (solo en ambientes con captura diferida)
capture, err := mallTx.Capture("token", "597055555552", "orden-detalle-1234", "codigo-de-autorizacion", 10000)
if err != nil {
	// manejar el error
}
```

### PatPass

PatPass usa su propio paquete (`patpass`) y credenciales: el código de comercio y una llave de autorización, no la `ApiKey` de los demás servicios. El ambiente se configura con `patpass.Integration` o `patpass.Production`.

```go
import (
	"github.com/ppastene/transbank-sdk-go/patpass"
)

opts := patpass.Options{
	CommerceCode:  "28299257",      // su código de comercio PatPass
	Authorization: "cxxXQgGD9vrVe4M41FIt", // su llave de autorización
	Environment:   patpass.Integration,
}

ins, err := patpass.NewInscription(opts)
if err != nil {
	// manejar el error
}

// Iniciar la inscripción del cliente
start, err := ins.Start(
	"https://www.mi-tienda.cl/finalizar-suscripcion", // url
	"Diego",           // nombre
	"Sanchez",         // pApellido
	"Valdovinos",      // sApellido
	"12345678-9",      // rut
	"323123",          // serviceId
	"https://www.mi-tienda.cl/voucher", // finalUrl
	"",                // montoMaximo (puede ir vacío)
	"57508624",        // telefonoFijo
	"57508624",        // telefonoCelular
	"Help - 8050014",  // nombrePatPass
	"persona@test.cl", // correoPersona
	"comercio@test.cl", // correoComercio
	"Merced 156, Santiago, Chile", // direccion
	"Santiago",        // ciudad
)
if err != nil {
	// manejar el error
}
_ = start.Token // token de la inscripción
_ = start.Url   // URL del formulario de inscripción

// Consultar el estado de la inscripción con el token recibido en el retorno
status, err := ins.Status(start.Token)
if err != nil {
	// manejar el error
}
_ = status.Authorized // true si la inscripción fue aprobada
_ = status.VoucherUrl // URL del voucher
```

## Manejo de errores

El SDK devuelve errores tipados que reflejan dónde falló la operación, todos
discriminables con `errors.As`:

- `*transbank.ValidationError`: parámetros o credenciales inválidas; la API no se llamó.
- `*transbank.TransportError`: no se pudo completar el request o procesar la respuesta (red, encoding, parsing). `Err` es la causa raíz.
- `*transbank.HTTPError`: la API de Transbank respondió con un código distinto de 2xx. `StatusCode` y `Body` (respuesta cruda) permiten diagnosticar.

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
		fmt.Println("Error de validación del SDK:", valErr.Message)
	case errors.As(err, &tpErr):
		fmt.Println("Error de comunicación:", tpErr.Err)
	case errors.As(err, &httpErr):
		fmt.Printf("Error HTTP %d: %s\n", httpErr.StatusCode, httpErr.Body)
	}
}
```

## Inyección de cliente HTTP

El SDK se comunica con la API de Transbank a través de la interfaz `transbank.HTTPClient` (un método `Do(req *http.Request) (*http.Response, error)`), que el `http.Client` de la librería estándar satisface. Por defecto se usa un cliente interno con timeout de 30 segundos. Si necesita un comportamiento distinto (proxy, TLS, timeouts personalizados o un mock para tests), inyéctelo en `Options.HTTPClient`.

Con un `http.Client` de la librería estándar:

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
	// manejar el error
}
```

`*http.Client` ya implementa `Do`, así que se inyecta tal cual.

### Adaptar cualquier cliente HTTP

La interfaz `transbank.HTTPClient` solo pide un método: `Do(req *http.Request) (*http.Response, error)`. Si usted tiene su propio cliente o usa uno externo como [Resty](https://github.com/go-resty/resty), hay dos formas de adaptarlo.

**1. Usar el `*http.Client` subyacente de la librería.** Muchas librerías están construidas sobre `net/http` y exponen su cliente interno, que ya implementa `Do`. Resty lo expone con `GetClient()`:

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
	HTTPClient:   restyClient.GetClient(), // *http.Client subyacente de Resty
}
```

**2. Escribir un adaptador que implemente `Do`.** Si la librería usa tipos propios, o quiere que los requests pasen por su pipeline (retries y middlewares de Resty), escriba un wrapper que traduzca `*http.Request` a los tipos de la librería y devuelva un `*http.Response` válido (con `StatusCode` y `Body`), que es lo único que el SDK consume:

```go
import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/ppastene/transbank-sdk-go"
)

// AdaptadorResty envuelve un resty.Client para que cumpla con transbank.HTTPClient.
type AdaptadorResty struct {
	cliente *resty.Client
}

func (a *AdaptadorResty) Do(req *http.Request) (*http.Response, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	r := a.cliente.R().SetBody(body)
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
	HTTPClient:   &AdaptadorResty{cliente: restyClient},
}
```

Si no se especifica, el SDK usa su cliente interno con timeout de 30 segundos.

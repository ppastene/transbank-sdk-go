# Transbank SDK Go
Libreria de integración con la API de Transbank escrito en el lenguaje Go

**Español** | [English](./README_en.md)
## Índice
- [Implementación](#implementacion)
- [Requisitos](#requisitos)
- [Instalación](#instalación)
- [Primeros pasos](#primeros-pasos)
- [Uso](#uso)
  - [Webpay Plus](#webpay-plus)
  - [Webpay Plus Mall](#webpay-plus-mall)
  - [Oneclick](#oneclick)
- [Manejo de errores](#manejo-de-errores)
- [Inyección de cliente HTTP](#inyección-de-cliente-http)
    - [Ejemplo de inyección con Resty](#ejemplo-de-inyección-con-resty)
- [Ejemplo de Implementación](#ejemplo-de-implementacion)
## Implementación
|       Servicio       | Implementado | Documentacion | Pruebas Unitarias |
|:--------------------:|:------------:|:-------------:|:-----------------:|
| Webpay Plus          |      ✅      |       ❌      |         ✅        |
| Webpay Plus Mall     |      ✅      |       ❌      |         ✅        |
| OneClick             |      ✅      |       ❌      |         ❌        |
| PatPass              |      ❌      |       ❌      |         ❌        |
| Transaccion Completa |      ❌      |       ❌      |         ❌        |
## Requisitos
- Go 1.21.0
## Instalación
Asegurese que su proyecto esté usando Go Modules (deberia haber un archivo go.mod en la raiz):
```go
go mod init
```
Luego, importe transbank-sdk-go usando import
```go
import github.com/ppastene/transbank-sdk-go
```
Tambien puede usar el comando go get a traves de la terminal
```go
go get -u github.com/ppastene/transbank-sdk-go
```
## Primeros pasos
Declare las variables de su ambiente de Transbank en un struct transbank.Options{} y paselo como argumento en el servicio que quiera utilizar
```go
import "github.com/ppastene/transbank-sdk-go"

options := &transbank.Options{
    ApiKey: "su api key",
    CommerceCode: "su codigo de comercio",
    Environment: transbank.IntegrationURL, //Opciones: transbank.IntegrationURL para ambiente de integracion, transbank.ProductionURL para ambiente de produccion.
}
tx := transbank.NewTransaction(options)
```
Con eso ya puede usar los metodos del servicio indicado por la [documentación de Transbank](https://www.transbankdevelopers.cl/documentacion/como_empezar).
## Uso
### Webpay Plus
```go
import "github.com/ppastene/transbank-sdk-go"

options := &transbank.Options{
    ApiKey: "597055555532",
    CommerceCode: "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
    Environment: transbank.IntegrationURL,
}
transaction := transbank.NewTransaction(options)
res, err := transaction.Create("buy_order", "session_id", "amount", "http://url-de-retorno.cl")
res, err := transaction.Commit("token")
res, err := transaction.Status("token")
res, err := transaction.Refund("token", "amount")
res, err := transaction.Capture("token", "buy_order", "authorization_code", "amount") // Solo en ambientes con opción diferido
```
### Webpay Plus Mall
```go
import "github.com/ppastene/transbank-sdk-go"

options := &transbank.Options{
    ApiKey: "597055555532",
    CommerceCode: "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
    Environment: transbank.IntegrationURL,
}
mallTransaction := transbank.NewMallTransaction(options)
details := transbank.[]WebpayPlusMallDetails{
    {
        Amount: 10000,
        CommerceCode: "commerceCodeStoreOne",
        BuyOrder: "ordenCompraDetalle1234",
    },
    {
        Amount: 10000,
        CommerceCode: "commerceCodeStoreTwo",
        BuyOrder: "ordenCompraDetalle4321",
    },
}
res, err := mallTransaction.Create("buy_order", "session_id", "http://url-de-retorno.cl", details)
res, err := mallTransaction.Commit("token")
res, err := mallTransaction.Status("token")
res, err := mallTransaction.Refund("token", "buy_order", "amount", "child_commerce_code")
res, err := mallTransaction.Capture("token", "child_commerce_code", "buy_order", "authorization_code", "amount") // Solo en ambientes con opción diferido
```
### Oneclick
```go
import "github.com/ppastene/transbank-sdk-go"

inscriptionOptions := &transbank.Options{
    ApiKey: "597055555532",
    CommerceCode: "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
    Environment: transbank.IntegrationURL,
}
oneclickInscrpition := transbank.NewOneclickMallInscription(inscriptionOptions)
res, err := oneclickInscrpition.Start("user", "email", "http://url-de-retorno.cl")
res, err := oneclickInscrpition.Finish("token")
res, err := oneclickInscription.Delete("token", "user")

transactionOptions := &transbank.Options{
    ApiKey: "597055555532",
    CommerceCode: "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
    Environment: transbank.IntegrationURL,
}

details := transbank.[]OneclickMallDetails{
    {
        CommerceCode: "commerceCodeStoreOne",
	    BuyOrder: "ordenCompraDetalle1234",
	    Amount: 10000,
	    InstallmentsNumber 0,
    },
    {
        CommerceCode: "commerceCodeStoreTwo",
	    BuyOrder: "ordenCompraDetalle4321",
	    Amount: 10000,
	    InstallmentsNumber 0,
    },
}

oneclickTransaction := transbank.NewOneclickMallTransaction(transactionOptions)
res, err := oneclickTransaction.Authorize("username", "token", "buy_order", details)
res, err := oneclickTransaction.Status("buy_order")
res, err := oneclickTransaction.Refund("buy_order", "child_commerce_code", "child_buy_order", "amount")
res, err := oneclickTransaction.Capture("buy_order", "commerce_code", "authorization_code", "amount") // Solo en ambientes con opción diferido
```
### Manejo de errores
Transbank SDK Go maneja un struct para todos los errores que puedan haber
```go
type WebpayError struct {
	Code           int      
	ServiceMessage string   
	Cause          error
}
```
- Code:
    - -1: Errores ocurridos en el SDK como validaciones
    - 0: Errores con el cliente HTTP, fallos en la DNS o hostname
    -  1 hacia arriba: Corresponde al codigo HTTP devueltos por la API de Transbank
- ServiceMessage: Mensaje de error a modo ilustrativo
- Cause: El error como tal

Si se imprime el error se usará el metodo Error() que devuelve un string de ServiceMessage + Cause
```go
options := &transbank.Options{
    ApiKey: "su api key",
    CommerceCode: "su codigo de comercio",
    Environment: transbank.IntegrationURL,
}
transaction := transbank.NewTransaction(options)

res, err := transaction.Status("")
if err != nil {
    fmt.Println(err)
}
/*
    Imprime en consola
    SDK Validation Error: 'token' cannot be empty
*/
```
Puede usar errors.Unwrap(err) para obtener directamente el error
```go
options := &transbank.Options{
    ApiKey: "su api key",
    CommerceCode: "su codigo de comercio",
    Environment: transbank.IntegrationURL,
}
transaction := transbank.NewTransaction(options)

res, err := transaction.Status("un-token-que-no-existe")
if err != nil {
    fmt.Println(errors.Unwrap(err))
}
/*
    Imprime en consola
    Invalid value for parameter: token
*/
```
## Inyección de cliente HTTP
El SDK posee un cliente HTTP para la comunicación con la API de Transbank. Puede reemplazar aquel cliente con el que desee usar en el SDK.

Estos son los servicios de los cuales puede inyectar un cliente HTTP
```go
NewTransactionWithClient(client shared.HTTPClientInterface, opt *shared.Options) // Para Webpay Plus
NewMallTransactionWithClient(client shared.HTTPClientInterface, opt *shared.Options) // Para Webpay Plus Mall
NewOneclickMallInscriptionWithClient(client shared.HTTPClientInterface, opt *shared.Options) // Para inscrpitiones con Oneclick
NewOneclickMallTransactionWithClient(client shared.HTTPClientInterface, opt *shared.Options) // Para transacciones con Oneclick
```
La interface contiene el siguiente metodo
```go
type HTTPClientInterface interface {
	Request(method string, url string, headers map[string]string, payload any) ([]byte, int, error)
}
```
Los parametros del metodo son:
- method: El metodo HTTP a ejecutar. Se escribe en mayuscula
- url: La url a consultar
- headers: Los headers de la peticion
- payload: El payload de la peticion. Tiene que ser capaz de convertirse a JSON usando json.Marshal(). Puede venir vacio

Las respuestas que devuelve el metodo son:
- []byte: La respuesta en crudo de la petición.
- int: El codigo HTTP de la respuesta
- error: Error en caso de problemas de comunicación, DNS, hostname, unmarshall
### Ejemplo de inyección con Resty
```go
type RestyClient struct {
	*resty.Client
}

func (c *RestyClient) Request(method string, url string, headers map[string]string, payload any) ([]byte, int, error) {
	req := c.Client.R().
		SetHeaders(headers).
		SetBody(payload)

	resp, err := req.Execute(method, url)
	if err != nil {
		return nil, 0, err
	}

	return resp.Body(), resp.StatusCode(), nil
}

func main() {
	client := &RestyClient{resty.New()}
	options := &transbank.Options{
		ApiKey:       "597055555532",
		CommerceCode: "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
		Environment:  transbank.IntegrationURL,
	}
	tx := transbank.NewTransactionWithClient(client, options)
	resp, err := tx.Status("token")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}
```
## Ejemplo de implementacion
En [este repositorio](https://github.com/ppastene/transbank-sdk-go-example) puede encontrar un ejemplo de implementación del SDK usando Goravel. Siga las instrucciones del README y la [documentación de Goravel](https://www.goravel.dev/getting-started/installation.html) para mas información.

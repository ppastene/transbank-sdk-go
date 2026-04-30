# Transbank SDK Go
Library to connect with the Transbank API written in the Go language

[Spanish](./README.md) | **English**
## Index
- [Implementation](#implementation)
- [Requirements](#requirements)
- [Installation](#installation)
- [First steps](#first-steps)
- [Use](#use)
  - [Webpay Plus](#webpay-plus)
  - [Webpay Plus Mall](#webpay-plus-mall)
  - [Oneclick](#oneclick)
- [Error management](#error-management)
- [HTTP Client injection](#http-client-injection)
    - [Example using Resty client](#example-using-resty-client)
- [Implementation Example](#implementation-example)
## Implementation
|        Service       | Implemented | Documentation | Unit Testing |
|:--------------------:|:-----------:|:-------------:|:------------:|
| Webpay Plus          |      ✅     |       ❌      |      ✅      |
| Webpay Plus Mall     |      ✅     |       ❌      |      ✅      |
| OneClick             |      ✅     |       ❌      |      ❌      |
| PatPass              |      ❌     |       ❌      |      ❌      |
| Transaccion Completa |      ❌     |       ❌      |      ❌      |
## Requirements
- Go 1.21.0
## Installation
Be sure your project is using Go modules (your project must have a go.mod file in your root):
```go
go mod init
```
Then, import transbank-sdk-go using import
```go
import github.com/ppastene/transbank-sdk-go
```
Also you can use the go get command in the terminal
```go
go get -u github.com/ppastene/transbank-sdk-go
```
## First steps
Declare the variables of your Transbank environment in the transbank.Options{} struct, and pass it as an argumento on the service you wanna use
```go
import webpay "github.com/ppastene/transbank-sdk-go"

options := &transbank.Options{
    ApiKey: "the api key",
    CommerceCode: "the commerce code",
    Environment: transbank.IntegrationURL, //Options: transbank.IntegrationURL for integration environment, transbank.ProductionURL for production environment.
}
tx := transbank.NewTransaction(options)
```
With that you can use the service methods displayed on the [official Transbank documentation](https://www.transbankdevelopers.cl/documentacion/como_empezar)
## Use
### Webpay Plus
```go
import webpay "github.com/ppastene/transbank-sdk-go"

options := &transbank.Options{
    ApiKey: "597055555532",
    CommerceCode: "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
    Environment: transbank.IntegrationURL,
}

transaction := transbank.NewTransaction(options)
res, err := transaction.Create("buy_order", "session_id", "amount", "http://return-url.com")
res, err := transaction.Commit("token")
res, err := transaction.Status("token")
res, err := transaction.Refund("token", "amount")
res, err := transaction.Capture("token", "buy_order", "authorization_code", "amount") // Only in environments with differ option
```
### Webpay Plus Mall
```go
import webpay "github.com/ppastene/transbank-sdk-go"

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
res, err := mallTransaction.Create("buy_order", "session_id", "http://return-url.com", details)
res, err := mallTransaction.Commit("token")
res, err := mallTransaction.Status("token")
res, err := mallTransaction.Refund("token", "buy_order", "amount", "child_commerce_code")
res, err := mallTransaction.Capture("token", "child_commerce_code", "buy_order", "authorization_code", "amount") // Only in environments with differ option
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
res, err := oneclickInscrpition.Start("user", "email", "http://return-url.com")
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
res, err := oneclickTransaction.Capture("buy_order", "commerce_code", "authorization_code", "amount") // Only in environments with differ option
```
### Error management
Transbank SDK Go uses a struct to manage all the errors encountered
```go
type WebpayError struct {
	Code           int      
	ServiceMessage string   
	Cause          error
}
```
- Code:
    - -1: Validation errors happened on the SDK before the API call
    - 0: HTTP client errors
    -  1 and up: HTTP error code returned on the API response
- ServiceMessage: Descriptive error
- Cause: The error as-is

If you print the error it will use the Error() method which will return a string of ServiceMessage + Cause
```go
options := &transbank.Options{
    ApiKey: "the api key",
    CommerceCode: "the commerce code",
    Environment: transbank.IntegrationURL,
}
transaction := transbank.NewTransaction(options)

res, err := transaction.Status("")
if err != nil {
    fmt.Println(err)
}
/*
    It will print on console
    SDK Validation Error: 'token' cannot be empty
*/
```
You can use errors.Unwrap(err) to get the error directly
```go
options := &transbank.Options{
    ApiKey: "the api key",
    CommerceCode: "the commerce code",
    Environment: transbank.IntegrationURL,
}
transaction := transbank.NewTransaction(options)

res, err := transaction.Status("non-registered-token")
if err != nil {
    fmt.Println(errors.Unwrap(err))
}
/*
    It will print on console
    Invalid value for parameter: token
*/
```
## HTTP Client injection
The SDK includes a HTTP client to communicate with the Transbank API. You can replace that client by injecting a HTTP client of your choice.

These are the functions of the services you can inject a HTTP client
```go
NewTransactionWithClient(client shared.HTTPClientInterface, opt *shared.Options) // For Webpay Plus
NewMallTransactionWithClient(client shared.HTTPClientInterface, opt *shared.Options) // For Webpay Plus Mall
NewOneclickMallInscriptionWithClient(client shared.HTTPClientInterface, opt *shared.Options) // For Oneclick Inscripions
NewOneclickMallTransactionWithClient(client shared.HTTPClientInterface, opt *shared.Options) // For Oneclick Transactions
```
The interface contains the following methhod
```go
type HTTPClientInterface interface {
	Request(method string, url string, headers map[string]string, payload any) ([]byte, int, error)
}
```
The parameters are:
- method: The HTTP method. Must be written uppercase
- url: The url to call
- headers: The request hearders
- payload: The request payload. Must be capable to be converted to JSON using json.Marshal(). It can be empty.

The responses returned by the interface are:
- []byte: The URL raw response.
- int: The HTTP code of the response
- error: The error in case of problems with the DNS, hostname, communication, unmarshalling
### Example using Resty client
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
## Implementation Example
In [this repository](https://github.com/ppastene/transbank-sdk-go-example) you will find an implementation example of the SDK using Goravel. Follow the README instructions and the [Goravel documentation](https://www.goravel.dev/getting-started/installation.html) for further information.

# Transbank SDK Go
Library to connect with the Transbank API written in the Go language

[Spanish](./README.md) | **English**
## Index
- [Requirements](#requirements)
- [Installation](#installation)
- [First steps](#first-steps)
- [Use](#use)
  - [Webpay Plus](#webpay-plus)
  - [Webpay Plus Mall](#webpay-plus-mall)
- [Error management](#error-management)
- [Examples](#examples)
    - [HTTP Client injection](#http-client-injection)
        - [Example using Resty client](#example-using-resty-client)
- [Roadmap](#roadmap)
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
Declare the function of the service you wanna use along with the credentials
```go
options := &webpay.Options{
    ApiKey: "su api key",
    CommerceCode: "su codigo de comercio",
    Environment: webpay.IntegrationURL, //Options: webpay.IntegrationURL for integration environment, webpay.ProductionURL for production environment.
}
tx := webpay.NewTransaction(options)
```
With that you can use the service methods displayed on the [official Transbank documentation](https://www.transbankdevelopers.cl/documentacion/como_empezar)
## Use
### Webpay Plus
```go
options := &webpay.Options{
    ApiKey: "597055555532",
    CommerceCode: "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
    Environment: webpay.IntegrationURL,
}
webpayplus := webpay.NewTransaction(options)
details := []webpay.MallDetail{

}
res, err := webpayplus.Create("buy_order", "session_id", "amount", "http://url-de-retorno.cl")
res, err := webpayplus.Commit("token")
res, err := webpayplus.Status("token")
res, err := webpayplus.Refund("token", "amount")
res, err := webpayplus.Capture("token", "buy_order", "authorization_code", "amount") // Only in environments with differ option
```
### Webpay Plus Mall
```go
options := &webpay.Options{
    ApiKey: "597055555532",
    CommerceCode: "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
    Environment: webpay.IntegrationURL,
}
webpayplus := webpay.NewMallTransaction(options)
details := webpay.[]MallDetail{
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
res, err := webpayplus.Create("buy_order", "session_id", "http://url-de-retorno.cl", details)
res, err := webpayplus.Commit("token")
res, err := webpayplus.Status("token")
res, err := webpayplus.Refund("token", "buy_order", "amount", "child_commerce_code")
res, err := webpayplus.Capture("token", "child_commerce_code", "orden_compra", "authorization_code", "amount") // Only in environments with differ option
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
res, err := webpayplus.Status("")
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
res, err := webpayplus.Status("non-registered-token")
if err != nil {
    fmt.Println(errors.Unwrap(err))
}
/*
    It will print on console
    Invalid value for parameter: token
*/
```
## Examples
### HTTP Client injection
The SDK includes a HTTP client to communicate with the Transbank API. You can replace that client by injecting a HTTP client of your choice.

These are the functions of the services you can inject a HTTP client
```go
func NewTransactionWithClient(client shared.HTTPClientInterface, opt *shared.Options)
func NewMallTransactionWithClient(client shared.HTTPClientInterface, opt *shared.Options)
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
- payload: The request payload. Must be capable to be converted o JSON using json.Marshal(). It can be empty.

Las respuestas que devuelve el metodo son:
- []byte: The URL raw response.
- int: The HTTP code of the response
- error: The error in case of problems with the DNS, hostname, communication, unmarshalling
#### Example using Resty client
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
	options := &webpay.Options{
		ApiKey:       "597055555532",
		CommerceCode: "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
		Environment:  webpay.IntegrationURL,
	}
	tx := webpay.NewTransactionWithClient(client, options)
	resp, err := tx.Status("token")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}
```
## Roadmap
To see what is yet to implement check the [TODO.md](./TODO.md) file.
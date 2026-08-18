# Testing del SDK

Cómo se testea el SDK de Transbank sin consumir la API real.

## Objetivo

Los tests deben verificar el comportamiento del SDK (headers, method, path, body,
decode de respuestas, errores) **sin tocar la red**. Para lograrlo se usa
`net/http/httptest` de la stdlib, cero dependencias externas.

## Cómo funciona la inyección HTTP

El SDK arma la URL de los requests desde el `Environment`, no desde el cliente:

- `internal/config.go` — por defecto `&http.Client{Timeout: 30 * time.Second}`.
- Cada constructor de servicio (webpayplus, oneclick, transaccioncompleta) setea `Config.Headers` con `Tbk-Api-Key-Id`/`Tbk-Api-Key-Secret`; PatPass usa `Commercecode`/`Authorization`.
- `webpayplus/transaction.go:28-32` — si `Options.HTTPClient != nil`, lo reemplaza. Ese es el punto de inyección. `Options.HTTPClient` es la interfaz `transbank.HTTPClient` (método `Do`), así que acepta cualquier implementación custom (no solo `*http.Client`).
- `internal/requestor.go:64` — `r.config.HTTP.Do(req)` ejecuta el request con ese cliente. Los headers salen de `config.Headers`.

La idea de los tests: inyectar un cliente que intercepte cualquier request (aunque
el host diga `webpay3gint`) y lo mande a un mock local.

## El gotcha: `httptest.Server.Client()`

Suposición inicial (incorrecta): `httptest.Server.Client()` devuelve un cliente que
enruta **cualquier URL** al mock. Es falso en la práctica:

- `GET http://127.0.0.1:puerto/direct` → sí llega al mock (200).
- `GET https://webpay3gint.transbank.cl/...` → **no** llega al mock; el request
  sale a internet y golpea la API real (responde 401/400/422).

Consecuencia: un request a `https://webpay3gint.transbank.cl/...` con ese cliente
termina en la API real de Transbank. Por eso durante una fase de debug los tests
fallaban con mensajes reales como `"Commerce is CHILEAN_PESO so decimal amounts are
not allowed"`.

## La solución: un RoundTripper propio

En `webpayplus/mock_test.go`, `mockTransport` reescribe la URL del request al mock:

1. Toma el request entrante (sea cual sea el host, ej. `webpay3gint`).
2. Reescribe `scheme` + `host` a la dirección del mock (`127.0.0.1:puerto`).
3. Delega en `http.DefaultTransport.RoundTrip`.

Como `Options.HTTPClient = mock.Client()`, el requestor hace `client.Do(req)` → el
transport lo reescribe a localhost → el handler del mock registra method/path/
headers/body y devuelve el JSON canjeado. El request jamás sale a la red.

## Archivos de tests

| Archivo | Cubre |
|---|---|
| `webpayplus/mock_test.go` | Helper `mockServer` (httptest + mockTransport), asserts de requests |
| `webpayplus/transaction_test.go` | 5 métodos de Transaction: headers, path, body, errores HTTP, validación sin llamar a la API |
| `webpayplus/malltransaction_test.go` | 5 métodos de MallTransaction, validación de `MallDetails` |
| `webpayplus/transaction_responses_test.go` | `IsApproved()` en los 3 niveles |
| `oneclick/mock_test.go` | Helper `mockServer` (httptest + mockTransport), asserts de requests |
| `oneclick/inscription_test.go` | 3 métodos de MallInscription: headers, path, body, errores HTTP, validación sin llamar a la API |
| `oneclick/transaction_test.go` | 4 métodos de MallTransaction (incluye verificación de que `capture` usa PUT), validación de `MallDetails` |
| `oneclick/responses_test.go` | `IsApproved()` en los 2 niveles |
| `transaccioncompleta/mock_test.go` | Helper `mockServer` (httptest + mockTransport), asserts de requests incl. body exacto (`assertMallRequest` valida `Tbk-Api-Key-Id` del comercio mall) |
| `transaccioncompleta/transaction_test.go` | 6 métodos de Transaction: headers, path, body exacto, errores HTTP, validación sin llamar a la API |
| `transaccioncompleta/responses_test.go` | Pruebas estrictas de tipos: installments/deferred_periods como números, `type` en refund, `IsApproved()`, fallo documentado si la API devuelve strings |
| `transaccioncompleta/malltransaction_test.go` | 6 métodos de MallTransaction: body exacto (installments como array de tiendas, commit que omite los campos opcionales sin serializar `null`), validación sin llamar a la API |
| `transaccioncompleta/malltransaction_responses_test.go` | Decode de installments/commit/status/refund/capture mall, `IsApproved()` a nivel detalle y transacción |
| `patpass/mock_test.go` | Helper `mockServer` (httptest + mockTransport), asserts de requests (valida headers `Commercecode`/`Authorization`, y que `Tbk-Api-Key-Id` NO esté presente) |
| `patpass/inscription_test.go` | `NewInscription` (validación de credenciales), `Start` (body exacto de las 16 claves, con `commerceCode` dinámico), `Status`, validaciones sin llamar a la API, errores HTTP |
| `patpass/responses_test.go` | Decode de start/status, `Authorized` como booleano, fallo documentado si la API devuelve el campo como string |
| `internal/validation_test.go` | Validadores (token, buy_order, session_id, amount, return_url, commerce_code, username, email, tbk_user, card_number, card_expiration_date, cvv, installments_number) |
| `options_test.go` | `Options.Validate()` + `Error.Error()`/`Unwrap()` |

## Tests de integración (API real)

Los tests mock no tocan la red. Para verificar el comportamiento real del SDK
contra el sandbox de integración de Transbank, cada paquete tiene un
`integration_test.go` opcional que se ejecuta solo cuando `TBK_INTEGRATION=1`:

```sh
TBK_INTEGRATION=1 go test ./...
```

Cubren el flujo real que describe la documentación de Transbank:

- **`create`/`start` siempre responde**: se aserta que devuelve token y URL no
  vacíos. Webpay Plus (`Create`), Transacción Completa (`Create`, tarjeta de
  test `4051885600446623`), OneClick (`Start`) y PatPass (`Start`).
- **Los servicios dependientes consumen el token retornado** (`Commit`,
  `Status`, `Refund`, `Capture`, `Finish`, `Delete`, `Authorize`,
  `Installments`). El token depende de que el cliente complete el flujo en el
  portal, así que sin ese paso la API responde **o** con JSON de propiedades
  vacías **o** con un mensaje de error.
- El helper `assertServed` acepta ambos resultados (respuesta decodificada con
  campos vacíos o `*transbank.HTTPError`) y falla solo si el SDK paniquea o
  devuelve un tipo de error inesperado.

## Nota de transparencia

Durante el debug de la fase intermedia se consumió el sandbox de integración de
Transbank (`webpay3gint.transbank.cl`, no producción), con credenciales de test
oficiales y tokens inválidos. No debería tener efecto alguno, pero queda registrado.
La suite final **no toca la red**: todos los requests quedan en localhost salvo
que se ejecuten los tests de integración con `TBK_INTEGRATION=1`.

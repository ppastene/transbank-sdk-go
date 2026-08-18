# Security Audit: transbank-sdk-go

Fecha: 2026-08-18

## Resumen Ejecutivo

El SDK tiene una postura de seguridad aceptable para un cliente HTTP. Cero dependencias externas, validación de inputs completa, y URLs HTTPS por defecto. Las vulnerabilidades encontradas son de severidad media y baja, ninguna crítica.

## Hallazgos MEDIUM

- [ ] **1. `io.ReadAll` sin límite en cuerpo de error HTTP**
  - **Archivo:** `internal/requestor.go:70`
  - **Problema:** Un servidor malicioso podría enviar un cuerpo gigante causando OOM (DoS). El error de `ReadAll` también se descarta con `_`.
  - **Solución:** Envolver `resp.Body` con `io.LimitReader(resp.Body, 1<<20)` (1 MB) antes de `ReadAll`. Manejar el error de `ReadAll`.

- [ ] **2. URL construida por concatenación de strings**
  - **Archivo:** `internal/requestor.go:42`
  - **Problema:** Sin uso de `url.Parse()` o `url.Join()`. Si `BaseURL` no termina con `/` o `path` no empieza con `/`, la URL resultante podría ser inválida. PatPass funciona "por coincidencia" (su base URL tiene `/` al final).
  - **Solución:** Usar `url.Join()` o `url.Parse()` para construir el endpoint de forma segura.

- [ ] **3. `HTTPClient` inyectado sin validación de timeout**
  - **Archivo:** `options.go:18-27`
  - **Problema:** Un usuario podría inyectar `http.DefaultClient` (sin timeout), vulnerable a slow-loris DoS.
  - **Solución:** Documentar en godoc que el `HTTPClient` inyectado debe tener un timeout razonable. Opcionalmente, si el cliente es `*http.Client`, verificar que tenga timeout configurado.

- [ ] **4. Credenciales en headers HTTP en texto plano**
  - **Archivo:** Todos los constructores `New*()`
  - **Problema:** Si un usuario inyecta un HTTPClient con `InsecureSkipVerify: true` o usa una URL `http://`, las API keys viajarían sin cifrar. No hay advertencia documentada.
  - **Solución:** Agregar advertencia en la documentación de `Options.HTTPClient` sobre el uso de TLS. Opcionalmente, validar que las URLs por defecto usen HTTPS.

## Hallazgos LOW

- [ ] **5. `HTTPError.Body` expone respuesta completa del servidor**
  - **Archivo:** `error.go:39-43`
  - **Problema:** Si se loguea el error, información sensible del servidor de Transbank se filtra en logs.
  - **Solución:** Truncar el body almacenado en `HTTPError` a un máximo razonable (ej. 512 bytes), o proporcionar un método para acceder al body completo separado del string de `Error()`.

- [ ] **6. `ValidateToken` solo verifica largo, no contenido**
  - **Archivo:** `internal/validation.go:32-37`
  - **Problema:** Un token de 64 caracteres con `/`, `?` o `#` pasaría validación y causaría path traversal en URLs construidas con `fmt.Sprintf`.
  - **Solución:** Agregar regex `^[A-Za-z0-9]+$` además de verificar largo.

- [ ] **7. `ValidateEmail` es mínimo**
  - **Archivo:** `internal/validation.go:82-93`
  - **Problema:** Solo verifica presencia de `@`. Strings como `@@` o `a@b` pasan.
  - **Solución:** Usar una regex más robusta o una librería de validación de emails.

- [ ] **8. `ValidateSessionID` no tiene whitelist de caracteres**
  - **Archivo:** `internal/validation.go:52-59`
  - **Problema:** Acepta cualquier carácter Unicode, a diferencia de `BuyOrder` que tiene regex.
  - **Solución:** Agregar un patrón regex similar al de `BuyOrder`.

- [ ] **9. `TransportError.Error()` descarta la causa raíz**
  - **Archivo:** `error.go:25-27`
  - **Problema:** Solo retorna `Message`, no incluye `Err`. Dificulta debugging en producción.
  - **Solución:** Cambiar a `return fmt.Sprintf("%s: %v", e.Message, e.Err)`.

- [ ] **10. `io.ReadAll` error descartado**
  - **Archivo:** `internal/requestor.go:70`
  - **Problema:** Si la lectura del body falla, `HTTPError.Body` queda vacío perdiendo información del servidor.
  - **Solución:** Manejar el error de `io.ReadAll` y propagarlo como `TransportError` o incluir un mensaje por defecto.

## Hallazgos POSITIVOS

- Cero dependencias externas — Elimina riesgo de supply chain.
- Validación de inputs exhaustiva — 12 funciones de validación con límites estrictos.
- URLs HTTPS por defecto — Todos los endpoints hardcoded usan HTTPS.
- Timeout de 30s en cliente por defecto — Protege contra slow-loris.
- Sin `panic()`, `fmt.Print` ni logging en código producción.
- Response bodies cerrados correctamente con `defer resp.Body.Close()`.
- Sin race conditions detectadas en código compartido.
- Credenciales de sandbox hardcoded solo en archivos de test, nunca en producción.
- `os.Getenv` solo usado en tests de integración (gated por `TBK_INTEGRATION`).

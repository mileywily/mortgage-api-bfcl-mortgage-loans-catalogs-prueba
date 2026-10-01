
# `http_fif`

## Descripción

`http_fif` es un paquete de Go que proporciona un cliente HTTP con características adicionales como logging y tracing. Simplifica las solicitudes HTTP mediante métodos para los verbos HTTP comunes y el manejo de errores estándar.

## Características

- **HTTP Client**: Un cliente HTTP personalizable con soporte para configurar encabezados, URL base y tiempo de espera.
- **REST Client**: Un cliente REST personalizable con soporte para configurar encabezados, URL base y tiempo de espera.
- **Logging**: Integración de logs utilizando `slog`.
- **Tracing**: Soporte para tracing distribuido utilizando `dd-trace-go` de DataDog.
- **Headers**: Soporte para el reenvio de cabeceras obligatorias
- **Authorization**: Soporte para Autorización
- **Manejo de Errores**: Errores predefinidos para problemas comunes como URLs inválidas, solicitudes nulas y parámetros incorrectos.

## Instalación

Para instalar el paquete, utiliza el siguiente comando:

```sh
go get -u https://github.com/falabella-regulado/go-lib-http-fif
```

## Uso

### Crear un Cliente HTTP

Puedes crear un nuevo cliente HTTP con opciones predeterminadas o personalizarlo usando `HttpClientOption`.

```go
import (
    "context"
    "github.com/yourusername/http_fif"
    "time"
    "log"
    "slog"
)

func main() {
    client := http_fif.NewHTTPClient(
        http_fif.BaseURL("https://api.example.com"),
        http_fif.Timeout(30 * time.Second),
        http_fif.Logger(slog.Default()),
    )

    // Ejemplo de solicitud GET
    req, _ := http.Request("GET", "/", nil)
    response, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()
}
```

### Realizar Solicitudes

El cliente proporciona un metodo para ejecutar un request.

#### Do Request

```go
req, _ := http.Request("GET", "/", nil)
response, err := client.Do(req)
if err != nil {
    log.Fatal(err)
}
defer response.Body.Close()
```

### Opciones del cliente HTTP

las opciones para un cliente HTTP son:

#### BaseURL
```go
client := http_fif.NewHTTPClient(
    http_fif.BaseURL("https://api.example.com"),
)
```

#### Header
```go
client := http_fif.NewHTTPClient(
    http_fif.Header("header-key", "header-value"),
)
```

#### Timeout
```go
client := http_fif.NewHTTPClient(
    http_fif.Timeout(5 * time.Millisecond),
)
```

#### TLSConfig
```go
client := http_fif.NewHTTPClient(
    http_fif.TLSConfig(&tls.Config{InsecureSkipVerify: true}),
)
```

#### Logger
```go
import "github.com/falabella-regulado/go-lib-logger-fif/loggers"
...
client := http_fif.NewHTTPClient(
    http_fif.Logger(loggers.NewJsonSlogLogger(loggerFif.InfoLevel, nil)),
)
...
```

#### WithTokenProvider
```go
provider := http_fif.NewBearerAuthorizationProvider(.....)
client := http_fif.NewHTTPClient(
    http_fif.WithTokenProvider(provider, "Authorization"),
)

```

#### WithDatadogTracer
```go
client := http_fif.NewHTTPClient(
    http_fif.WithDatadogTracer(),
)

```

#### WithMandatoryHeaders
```go
client := http_fif.NewHTTPClient(
    http_fif.WithMandatoryHeaders(),
)

```

### Crear un Cliente REST

Puedes crear un nuevo cliente HTTP con opciones predeterminadas o personalizarlo usando `HttpClientOption` o `RestClientOption`.

```go
import (
    "context"
    "github.com/yourusername/http_fif"
    "time"
    "log"
    "slog"
)

func main() {
    client := http_fif.NewRestClient(
        http_fif.BaseURL("https://api.example.com"),
        http_fif.Timeout(30 * time.Second),
        http_fif.Logger(slog.Default()),
    )

    // Ejemplo de solicitud GET
    response, err := client.Get(context.Background(), "/path")
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()
}
```


### Realizar Solicitudes

El cliente proporciona métodos para realizar solicitudes `GET`, `POST`, `PUT` y `DELETE`.

#### GET Request

```go
response, err := client.Get(context.Background(), "/path", "param1", "value1")
if err != nil {
    log.Fatal(err)
}
defer response.Body.Close()
```

#### POST Request

```go
response, err := client.Post(context.Background(), "/path", requestBody)
if err != nil {
    log.Fatal(err)
}
defer response.Body.Close()
```

#### PUT Request

```go
response, err := client.Put(context.Background(), "/path", requestBody)
if err != nil {
    log.Fatal(err)
}
defer response.Body.Close()
```

#### DELETE Request

```go
response, err := client.Delete(context.Background(), "/path")
if err != nil {
    log.Fatal(err)
}
defer response.Body.Close()
```

### Opciones del cliente REST

Ademas de las opciones soportadas por el cliente HTTP, el cliente REST soporta las siguientes opciones

#### Encoder
```go
var encoder DataEncoder
...
client := http_fif.NewHTTPClient(
    http_fif.Encoder(encoder),
)
...
```

### Manejo de Errores

El paquete define varios errores para problemas comunes:

- `ErrInvalidURLAddress`: Se devuelve cuando la URL proporcionada es inválida.
- `ErrRequestNil`: Se devuelve cuando el objeto de la solicitud es nulo.
- `ErrResponseNil`: Se devuelve cuando el objeto de la respuesta es nulo.
- `ErrFailedValidation`: Se devuelve cuando el proceso de validación falla.
- `ErrNilBody`: Se devuelve cuando el cuerpo de la solicitud o la respuesta es nulo.
- `ErrInvalidBody`: Se devuelve cuando el cuerpo de la solicitud o la respuesta es inválido.
- `ErrInvalidParametersCount`: Se devuelve cuando la cantidad de parámetros proporcionados es incorrecta.

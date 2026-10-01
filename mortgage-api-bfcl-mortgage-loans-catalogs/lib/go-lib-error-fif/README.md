# error-fif

## Descripción

`error-fif` es un paquete de Go que proporciona un conjunto de errores y sus respectivas respuestas http.

## Características

- **Errors**: Conjunto de errores tipificados para cada caso
- **Interfaces**: Conjunto de interfaces para utilizar los errores
- **Error Handler**: implementacion de la interface para manejar errores
- **Response Interceptor**: implementacion de la interface para agregar trace-id and date-time a la respuesta de error
- **Nice Recovery**: implementacion de la interface para manejar el panic y retornar una respuesta de error

## Instalación

Para instalar el paquete, utiliza el siguiente comando:

```sh
go get -u https://github.com/falabella-regulado/go-lib-error-fif
```

## Uso

### Errores que ofrece la libreria

```go
import (
    errorFif "github.com/falabella-regulado/go-lib-error-fif"
)
func main() {
    errBadRequest := ErrBadRequest{Code: 400, Message: "Mensaje del error", Resource: "Recurso"}
    errUnauthorized := ErrUnauthorized{Code: 401, Message: "Mensaje del error", Resource: "Recurso"}
    errForbidden := ErrForbidden{Code: 403, Message: "Mensaje del error", Resource: "Recurso"}
    errNotFound := ErrNotFound{Code: 404, Message: "Mensaje del error", Resource: "Recurso"}
    errNotAcceptable := ErrNotAcceptable{Code: 406, Message: "Mensaje del error", Resource: "Recurso"}
    errUnprocessableEntity := ErrUnprocessableEntity{Code: 422, Message: "Mensaje del error", Resource: "Recurso"}
    errPreconditionRequired := ErrPreconditionRequired{Code: 428, Message: "Mensaje del error", Resource: "Recurso"}
    errInternalServerError := ErrInternalServerError{Code: 500, Message: "Mensaje del error", Resource: "Recurso"}
    errBadGateway := ErrBadGateway{Code: 502, Message: "Mensaje del error", Resource: "Recurso"}
    errTimeout := ErrTimeout{Code: 504, Message: "Mensaje del error", Resource: "Recurso"}
}
```

### Interface ErrorHandler
```go
import (
    errorFif "github.com/falabella-regulado/go-lib-error-fif"
)

//type ErrorHandler func(error) (int, interface{})
errorFif.ErrorHandler
```

### Crear un nuevo ErrorHandler por defecto

```go
import (
    errorFif "github.com/falabella-regulado/go-lib-error-fif"
)

func main() {
    ...
    errorHandler := errorFif.DefaultErrorHandler()
    ...
}
```

### Crear un nuevo ErrorHandler con los errores FIF

```go
import (
    errorFif "github.com/falabella-regulado/go-lib-error-fif"
)

func main() {
    ...
    errorHandler := errorFif.DefaultErrorHandler()
    errorHandler = errorFif.FifErrorHandler(errorHandler)
    ...
}
```

### Crear un nuevo ErrorHandler que maneje errores propios

```go
import (
    errorFif "github.com/falabella-regulado/go-lib-error-fif"
)

var customErr = errors.new("esto es un error_handler custom")

func NewCustomErrorToHttp(next errorFif.ErrorHandler) errorFif.ErrorHandler {
    return func(err error) (int, interface{}) {
        switch {
        case errors.Is(err, &customErr):
            return http.StatusInternalServerError, errorFif.ResponseBody{
                Code:    http.StatusInternalServerError,
                Message: err.Error(),
                Status:  errorFif.INTERNAL_ERROR,
            }
        }
        return next(err)
    }
}

func main() {
    ...
    errorHandler := errorFif.DefaultErrorHandler()
    errorHandler = NewCustomErrorToHttp(errorHandler)
    ...
}
```

### Interface ResponseInterceptor
```go
import (
    errorFif "github.com/falabella-regulado/go-lib-error-fif"
)

//type ResponseInterceptor func(*http.Request, int, interface{}) interface{}
errorFif.ResponseInterceptor
```

### Crear un nuevo Interceptor que agregue trace-id y date-time

```go
import (
    errorFif "github.com/falabella-regulado/go-lib-error-fif"
)

func main() {
    ...
    errorInterceptor := errorFif.ContextResponseInterceptor()
    ...
}
```

### Crear un nuevo Interceptor que modifique el body con una funcion custom

```go
import (
    errorFif "github.com/falabella-regulado/go-lib-error-fif"
)

func CustomBodyModifier() errorFif.ResponseInterceptor {
    return func(r *http.Request, status int, b interface{}) interface{} {
        ...
		return struct{}
    }
}

func main() {
	errorInterceptor := CustomBodyModifier()
}
```

### Interface PanicRecoveryHandler
```go
import (
    errorFif "github.com/falabella-regulado/go-lib-error-fif"
)

//type PanicRecoveryHandler func(interface{}) (int, interface{})
errorFif.PanicRecoveryHandler
```

### Usar el NicePanicRecovery

```go
import (
    errorFif "github.com/falabella-regulado/go-lib-error-fif"
)

func main() {
    ...
    panicRecovery := errorFif.FifPanicRecoveryHandler()
    ...
}
```

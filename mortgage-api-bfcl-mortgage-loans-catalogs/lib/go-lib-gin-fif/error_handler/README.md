# Error Handler Middleware

## Descripción

`error-handler` es un paquete de gin-fif que proporciona un middleware para manejar errores.

## Instalación

Para instalar el paquete, utiliza el siguiente comando:

```sh
go get -u https://github.com/falabella-regulado/go-lib-gin-fif
```

Añade la dependencia de `gin-gonic` a tu proyecto:

```sh
go get -u github.com/gin-gonic/gin
```

## Uso

un middleware por defecto

```go
import (
    "github.com/gin-gonic/gin"
    errorHandler "github.com/falabella-regulado/go-lib-gin-fif/error_handler"
)
func main() {
    engine := gin.New()

    // Agregar el middleware a Gin
    engine.Use(errorHandler.NewDefaultErrorHandlerMiddleware())

    // Definir rutas
    r.GET("/example", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "success"})
    })

    // Ejecutar el servidor
    r.Run(":8080")
}
```

un error handler customizado

```go
import (
    "github.com/gin-gonic/gin"
    errorHandler "github.com/falabella-regulado/go-lib-gin-fif/error_handler"
)
func main() {
    engine := gin.New()
	
	// Definir una funcion para manejar errores
    h := func(err error) (int, interface{}) {
        return http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()}
    }

    // Agregar el middleware a Gin
    engine.Use(errorHandler.NewErrorHandlerMiddleware(h))

    // Definir rutas
    r.GET("/example", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "success"})
    })

    // Ejecutar el servidor
    r.Run(":8080")
}
```

el error handler de fif

```go
import (
    "github.com/gin-gonic/gin"
    errorHandler "github.com/falabella-regulado/go-lib-gin-fif/error_handler"
)

func main() {
    engine := gin.New()
	
    // Agregar el middleware a Gin
    engine.Use(errorHandler.NewFifErrorHandlerMiddleware())

    // Definir rutas
    r.GET("/example", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "success"})
    })

    // Ejecutar el servidor
    r.Run(":8080")
}
```

un middleware customizado utilizando los errores de fif

```go
import (
    "github.com/gin-gonic/gin"
    errorHandler "github.com/falabella-regulado/go-lib-gin-fif/error_handler"
)

type RepositoryError struct {
    Message string
    Cause   error
}

func (e *RepositoryError) Error() string {
    return e.Message
}

func main() {
    engine := gin.New()
	
	// Definir una funcion para manejar errores
    h := func(ctx *gin.Context, err error) (int, interface{}) {
        switch (err).(type) {
        case *domain.RepositoryError:
            newErr := err.(*domain.RepositoryError)
            return http.StatusBadRequest, errorFIF.ErrBadRequest{Message: newErr.Cause.Error(), Resource: "repostory"}
        default:
            return next(ctx, err)
        }
    }
	
    // Agregar el middleware a Gin
    engine.Use(errorHandler.NewFifErrorHandlerMiddleware(h))

    // Definir rutas
    r.GET("/example", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "success"})
    })

    // Ejecutar el servidor
    r.Run(":8080")
}
```
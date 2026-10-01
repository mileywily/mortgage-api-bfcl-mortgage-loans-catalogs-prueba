# Panic Recovery Middleware

## Descripción

`panic-recovery` es un paquete de gin-fif que proporciona un middleware para manejar panic's.

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
    panicRecovery "github.com/falabella-regulado/go-lib-gin-fif/panic_recovery"
)
func main() {
    engine := gin.New()

    // Agregar el middleware a Gin
    engine.Use(panicRecovery.NewDefaultPanicRecoveryMiddleware())

    // Definir rutas
    r.GET("/example", func(c *gin.Context) {
        panic("fail")
    })

    // Ejecutar el servidor
    r.Run(":8080")
}
```

con un middleware customizado

```go
import (
    "github.com/gin-gonic/gin"
    panicRecovery "github.com/falabella-regulado/go-lib-gin-fif/panic_recovery"
)
func main() {
    engine := gin.New()
	
	// Definir una funcion para manejar panic
    h := func(recovered interface{}) (int, interface{}) {
        return http.StatusInternalServerError, fmt.Sprintf("%s", recovered)
    }

    // Agregar el middleware a Gin
    engine.Use(panicRecovery.NewPanicRecoveryMiddleware(h))

    // Definir rutas
    r.GET("/example", func(c *gin.Context) {
        panic("fail")
    })

    // Ejecutar el servidor
    r.Run(":8080")
}
```


con el panic recovery de error fif

```go
import (
    "github.com/gin-gonic/gin"
    panicRecovery "github.com/falabella-regulado/go-lib-gin-fif/panic_recovery"
)
func main() {
    engine := gin.New()
	
    // Agregar el middleware a Gin
    engine.Use(panicRecovery.NewFifErrorPanicRecoveryMiddleware())

    // Definir rutas
    r.GET("/example", func(c *gin.Context) {
        panic("fail")
    })

    // Ejecutar el servidor
    r.Run(":8080")
}
```

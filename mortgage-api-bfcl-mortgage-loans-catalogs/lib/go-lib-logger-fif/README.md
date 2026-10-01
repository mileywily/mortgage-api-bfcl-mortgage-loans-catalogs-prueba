
# Middleware HTTP para Gin-Gonic con Logging y Tracing

## Descripción

Este middleware proporciona un sistema de logging y tracing avanzado para aplicaciones construidas con el framework `Gin-Gonic`. Está diseñado para capturar información detallada de las solicitudes y respuestas, así como integrarse con DataDog para tracing distribuido.

## Características

- **Logging**: Captura detalles de cada solicitud y respuesta, incluyendo encabezados, cuerpo y otros atributos configurables.
- **Tracing**: Se integra con `dd-trace-go` de DataDog para rastrear las solicitudes a través de diferentes servicios.
- **Filtros Personalizables**: Permite aplicar filtros para definir cuándo registrar o ignorar ciertas solicitudes.
- **Manejo de Cuerpos de Solicitud y Respuesta**: Limita el tamaño de los cuerpos que se registran para evitar problemas de rendimiento.

## Instalación

Para instalar el paquete, utiliza el siguiente comando:

```sh
go get -u github.com/falabella-regulado/go-lib-logger-fif
```


Añade la dependencia de `dd-trace-go` y `gin-gonic` a tu proyecto:

```sh
go get -u gopkg.in/DataDog/dd-trace-go.v1
go get -u github.com/gin-gonic/gin
```

## Uso

### Configuración Básica

Para utilizar el middleware, primero debes configurarlo con un logger (`slog.Logger`) y una configuración (`Config`).

```go
import (
    "github.com/gin-gonic/gin"
     middlewareLogger "github.com/falabella-regulado/go-lib-logger-fif/middleware"
    "log/slog"
)

func main() {
    r := gin.Default()

    // Crear un logger predeterminado
    logger := slog.New(slog.NewJSONHandler(os.Stdout))

    // Configuración del middleware
    config := middleware.Config{
        DefaultLevel:     slog.LevelInfo,
        ClientErrorLevel: slog.LevelWarn,
        ServerErrorLevel: slog.LevelError,
        WithRequestID:     true,
        WithRequestBody:   true,
        WithResponseBody:  true,
    }

    // Agregar el middleware a Gin
    r.Use(middlewareLogger.NewWithConfig(logger, config))

    // Definir rutas
    r.GET("/example", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "success"})
    })

    // Ejecutar el servidor
    r.Run(":8080")
}
```

### Opciones de Configuración

El middleware permite una serie de configuraciones a través de la estructura `Config`:

- **`DefaultLevel`**: Nivel de logging predeterminado (`slog.LevelInfo`).
- **`ClientErrorLevel`**: Nivel de logging para errores del cliente (`slog.LevelWarn`).
- **`ServerErrorLevel`**: Nivel de logging para errores del servidor (`slog.LevelError`).
- **`WithRequestID`**: Habilita la inclusión del ID de solicitud en los logs.
- **`WithRequestBody`**: Habilita la captura del cuerpo de la solicitud.
- **`WithRequestHeader`**: Habilita la captura de los encabezados de la solicitud.
- **`WithResponseBody`**: Habilita la captura del cuerpo de la respuesta.

### Ejemplo de Filtros

Los filtros permiten controlar qué solicitudes se deben registrar. Los filtros son funciones que reciben el contexto de Gin (`gin.Context`) y devuelven un valor booleano.

```go
config.Filters = []middleware.Filter{
    func(c *gin.Context) bool {
        // Ignorar solicitudes a /health
        return c.Request.URL.Path != "/health"
    },
}
```

## Integración con DataDog

Para habilitar el tracing con DataDog, asegúrate de configurar el `tracer` de DataDog antes de iniciar el servidor:

```go
import "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

func main() {
    // Iniciar el tracer de DataDog
    tracer.Start()
    defer tracer.Stop()

    // Resto de la configuración de Gin y el middleware
}
```





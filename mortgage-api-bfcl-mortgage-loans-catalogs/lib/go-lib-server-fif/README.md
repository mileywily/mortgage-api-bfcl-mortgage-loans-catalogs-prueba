# Server-Fif

Libreria que provee un servidor http con graceful shutdown.

### Configuracion

| Variables de entorno aceptadas | Valor por defecto |                                Descripción                                 |
|--------------------------------|:-----------------:|:--------------------------------------------------------------------------:|
| PORT                           |       8080        |                    Puerto que escucha el servicio http                     |
| WAIT_TOLERANCE_SERVER          |        30         |                 Tolerancia de espera al cerrar el servidor                 |
| WRITE_TOLERANCE_SERVER         |        15         |       Tiempo máximo permitido para escribir una respuesta al cliente       |
| READ_TOLERANCE_SERVER          |        15         |         Tiempo máximo permitido para leer la solicitud del cliente         |
| IDLE_TOLERANCE_SERVER          |        60         | Tiempo máximo permitido para que una conexión permanezca inactiva|

### Modo de uso

Primero descargamos la libreria y sus dependencias utilizando el comando:
```
go get "github.com/falabella-regulado/go-lib-server-fif"
```

Para utilizarla se debe importar la librería:

```
import (
    "github.com/falabella-regulado/go-lib-server-fif"
)
```

Crear un ServerOption
```
type ServerOptions struct {
	Logger  loggerFif.Logger
	Handler http.Handler
}
```

Llamar al metodo StartServer
```
 _ = serverFif.StartServer(serverOption)
```

# Librerías Corporativas y Variables

El proyecto usa la estructura hexagonal de la plantilla. Ya están integrados `config-fif`, `logger-fif`, middlewares de `gin-fif` y Datadog para el servidor. Los errores HTTP se mapean localmente para conservar el contrato legado.

Los adaptadores Finnflow/Java aún usan `net/http` con timeout configurable. Antes de producción hay que validar con el equipo corporativo la integración saliente de Datadog y si corresponde reemplazarlos por `http-fif`; el cliente debe conservar GET + Basic Auth para Finnflow y POST + headers para Java.

## Dependencias Privadas

Las dependencias declaradas por la plantilla usan el namespace `github.com/falabella-regulado/*`. Configura el acceso privado en una máquina corporativa y autentica Git con el mecanismo aprobado por tu organización:

```powershell
go env -w 'GOPRIVATE=github.com/falabella-regulado/*'
go mod tidy
go test ./...
go build ./cmd/api
```

No almacenes tokens o secretos en el repositorio. Si la organización sirve estas librerías desde otro host, confirma primero el namespace correcto y actualiza `go.mod` con el equipo de plataforma.

## Variables de Ejecución

- `PORT`: puerto HTTP; default `8080`. Para comparar localmente con Java, usa `8082`.
- `DEFAULT_BACKEND`: `real`, `java` o `dummy`; default `real`.
- `FINNFLOW_URL`: URL base de Finnflow; requerida si el default es `real`.
- `FINNFLOW_KEY` y `FINNFLOW_SECRET`: credenciales Finnflow, inyectadas como secretos.
- `JAVA_LEGACY_URL`: URL base del proxy Java; requerida si el default es `java`.
- `TIMEOUT`: timeout en segundos; default `10`.
- `DD_AGENT_HOST`, `DD_AGENT_PORT`, `DD_SERVICE_NAME` y `DD_PROFILE_ENABLED`: configuración Datadog.

`X-Backend-Env` permite seleccionar un backend por request. Las cabeceras `X-Channel`, `X-Commerce` y `X-Transaction-ID` se propagan al proxy Java, pero no son rechazadas cuando faltan por compatibilidad.

## Pipeline y Backstage

El proyecto conserva los archivos de pipeline de la plantilla. Antes de desplegar, confirmar las variables CI/CD del registry/escáner y revisar en `catalog-info.yaml` que owner, system y enlaces Datadog correspondan al equipo de Créditos Hipotecarios.

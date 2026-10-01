# Mortgage Loan Catalogs

Microservicio Gin para consultar catálogos hipotecarios, adaptado a la arquitectura hexagonal de la plantilla corporativa FIF.

## Endpoints

- `GET /health`: estado del servicio.
- `POST /v1/bfcl/mortgage-loan/catalogs/:catalog`: consulta un catálogo.
- `GET /swagger/index.html`: documentación OpenAPI.

La ruta legacy acepta `X-Channel`, `X-Commerce` y `X-Transaction-ID`, y los propaga al backend Java. Para mantener paridad con el servicio anterior, el servidor no rechaza requests cuando esas cabeceras faltan.

`X-Backend-Env` permite seleccionar `real`, `java` o `dummy` para una request. Si no se envía, se usa `DEFAULT_BACKEND`.

## Configuración

Las variables se leen con `config-fif`; credenciales deben inyectarse como secretos del entorno y no guardarse en el repositorio.

| Variable | Uso | Valor por defecto |
|---|---|---|
| `PORT` | Puerto HTTP | `8080` |
| `DEFAULT_BACKEND` | Backend (`real`, `java`, `dummy`) | `real` |
| `FINNFLOW_URL` | URL base Finnflow; requerida para backend `real` | vacío |
| `FINNFLOW_KEY` | Client ID de Finnflow | vacío |
| `FINNFLOW_SECRET` | Client secret de Finnflow | vacío |
| `JAVA_LEGACY_URL` | URL base del proxy Java | vacío |
| `TIMEOUT` | Timeout HTTP en segundos | `10` |
| `APP_NAME` | Nombre del servicio | `mortgage-api-bfcl-mortgage-loans-catalogs` |
| `APP_ENV` / `APP_VERSION` | Ambiente y versión para Datadog | `dev` / `0.0.0` |
| `GIN_MODE` | Modo Gin | `DEBUG` |
| `LOGGING_LEVEL` | Nivel de log | `info` |
| `DD_AGENT_HOST` / `DD_AGENT_PORT` | Agente Datadog | `localhost` / `8126` |
| `DD_PROFILE_ENABLED` | Habilita profiler Datadog | `false` |

Si `DEFAULT_BACKEND=java`, `JAVA_LEGACY_URL` es obligatoria. El backend dummy no requiere URLs externas.

## Desarrollo

Requiere Go 1.26.6 y acceso corporativo a las dependencias privadas.

```powershell
go test ./...
go build ./cmd/api
```

Para ejecutar localmente con datos dummy:

```powershell
$env:DEFAULT_BACKEND = "dummy"
go run ./cmd/api
```

La configuración por defecto expone la API en el puerto `8080`.

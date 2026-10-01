# API Hello World - Hexagonal Template

Esta es una API de ejemplo basada en una **Arquitectura Hexagonal** (puertos y adaptadores), diseñada para ser utilizada como plantilla base en proyectos de Go dentro de la plataforma.

## Estructura del Proyecto

El proyecto sigue una estructura modular que separa las responsabilidades en capas claras:

- **`cmd/`**: Punto de entrada de la aplicación.
  - `api/main.go`: Configuración inicial, inyección de dependencias y arranque del servidor Gin.
  - `config/config.go`: Gestión de configuración mediante variables de entorno y flags.
- **`internal/app/api/`**: Capa de aplicación (Interfaces de entrada).
  - `handler/`: Manejadores de rutas HTTP (Gin).
  - `dto/`: Objetos de transferencia de datos (Data Transfer Objects).
  - `mapper/`: Transformadores de datos entre la capa de entrada y el dominio.
- **`internal/core/`**: Núcleo de la aplicación (Lógica de negocio).
  - `domain/`: Entidades y errores del dominio.
  - `ports/`: Interfaces que definen los contratos de entrada (`in`) y salida (`out`).
  - `usecases/`: Implementación de la lógica de negocio.
- **`internal/infra/`**: Adaptadores de infraestructura (Implementaciones externas).
  - `rest/external_service/`: Cliente para consumir servicios REST externos.
  - `tracerMiddleware/`: Middlewares para observabilidad (Datadog).
- **`ias/`**: Configuración de infraestructura como código (Dockerfile, etc.).

## Requisitos Previos

- Go 1.22 o superior.
- Acceso a las librerías privadas de `gitlab.falabella.tech`.

## Configuración

La aplicación se configura mediante variables de entorno. Los valores por defecto se definen en `cmd/config/config.go`:

| Variable | Descripción | Valor por Defecto |
|----------|-------------|-------------------|
| `APP_NAME` | Nombre de la API | `api-hello-world` |
| `GIN_MODE` | Modo de Gin (`DEBUG` o `RELEASE`) | `DEBUG` |
| `LOGGING_LEVEL` | Nivel de log (`info`, `debug`, `error`) | `info` |
| `URI_PREFIX` | Prefijo de las rutas | `/fifcl/v1` |
| `COUNTRY` | País de despliegue (`CL`) | `CL` |
| `EXTERNAL_SERVICE_URL` | URL del servicio externo Customer | (vacío) |
| `DD_PROFILE_ENABLED` | Habilitar Profiler de Datadog | `false` |

## Comandos Útiles

El proyecto incluye un `Makefile` para facilitar las tareas comunes:

- **Compilar la aplicación**:
  ```bash
  make build
  ```
  Genera el binario en `build/bin/dist`.

- **Ejecutar tests**:
  ```bash
  make test
  ```

- **Verificar cobertura**:
  ```bash
  make coverage
  ```

## Endpoints Principales

- **Health Check**: `GET /health`
- **Ejemplo**: `POST /fifcl/v1/anything/:something/example` (Requiere cabeceras obligatorias según configuración).

## Tecnologías Utilizadas

- **Framework Web**: [Gin Gonic](https://github.com/gin-gonic/gin).
- **Observabilidad**: Datadog (Tracer & Profiler).
- **Librerías Base**: `go-libs` de Falabella (error-fif, logger-fif, http-fif, etc.).
- **Arquitectura**: Hexagonal / Puertos y Adaptadores.

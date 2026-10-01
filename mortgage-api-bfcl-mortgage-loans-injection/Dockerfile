# ==========================================
# STAGE 1: Build & Dependencies
# ==========================================
FROM golang:1.26.6-alpine AS builder

# Instalar certificados CA y zonas horarias
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Configuración de entorno de compilación Go
ARG GOPRIVATE="github.com/falabella-regulado/*"
ENV GOPRIVATE=${GOPRIVATE}
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Descarga e instalación de dependencias
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copia del código fuente
COPY . .

# Compilación del binario estático optimizado
RUN go build \
    -ldflags="-w -s -X main.version=${APP_VERSION:-1.0.0}" \
    -o /bin/api \
    ./cmd/api

# ==========================================
# STAGE 2: Minimal Runtime (Nullplatform)
# ==========================================
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

# Copiar certificados y zonas horarias
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copiar binario ejecutable
COPY --from=builder /bin/api /app/api

# Variables de entorno predeterminadas
ENV PORT=8080 \
    GIN_MODE=release

# Puerto expuesto por el contenedor
EXPOSE 8080

# Usuario sin privilegios root para mayor seguridad
USER 65532:65532

ENTRYPOINT ["/app/api"]

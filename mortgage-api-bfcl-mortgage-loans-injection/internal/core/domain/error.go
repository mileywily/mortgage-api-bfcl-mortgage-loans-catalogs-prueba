package domain

import "fmt"

// ServiceError cubre fallas de integración que no dependen del status HTTP:
// mapping de request, transporte, body ilegible o respuesta no decodificable.
type ServiceError struct {
	Message string
	Cause   error
}

func (e *ServiceError) Error() string { return e.Message }

func (e *ServiceError) Unwrap() error { return e.Cause }

// ServiceUnavailableError representa un status >= 500 del servicio externo.
type ServiceUnavailableError struct {
	Message string
	Body    string
}

func (e *ServiceUnavailableError) Error() string { return e.Message }

// ServiceClientError representa un status 4xx del servicio externo: la
// petición que armó esta API fue rechazada, no es corregible por el cliente.
type ServiceClientError struct {
	StatusCode int
	Body       string
}

func (e *ServiceClientError) Error() string {
	return fmt.Sprintf("external service rejected the request with status %d", e.StatusCode)
}

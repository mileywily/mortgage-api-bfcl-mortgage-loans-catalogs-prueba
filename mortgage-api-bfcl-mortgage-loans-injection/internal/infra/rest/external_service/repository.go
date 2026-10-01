package external_service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	httpFIF "github.com/falabella-regulado/go-lib-http-fif"

	"api-hello-world/internal/core/domain"
	"api-hello-world/internal/core/ports/out"
)

const resource = "customer Information"

type externalServiceRepository struct {
	client         httpFIF.RestClient
	requestMapper  RequestMapper
	responseMapper ResponseMapper
}

func NewExternalServiceRepository(
	client httpFIF.RestClient,
	requestMapper RequestMapper,
	responseMapper ResponseMapper,
) out.ExternalServiceRepository {
	return &externalServiceRepository{
		client:         client,
		requestMapper:  requestMapper,
		responseMapper: responseMapper,
	}
}

func (r *externalServiceRepository) Do(ctx context.Context, entity domain.Entity) (domain.EntityOut, error) {
	// 1. Mapear dominio → request del servicio externo
	if _, err := r.requestMapper(entity); err != nil {
		return domain.EntityOut{}, &domain.ServiceError{
			Message: fmt.Sprintf("mapping request for %s failed", resource),
			Cause:   err,
		}
	}

	// 2. Headers propios de la llamada: ninguno. Los corporativos los propaga
	//    WithMandatoryHeaders() y el tracing WithDatadogTracer() (wire-up).
	// 3. Llamada HTTP
	res, err := r.client.Get(ctx, "/", nil)
	if err != nil {
		return domain.EntityOut{}, &domain.ServiceError{
			Message: fmt.Sprintf("request to %s failed", resource),
			Cause:   err,
		}
	}
	defer res.Body.Close()

	// 4. Leer el body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return domain.EntityOut{}, &domain.ServiceError{
			Message: fmt.Sprintf("can't read %s response", resource),
			Cause:   err,
		}
	}

	// 5. Status → errores tipados
	switch {
	case res.StatusCode >= http.StatusInternalServerError:
		return domain.EntityOut{}, &domain.ServiceUnavailableError{
			Message: fmt.Sprintf("%s returned status %d", resource, res.StatusCode),
			Body:    string(body),
		}
	case res.StatusCode >= http.StatusBadRequest:
		return domain.EntityOut{}, &domain.ServiceClientError{
			StatusCode: res.StatusCode,
			Body:       string(body),
		}
	case res.StatusCode >= http.StatusMultipleChoices:
		return domain.EntityOut{}, &domain.ServiceError{
			Message: fmt.Sprintf("%s returned unexpected status %d", resource, res.StatusCode),
		}
	}

	// 6. Decodificar
	var apiResponse ExternalServiceResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return domain.EntityOut{}, &domain.ServiceError{
			Message: fmt.Sprintf("can't decode %s response", resource),
			Cause:   err,
		}
	}

	// 7. Mapear response → dominio
	return r.responseMapper(apiResponse)
}

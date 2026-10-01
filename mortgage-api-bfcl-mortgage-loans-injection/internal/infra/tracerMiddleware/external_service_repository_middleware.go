package tracerMiddleware

import (
	"context"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/ext"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	loggerFif "github.com/falabella-regulado/go-lib-logger-fif"

	"api-hello-world/internal/core/domain"
	"api-hello-world/internal/core/ports/out"
)

type externalServiceRepositoryMiddleware struct {
	next   out.ExternalServiceRepository
	logger loggerFif.Logger
}

func NewExternalServiceRepositoryMiddleware(next out.ExternalServiceRepository, logger loggerFif.Logger) out.ExternalServiceRepository {
	return &externalServiceRepositoryMiddleware{next: next, logger: logger}
}

func (m *externalServiceRepositoryMiddleware) Do(ctx context.Context, entity domain.Entity) (res domain.EntityOut, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "rest.external_service.Do")
	defer func() {
		if err != nil {
			m.logger.Error("Do failed", "cause:", err.Error(),
				ext.LogKeyTraceID, span.Context().TraceID(),
				ext.LogKeySpanID, span.Context().SpanID())
			span.SetTag(statusTag, failed)
			span.Finish(tracer.WithError(err))
			return
		}
		span.SetTag(statusTag, passed)
		span.Finish()
	}()
	return m.next.Do(ctx, entity)
}

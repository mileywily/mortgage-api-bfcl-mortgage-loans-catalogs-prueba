package tracerMiddleware

import (
	"context"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/ext"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	loggerFif "github.com/falabella-regulado/go-lib-logger-fif"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/ports/out"
)

type catalogRepositoryMiddleware struct {
	next   out.CatalogRepository
	logger loggerFif.Logger
}

// NewCatalogRepositoryMiddleware decora un CatalogRepository con tracing DataDog y logging de errores.
func NewCatalogRepositoryMiddleware(next out.CatalogRepository, logger loggerFif.Logger) out.CatalogRepository {
	return &catalogRepositoryMiddleware{next: next, logger: logger}
}

func (m *catalogRepositoryMiddleware) GetCatalog(ctx context.Context, req domain.GetCatalogRequest) (res interface{}, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "rest.catalog_service.GetCatalog")
	defer func() {
		if err != nil {
			m.logger.Error("GetCatalog repository failed",
				"catalog", req.CatalogName,
				"cause", err.Error(),
				ext.LogKeyTraceID, span.Context().TraceID(),
				ext.LogKeySpanID, span.Context().SpanID(),
			)
			span.SetTag(statusTag, failed)
			span.Finish(tracer.WithError(err))
			return
		}
		span.SetTag(statusTag, passed)
		span.Finish()
	}()
	return m.next.GetCatalog(ctx, req)
}

package tracerMiddleware

import (
	"context"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/ext"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	loggerFif "github.com/falabella-regulado/go-lib-logger-fif"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/ports/in"
)

type catalogUsecaseMiddleware struct {
	next   in.CatalogUsecase
	logger loggerFif.Logger
}

// NewCatalogUsecaseMiddleware decora un CatalogUsecase con tracing DataDog y logging de errores.
func NewCatalogUsecaseMiddleware(next in.CatalogUsecase, logger loggerFif.Logger) in.CatalogUsecase {
	return &catalogUsecaseMiddleware{next: next, logger: logger}
}

func (m *catalogUsecaseMiddleware) GetCatalog(ctx context.Context, req domain.GetCatalogRequest) (res interface{}, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "usecase.catalog.GetCatalog")
	defer func() {
		if err != nil {
			m.logger.Error("GetCatalog failed",
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

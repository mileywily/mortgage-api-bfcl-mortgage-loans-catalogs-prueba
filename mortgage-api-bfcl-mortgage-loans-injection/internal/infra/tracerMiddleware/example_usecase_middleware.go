package tracerMiddleware

import (
	"context"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/ext"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	loggerFif "github.com/falabella-regulado/go-lib-logger-fif"

	"api-hello-world/internal/core/domain"
	"api-hello-world/internal/core/ports/in"
)

type exampleUsecaseMiddleware struct {
	next   in.ExampleUsecase
	logger loggerFif.Logger
}

func NewExampleUsecaseMiddleware(next in.ExampleUsecase, logger loggerFif.Logger) in.ExampleUsecase {
	return &exampleUsecaseMiddleware{next: next, logger: logger}
}

func (m *exampleUsecaseMiddleware) Execute(ctx context.Context, entity domain.Entity) (res domain.EntityOut, err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "usecase.example.Execute")
	defer func() {
		if err != nil {
			m.logger.Error("Execute failed", "cause:", err.Error(),
				ext.LogKeyTraceID, span.Context().TraceID(),
				ext.LogKeySpanID, span.Context().SpanID())
			span.SetTag(statusTag, failed)
			span.Finish(tracer.WithError(err))
			return
		}
		span.SetTag(statusTag, passed)
		span.Finish()
	}()
	return m.next.Execute(ctx, entity)
}

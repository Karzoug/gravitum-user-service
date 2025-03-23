package middleware

import (
	"context"
	"net/http"

	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"go.opentelemetry.io/otel/trace"

	"github.com/Karzoug/meower-common-go/trace/otlp"

	gen "github.com/Karzoug/gravitum-user-service/internal/delivery/http/gen/user/v1"
)

// Otel starts the otel tracing and stores the trace id in the context.
func Otel(tracer trace.Tracer) gen.StrictMiddlewareFunc {
	return func(f nethttp.StrictHTTPHandlerFunc, operationID string) nethttp.StrictHTTPHandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (response any, err error) {
			return f(otlp.InjectTracing(ctx, tracer), w, r, request)
		}
	}
}

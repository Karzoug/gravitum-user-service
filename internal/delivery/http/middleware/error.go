package middleware

import (
	"context"
	"errors"
	"net/http"
	"syscall"

	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"

	"github.com/Karzoug/meower-common-go/ucerr"

	gen "github.com/Karzoug/gravitum-user-service/internal/delivery/http/gen/user/v1"
	"github.com/Karzoug/gravitum-user-service/internal/delivery/http/response"
)

type httpStatus interface {
	HTTPStatus() (int, string)
}

// Error is a middleware that handle errors and logs them.
func Error(logger zerolog.Logger) gen.StrictMiddlewareFunc {
	return func(f nethttp.StrictHTTPHandlerFunc, operationID string) nethttp.StrictHTTPHandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (resp any, err error) {
			defer func() {
				if nil == err {
					return
				}

				var switchErr error
				switch e := err.(type) { //nolint:errorlint
				case httpStatus:
					code, msg := e.HTTPStatus()
					logOnlyServiceError(ctx, err, operationID, logger)
					switchErr = response.JSON(w, code, gen.ErrorResponse{Error: msg})
				default:
					switch {
					//  got network error -> it's ok, so ignore it
					case isNetworkError(e):

					// it's unknown (= untrusted) error,
					// log it and return internal server error
					default:
						logger.Error().
							Err(e).
							Ctx(ctx). // for trace_id
							Str("operation_id", operationID).
							Msg("error handler")
						switchErr = response.JSON(w,
							http.StatusInternalServerError,
							gen.ErrorResponse{Error: http.StatusText(http.StatusInternalServerError)},
						)
					}
				}

				// finally, if we can't write error response, log it, unless it's network error
				if switchErr != nil && !isNetworkError(switchErr) {
					logger.Error().Err(switchErr).Msg("error handler: couldn't write error response")
				}

				err = nil
			}()

			return f(ctx, w, r, request)
		}
	}
}

func logOnlyServiceError(ctx context.Context, err error, method string, logger zerolog.Logger) {
	var e ucerr.Error
	if !errors.As(err, &e) {
		return
	}

	var ev *zerolog.Event
	switch e.Code() {
	case codes.OK,
		codes.Canceled,
		codes.InvalidArgument,
		codes.DeadlineExceeded,
		codes.NotFound,
		codes.AlreadyExists,
		codes.PermissionDenied,
		codes.FailedPrecondition,
		codes.OutOfRange,
		codes.Unimplemented, codes.Unauthenticated:
		return

	case codes.ResourceExhausted, codes.Aborted:
		ev = logger.Warn()

	case codes.Internal, codes.Unavailable, codes.Unknown, codes.DataLoss:
		ev = logger.Error()
	}

	ev.Err(e.Unwrap()).
		Ctx(ctx). // for trace_id
		Str("method", method).
		Msg("error handler")
}

func isNetworkError(err error) bool {
	// Ignore syscall.EPIPE and syscall.ECONNRESET errors which occurs
	// when a write operation happens on the http.ResponseWriter that
	// has simultaneously been disconnected by the client (TCP
	// connections is broken). For instance, when large amounts of
	// data is being written or streamed to the client.
	// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	// https://gosamples.dev/broken-pipe/
	// https://gosamples.dev/connection-reset-by-peer/

	switch {
	case errors.Is(err, syscall.EPIPE):

		// Usually, you get the broken pipe error when you write to the connection after the
		// RST (TCP RST Flag) is sent.
		// The broken pipe is a TCP/IP error occurring when you write to a stream where the
		// other end (the peer) has closed the underlying connection. The first write to the
		// closed connection causes the peer to reply with an RST packet indicating that the
		// connection should be terminated immediately. The second write to the socket that
		// has already received the RST causes the broken pipe error.
		return true

	case errors.Is(err, syscall.ECONNRESET):

		// Usually, you get connection reset by peer error when you read from the
		// connection after the RST (TCP RST Flag) is sent.
		// The connection reset by peer is a TCP/IP error that occurs when the other end (peer)
		// has unexpectedly closed the connection. It happens when you send a packet from your
		// end, but the other end crashes and forcibly closes the connection with the RST
		// packet instead of the TCP FIN, which is used to close a connection under normal
		// circumstances.
		return true
	}

	return false
}

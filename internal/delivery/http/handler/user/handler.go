package user

import (
	"net/http"

	"github.com/rs/zerolog"

	gen "github.com/Karzoug/gravitum-user-service/internal/delivery/http/gen/user/v1"
	"github.com/Karzoug/gravitum-user-service/internal/delivery/http/handler/errfunc"
	"github.com/Karzoug/gravitum-user-service/internal/delivery/http/middleware"
	"github.com/Karzoug/gravitum-user-service/internal/service"
)

const baseURL = "/api/web/v1"

var _ gen.StrictServerInterface = handlers{}

func RoutesFunc(us service.UserService, logger zerolog.Logger) func(mux *http.ServeMux) {
	logger = logger.With().
		Str("component", "http server: user handlers").
		Logger()

	hdl := handlers{
		userService: us,
		logger:      logger,
	}

	return func(mux *http.ServeMux) {
		gen.HandlerWithOptions(
			gen.NewStrictHandlerWithOptions(hdl,
				[]gen.StrictMiddlewareFunc{
					middleware.Recover,
					middleware.Error(logger),
					middleware.Logger(logger),
				},
				gen.StrictHTTPServerOptions{
					RequestErrorHandlerFunc:  errfunc.JSONRequest(logger),
					ResponseErrorHandlerFunc: errfunc.JSONResponse(logger),
				}),
			gen.StdHTTPServerOptions{
				BaseURL:    baseURL,
				BaseRouter: mux,
			})
	}
}

type handlers struct {
	userService service.UserService
	logger      zerolog.Logger
}

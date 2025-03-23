package app

import (
	"context"
	"runtime"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/Karzoug/gravitum-user-service/pkg/metric/prom"
	"github.com/Karzoug/gravitum-user-service/pkg/postgresql"

	"github.com/Karzoug/gravitum-user-service/internal/config"
	healthHandlers "github.com/Karzoug/gravitum-user-service/internal/delivery/http/handler/health"
	userHttp "github.com/Karzoug/gravitum-user-service/internal/delivery/http/handler/user"
	httpServer "github.com/Karzoug/gravitum-user-service/internal/delivery/http/server"
	userRepo "github.com/Karzoug/gravitum-user-service/internal/repo/pg"
	"github.com/Karzoug/gravitum-user-service/internal/service"
	"github.com/Karzoug/gravitum-user-service/pkg/buildinfo"
)

const (
	serviceName     = "UserService"
	metricNamespace = "user_service"
	initTimeout     = 10 * time.Second
	shutdownTimeout = 10 * time.Second
)

var serviceVersion = buildinfo.Get().ServiceVersion

func Run(ctx context.Context, logger zerolog.Logger) error {
	cfg, err := env.ParseAs[config.Config]()
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(cfg.LogLevel)

	logger.Info().
		Int("GOMAXPROCS", runtime.GOMAXPROCS(0)).
		Str("log level", cfg.LogLevel.String()).
		Msg("starting up")

	// set timeout for initialization
	ctxInit, closeCtx := context.WithTimeout(ctx, initTimeout)
	defer closeCtx()

	// set up meter
	shutdownMeter, err := prom.RegisterGlobal(ctxInit, serviceName, serviceVersion, metricNamespace)
	if err != nil {
		return err
	}
	defer doClose(shutdownMeter, logger)

	db, err := postgresql.NewDB(ctxInit, cfg.PG)
	if err != nil {
		return err
	}
	defer doClose(db.Close, logger)

	us := service.NewUserService(userRepo.NewUserRepo(db), logger)

	// set up http server
	httpSrv := httpServer.New(
		cfg.HTTP,
		[]httpServer.Routes{
			userHttp.RoutesFunc(us, logger),
			healthHandlers.RoutesFunc(logger),
		},
		logger)

	eg, ctx := errgroup.WithContext(ctx)
	// run service http server
	eg.Go(func() error {
		return httpSrv.Run(ctx)
	})
	// run prometheus metrics http server
	eg.Go(func() error {
		return prom.Serve(ctx, cfg.PromHTTP, logger)
	})

	return eg.Wait()
}

func doClose(fn func(context.Context) error, logger zerolog.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := fn(ctx); err != nil {
		logger.Error().
			Err(err).
			Msg("error closing")
	}
}

package config

import (
	"github.com/rs/zerolog"

	"github.com/Karzoug/gravitum-user-service/pkg/metric/prom"
	"github.com/Karzoug/gravitum-user-service/pkg/postgresql"

	httpConfig "github.com/Karzoug/gravitum-user-service/internal/delivery/http/config"
)

type Config struct {
	LogLevel zerolog.Level           `env:"LOG_LEVEL" envDefault:"info"`
	HTTP     httpConfig.ServerConfig `envPrefix:"HTTP_"`
	PromHTTP prom.ServerConfig       `envPrefix:"PROM_"`
	PG       postgresql.Config       `envPrefix:"PG_"`
}

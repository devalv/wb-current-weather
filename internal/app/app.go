package app

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/devalv/wb-current-weather/internal/config"
)

type Application struct {
	cfg *config.Config
}

func NewApplication(cfg *config.Config) *Application {
	app := &Application{cfg: cfg}
	return app
}

func (app *Application) Start(ctx context.Context) {
	log.Debug().Msg("Starting weather application")
	app.Stop(ctx)
}

func (app *Application) Stop(ctx context.Context) {
	log.Debug().Msg("Application stopped")
	os.Exit(0)
}

package app

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/devalv/wb-current-weather/internal/config"
	transport "github.com/devalv/wb-current-weather/internal/transport/http"
	"github.com/devalv/wb-current-weather/internal/usecase"
)

type Application struct {
	cfg *config.Config
}

func NewApplication(cfg *config.Config) *Application {
	app := &Application{cfg: cfg}
	return app
}

func (app *Application) getForecast(ctx context.Context) (fc usecase.Forecast, err error) {
	f, err := transport.GetForecast(ctx, app.cfg)
	if err != nil {
		return usecase.Forecast{}, err
	}
	return f, nil
}

func (app *Application) Start(ctx context.Context) {
	log.Debug().Msg("Starting weather application")
	fc, err := app.getForecast(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get forecast")
	}

	wo, err := usecase.NewWaybarOutput(fc.Text(), fc.TooltipInfo())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create Waybar output")
	}

	fmt.Println(wo)
	app.Stop(ctx)
}

func (app *Application) Stop(ctx context.Context) {
	log.Debug().Msg("Application stopped")
	os.Exit(0)
}

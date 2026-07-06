package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/devalv/wb-current-weather/internal/config"
	"github.com/devalv/wb-current-weather/internal/usecase"
)

const apiURL = "https://api.openweathermap.org/data/2.5/weather?id=%d&appid=%s&units=%s&lang=%s"

type weather struct {
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type main struct {
	Temp float32 `json:"temp"`
}

type wind struct {
	Speed float32 `json:"speed"`
}

type forecastResponse struct {
	Weather []weather `json:"weather"`
	Main    main      `json:"main"`
	Wind    wind      `json:"wind"`
}

func GetForecast(ctx context.Context, cfg *config.Config) (fr usecase.Forecast, err error) {
	requestURL := fmt.Sprintf(apiURL, cfg.CityID, cfg.WeatherAPIToken, cfg.Units, cfg.Lang)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, http.NoBody)
	if err != nil {
		return usecase.Forecast{}, fmt.Errorf("error creating http request: %w", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return usecase.Forecast{}, fmt.Errorf("making http request err: %w", err)
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("error closing request body")
		}
	}()

	log.Debug().Msgf("http response code: %v", res.StatusCode)
	if res.StatusCode != http.StatusOK {
		log.Error().Msgf("http response code: %v", res.StatusCode)

		return usecase.Forecast{}, fmt.Errorf("http response code: %v", res.StatusCode)
	}

	var f forecastResponse
	err = json.NewDecoder(res.Body).Decode(&f)
	if err != nil {
		return usecase.Forecast{}, fmt.Errorf("error decoding http response body: %w", err)
	}

	log.Debug().Msg(fmt.Sprintf("forecast: %v", f))

	return usecase.Forecast{
		Description: f.Weather[0].Description,
		Icon:        f.Weather[0].Icon,
		Temp:        f.Main.Temp,
		Wind:        f.Wind.Speed,
	}, nil
}

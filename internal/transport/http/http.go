package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/devalv/wb-current-weather/internal/config"
	"github.com/devalv/wb-current-weather/internal/usecase"
	"github.com/rs/zerolog/log"
)

const apiURL = "https://api.openweathermap.org/data/2.5/weather?id=%d&appid=%s&units=%s&lang=%s"

type weather struct {
	Descriotion string `json:"description"`
	Icon        string `json:"icon"`
}

type main struct {
	Temp float32 `json:"temp"`
}

type wind struct {
	Speed int `json:"speed"`
}

type forecastResponse struct {
	Weather []weather
	Main    main
	Wind    wind
}

func GetForecast(ctx context.Context, cfg *config.Config) (fr usecase.Forecast, err error) {
	requestURL := fmt.Sprintf(apiURL, cfg.CityID, cfg.WeatherAPIToken, cfg.Units, cfg.Lang)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, http.NoBody)
	if err != nil {
		log.Error().Err(err).Msg("error creating http request")
		return usecase.Forecast{}, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("error making http request")
		return usecase.Forecast{}, err
	}
	defer res.Body.Close()

	log.Debug().Msgf("http response code: %v", res.StatusCode)
	if res.StatusCode != http.StatusOK {
		log.Error().Msgf("http response code: %v", res.StatusCode)
		return usecase.Forecast{}, fmt.Errorf("http response code: %v", res.StatusCode)
	}

	var f forecastResponse
	err = json.NewDecoder(res.Body).Decode(&f)
	if err != nil {
		log.Error().Err(err).Msg("error decoding http response body")
		return usecase.Forecast{}, err
	}
	log.Debug().Msg(fmt.Sprintf("forecast: %v", f))
	return usecase.Forecast{Description: f.Weather[0].Descriotion, Icon: f.Weather[0].Icon, Temp: f.Main.Temp, Wind: f.Wind.Speed}, nil
}

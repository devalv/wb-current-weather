package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const apiURL = "https://api.openweathermap.org/data/2.5/weather?id=%d&appid=%s&units=%s&lang=%s"

type Config struct {
	Debug           bool   `yaml:"debug"`
	CityID          int    `yaml:"city_id"`
	WeatherAPIToken string `yaml:"weather_api_token"`
	Units           string `yaml:"units"`
	Lang            string `yaml:"lang"`
}

func NewConfig() (*Config, error) {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "./config.yml", "path to config file")
	flag.Parse()

	s, err := os.Stat(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if s.IsDir() {
		return nil, fmt.Errorf("'%s' is a directory, not a file", cfgPath)
	}

	log.Debug().Msgf("Config path: %s", cfgPath)

	cfg := &Config{
		Units:  "metric",
		Lang:   "ru",
		CityID: 524305, //nolint:mnd
	}

	err = cleanenv.ReadConfig(cfgPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config read err: %w", err)
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug mode is enabled")
	}

	log.Debug().Msgf("Config -> city_id: `%d`, units: `%s`, lang: `%s`", cfg.CityID, cfg.Units, cfg.Lang)

	return cfg, nil
}

type Forecast struct {
	Description string  `json:"description"`
	Icon        string  `json:"icon"`
	Temp        float32 `json:"temp"`
	Wind        float32 `json:"wind"`
}

func (f *Forecast) TooltipInfo() []string {
	return []string{
		fmt.Sprintf("🌡️: %.1f", f.Temp),
		fmt.Sprintf("🌬: %.1f", f.Wind),
		"📜: " + f.Description,
	}
}

func (f *Forecast) Text() string {
	iconRegistry := map[string]string{
		"01d": "☀️",
		"01n": "🌙",
		"02d": "🌤️",
		"02n": "☁️",
		"03d": "🌥️",
		"03n": "🌥️",
		"04d": "🌥️",
		"04n": "🌥️",
		"09d": "🌧️",
		"09n": "🌧️",
		"10d": "🌦️",
		"10n": "🌦️",
		"11d": "🌩️",
		"11n": "🌩️",
		"13d": "🌨️",
		"13n": "🌨️",
		"50d": "🌫️",
		"50n": "🌫️",
	}

	i := iconRegistry[f.Icon]
	if i == "" {
		i = "🌡️"
	}

	return fmt.Sprintf("%.0f | %s", f.Temp, i)
}

type WaybarOutput struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
}

func (wo *WaybarOutput) String() string {
	val, err := json.Marshal(wo)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to marshal waybar output")

		return ""
	}

	return string(val)
}

type weatherForecast struct {
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type mainForecast struct {
	Temp float32 `json:"temp"`
}

type windForecast struct {
	Speed float32 `json:"speed"`
}

type forecastResponse struct {
	Weather []weatherForecast `json:"weather"`
	Main    mainForecast      `json:"main"`
	Wind    windForecast      `json:"wind"`
}

func GetForecast(ctx context.Context, cfg *Config) (fc *Forecast, err error) {
	requestURL := fmt.Sprintf(apiURL, cfg.CityID, cfg.WeatherAPIToken, cfg.Units, cfg.Lang)

	// TODO: context with timeout
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %w", err)
	}

	// TODO: client with timeout
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making http request err: %w", err)
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("error closing request body")
		}
	}()

	if res.StatusCode != http.StatusOK {
		log.Error().Msgf("http response code: %v", res.StatusCode)

		return nil, fmt.Errorf("http response code: %v", res.StatusCode)
	}

	var f forecastResponse
	err = json.NewDecoder(res.Body).Decode(&f)
	if err != nil {
		return nil, fmt.Errorf("error decoding http response body: %w", err)
	}

	log.Debug().Msg(fmt.Sprintf("server forecast: %v", f))

	return &Forecast{
		Description: f.Weather[0].Description,
		Icon:        f.Weather[0].Icon,
		Temp:        f.Main.Temp,
		Wind:        f.Wind.Speed,
	}, nil
}

func run() error {
	cfg, err := NewConfig()
	if err != nil {
		return fmt.Errorf("failed to read config %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, syscall.SIGSEGV)
	defer cancel()

	fc, err := GetForecast(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to get forecast: %w", err)
	}
	if fc == nil {
		return errors.New("received nil forecast without error")
	}

	wo := &WaybarOutput{
		Text:    fc.Text(),
		Tooltip: strings.Join(fc.TooltipInfo(), "\n"),
	}

	fmt.Println(wo)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("application failed")
	}
}

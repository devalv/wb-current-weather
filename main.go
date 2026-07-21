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
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
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

	cfg := &Config{
		Units:  "metric",
		Lang:   "ru",
		CityID: 524305, //nolint:mnd
	}

	err = cleanenv.ReadConfig(cfgPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config read err: %w", err)
	}

	return cfg, nil
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

type forecastResponse struct {
	Weather []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Main struct {
		Temp float32 `json:"temp"`
	} `json:"main"`
	Wind struct {
		Speed float32 `json:"speed"`
	} `json:"wind"`
}

func (f *forecastResponse) TooltipInfo() []string {
	return []string{
		fmt.Sprintf("🌡️: %.1f", f.Main.Temp),
		fmt.Sprintf("🌬: %.1f", f.Wind.Speed),
		"📜: " + f.Weather[0].Description,
	}
}

func (f *forecastResponse) Text() string {
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

	i := iconRegistry[f.Weather[0].Icon]
	if i == "" {
		i = "🌡️"
	}

	return fmt.Sprintf("%.0f | %s", f.Main.Temp, i)
}

func GetForecast(ctx context.Context, cfg *Config) (fR *forecastResponse, err error) {
	const timeout = 15 * time.Second
	const apiURL = "https://api.openweathermap.org/data/2.5/weather?id=%d&appid=%s&units=%s&lang=%s"

	requestURL := fmt.Sprintf(apiURL, cfg.CityID, cfg.WeatherAPIToken, cfg.Units, cfg.Lang)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %w", err)
	}

	apiClient := &http.Client{
		Timeout: timeout,
	}

	res, err := apiClient.Do(req)
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
	if err := json.NewDecoder(res.Body).Decode(&f); err != nil {
		return nil, fmt.Errorf("error decoding http response body: %w", err)
	}

	return &f, nil
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

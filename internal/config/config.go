package config

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Debug           bool   `yaml:"debug"`
	CityID          int    `yaml:"city_id"`
	WeatherAPIToken string `yaml:"weather_api_token"`
	Units           string `yaml:"units"`
	Lang            string `yaml:"lang"`
	ConfigPath      string
}

func validateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("failed to read config file: %w", err)
		}

		return fmt.Errorf("failed to get config path stat: %w", err)
	}

	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a file", path)
	}

	return nil
}

func parseFlags() (path string, err error) {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "./config.yml", "path to config file")
	flag.Parse()

	if err := validateConfigPath(cfgPath); err != nil {
		return "", err
	}

	return cfgPath, nil
}

func NewConfig() (*Config, error) {
	cfgPath, err := parseFlags()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse flags")
	}

	var cfg Config
	err = cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("config read err: %w", err)
	}

	if cfg.Units == "" {
		cfg.Units = "metric"
	}

	if cfg.Lang == "" {
		cfg.Lang = "ru"
	}

	cfg.ConfigPath = cfgPath
	cfg.configureLogger()

	return &cfg, nil
}

func (cfg *Config) configureLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug mode enabled")
	}
}

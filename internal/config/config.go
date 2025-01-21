package config

import (
	"flag"
	"fmt"
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
		return err
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
		return nil, err
	}

	// TODO: default units
	// TODO: default lang
	// TODO: default city_id

	cfg.ConfigPath = cfgPath
	cfg.ConfigureLogger()
	return &cfg, nil
}

func (cfg *Config) ConfigureLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug mode enabled")
	}
}

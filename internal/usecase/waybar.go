package usecase

import (
	"encoding/json"
	"strings"

	"github.com/rs/zerolog/log"
)

type WaybarOutput struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
}

func (wo WaybarOutput) String() string {
	val, err := json.Marshal(wo)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal waybar output")
		return ""
	}
	return string(val)
}

func NewWaybarOutput(text string, tooltipInfo []string) (WaybarOutput, error) {
	return WaybarOutput{
		Text:    text,
		Tooltip: strings.Join(tooltipInfo, "\n"),
	}, nil
}

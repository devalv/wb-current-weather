package usecase

import "fmt"

type Forecast struct {
	Description string  `json:"description"`
	Icon        string  `json:"icon"`
	Temp        float32 `json:"temp"`
	Wind        float32 `json:"wind"`
}

func (f *Forecast) TooltipInfo() []string {
	return []string{
		fmt.Sprintf("Температура: %.1f", f.Temp),
		fmt.Sprintf("Скорость ветра: %.1f", f.Wind),
		"Описание: " + f.Description,
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

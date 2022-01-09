package riley

import "github.com/drTragger/rileyBot/storage"

type Config struct {
	BotToken    string `toml:"bot_token"`
	WeatherKey  string `toml:"weather_key"`
	LoggerLevel string `toml:"logger_level"`
	Storage     *storage.Config
}

func NewConfig() *Config {
	return &Config{
		LoggerLevel: "debug",
		Storage:     storage.NewConfig(),
	}
}

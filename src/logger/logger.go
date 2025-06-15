package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yamaki-87/mailbot/src/config"
)

func getConfig() *config.AppConfig {
	return config.GetConfig()
}

func loadLogLevel() string {
	config := getConfig()

	level := config.Log.Level
	if level == "" {
		level = "info"
	}
	return level
}

func Init() {
	levelStr := loadLogLevel()

	level, err := zerolog.ParseLevel(strings.ToLower(levelStr))
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	log.Info().Str("level", levelStr).Msg("Logger initalized")
}

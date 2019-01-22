package internal

import (
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

// Logger represents an active logging object
var Logger zerolog.Logger

func init() {
	nocolor, _ := strconv.ParseBool(GetEnv("LOG_NOCOLOR", "false"))
	Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC1123,
		NoColor:    nocolor,
	}).With().Timestamp().Logger()

	logLevel, err := zerolog.ParseLevel(GetEnv("LOG_LEVEL", "info"))
	if err != nil {
		Logger.Error().Err(err).Msgf("Unknown log level")
	}
	zerolog.SetGlobalLevel(logLevel)
}

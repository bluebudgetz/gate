package infra

import (
	glog "log"
	"os"
	"strings"

	"github.com/huandu/go-tls"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupLogging() {
	zerolog.ErrorStackMarshaler = stackTraceMarshaller

	switch os.Getenv("LOG_LEVEL") {
	case "disabled":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info", "":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		log.Fatal().Msg("log level must be one of: disabled, panic, fatal, error, warn, info or debug")
	}

	if boolEnv("LOG_PRETTY") {
		log.Logger = log.Output(newPrettyConsoleWriter())
	}
	if boolEnv("LOG_PID") {
		log.Logger = log.With().
			Int("pid", os.Getpid()).
			Logger().Hook(zerolog.HookFunc(threadIDHook))
	}
	if boolEnv("LOG_SVC") {
		log.Logger = log.With().
			Str("svc", "gate").
			Logger()
	}
	log.Logger = log.Logger.With().
		CallerWithSkipFrameCount(2).
		Stack().
		Logger()

	glog.SetFlags(0)
	glog.SetOutput(log.Logger)
}

func threadIDHook(e *zerolog.Event, _ zerolog.Level, _ string) {
	// Used by request logger
	e.Int64("tid", tls.ID())
}

func boolEnv(key string) bool {
	value := strings.ToLower(os.Getenv(key))
	return value == "1" || value == "y" || value == "yes" || value == "true"
}

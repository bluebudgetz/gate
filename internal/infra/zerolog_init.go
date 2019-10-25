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
	logPrettyEnv := strings.ToLower(os.Getenv("LOG_PRETTY"))
	logPretty := logPrettyEnv == "1" || logPrettyEnv == "y" || logPrettyEnv == "yes" || logPrettyEnv == "true"
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

	log.Logger = log.Logger.Hook(zerolog.HookFunc(threadIDHook))
	if logPretty {
		log.Logger = log.Output(newConsoleWriter()).With().Int("pid", os.Getpid()).Logger()
	} else {
		log.Logger = log.Logger.With().Str("svc", "gate").Logger()
	}
	log.Logger = log.Logger.With().CallerWithSkipFrameCount(2).Stack().Logger()

	glog.SetFlags(0)
	glog.SetOutput(log.Logger)
}

func threadIDHook(e *zerolog.Event, _ zerolog.Level, _ string) {
	// Used by request logger
	e.Int64("tid", tls.ID())
}

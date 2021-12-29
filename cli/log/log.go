package log

import (
	"io"
	"os"
	"reflect"
	"time"

	"github.com/rs/zerolog"
)

const (
	DEBUG_LEVEL = zerolog.DebugLevel
	INFO_LEVEL  = zerolog.InfoLevel
	WARN_LEVEL  = zerolog.WarnLevel
	ERR_LEVEL   = zerolog.ErrorLevel
	FATAL_LEVEL = zerolog.FatalLevel

	KeyModule = "mod"
	KeyEvent  = "event"

	ModuleConfirmer = "confirmer"
	ModuleMain      = "main"
)

// global
var Logger = zerolog.New(writer).With().Timestamp().Logger()
var ConsoleWriter = zerolog.ConsoleWriter{Out: os.Stdout}

var (
	writer            = &Writer{Out: os.Stderr}
	defaultLevel      = DEBUG_LEVEL
	defaultWriter     = ConsoleWriter
	defaultTimeFormat = time.RFC3339
)

// default config
func init() {
	SetLevel(defaultLevel)
	SetWriter(defaultWriter)
	SetTimeFormat(defaultTimeFormat)
}

func SetLevel(lv zerolog.Level) {
	Logger = Logger.Level(lv)
}

func SetWriter(w io.Writer) {
	writer.SetWriter(w)
}

func SetTimeFormat(format string) {
	zerolog.TimeFieldFormat = format
	if reflect.TypeOf(writer.Out).Name() == "ConsoleWriter" {
		// reset writer
		ConsoleWriter = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: format,
		}
		SetWriter(ConsoleWriter)
	}
}

func Confirmer(event string) zerolog.Logger {
	if event == "" {
		return Logger.With().Str(KeyModule, ModuleConfirmer).Logger()
	}
	return Logger.With().Str(KeyModule, ModuleConfirmer).Str(KeyEvent, event).Logger()
}

func Main(event string) zerolog.Logger {
	if event == "" {
		return Logger.With().Str(KeyModule, ModuleMain).Logger()
	}
	return Logger.With().Str(KeyModule, ModuleMain).Str(KeyEvent, event).Logger()
}

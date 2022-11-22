package zerolog

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type NoMsgHook struct{}

func (h NoMsgHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if msg == "" {
		e.Str("message", "no message")
	}
}

// Config setting log config
type Config struct {
	Debug  bool `yaml:"debug" mapstructure:"debug"`
	Pretty bool `yaml:"pretty" mapstructure:"pretty"`
}

// Init logger singleton
func Init(cfg Config) {
	var outWriter io.Writer = os.Stdout
	if cfg.Pretty {
		zerolog.ErrorStackMarshaler = ESPackMarshalStack
		outWriter = newConsoleWriter()
	} else {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	}
	logger := zerolog.New(outWriter).
		With().
		Timestamp().
		Caller().
		Stack().
		Logger().
		Hook(NoMsgHook{})
	log.Logger = logger
}

func newConsoleWriter() zerolog.ConsoleWriter {
	console := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006/01/02 15:04:05",
	}
	console.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("[ %s ]", i)
	}
	return console
}

package fxlogger

import (
	"strings"

	"github.com/rs/zerolog/log"
	"go.uber.org/fx/fxevent"
)

type fxLogger struct{}

func NewLogger() fxevent.Logger {
	return &fxLogger{}
}

// LogEvent logs the given event to the provided Zap logger.
func (l *fxLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		log.Info().
			Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStart hook executing")
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			log.Info().
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Msg("OnStart hook failed")
		} else {
			log.Info().
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msg("OnStart hook executed")
		}
	case *fxevent.OnStopExecuting:
		log.Info().
			Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStop hook executing")
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			log.Info().
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Err(e.Err).
				Msg("")
		} else {
			log.Info().
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msg("OnStop hook executed")
		}
	case *fxevent.Supplied:
		log.Info().
			Str("type", e.TypeName).
			Str("module_field", e.ModuleName).
			Err(e.Err).
			Msg("supplied")
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			log.Info().
				Str("constructor", e.ConstructorName).
				Str("module_field", e.ModuleName).
				Str("type", rtype).
				Msgf("provided %s", rtype)
		}
		if e.Err != nil {
			log.Error().
				Str("module_field", e.ModuleName).
				Err(e.Err).
				Msg("error encountered while applying options")
		}
	case *fxevent.Replaced:
		for _, rtype := range e.OutputTypeNames {
			log.Info().
				Str("module_field", e.ModuleName).
				Str("type", rtype).
				Msg("replaced")
		}
		if e.Err != nil {
			log.Error().
				Str("module_field", e.ModuleName).
				Err(e.Err).
				Msg("error encountered while replacing")
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			log.Info().
				Str("decorator", e.DecoratorName).
				Str("module_field", e.ModuleName).
				Str("type", rtype).
				Msg("decorated")
		}
		if e.Err != nil {
			log.Error().
				Str("module_field", e.ModuleName).
				Err(e.Err).
				Msg("error encountered while applying options")
		}
	case *fxevent.Invoking:
		// Do not log stack as it will make logs hard to read.
		log.Info().
			Str("function", e.FunctionName).
			Str("module_field", e.ModuleName).
			Msgf("invoking %s", e.FunctionName)
	case *fxevent.Invoked:
		if e.Err != nil {
			log.Error().
				Str("stack", e.Trace).
				Str("function", e.FunctionName).
				Str("module_field", e.ModuleName).
				Err(e.Err).
				Msg("invoke failed")
		}
	case *fxevent.Stopping:
		log.Info().
			Str("signal", strings.ToUpper(e.Signal.String())).
			Msg("received signal")
	case *fxevent.Stopped:
		if e.Err != nil {
			log.Error().Err(e.Err).Msg("stop failed")
		}
	case *fxevent.RollingBack:
		log.Error().Err(e.StartErr).Msg("start failed, rolling back")
	case *fxevent.RolledBack:
		if e.Err != nil {
			log.Error().Err(e.Err).Msg("rollback failed")
		}
	case *fxevent.Started:
		if e.Err != nil {
			log.Error().Err(e.Err).Msg("start failed")
		} else {
			log.Info().Msg("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			log.Error().Err(e.Err).Msg("custom logger initialization failed")
		} else {
			log.Info().Str("function", e.ConstructorName).Msg("initialized custom fxevent.Logger")
		}
	}
}

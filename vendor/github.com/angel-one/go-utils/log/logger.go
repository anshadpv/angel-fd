package log

import (
	"context"
	"io"
	"runtime/debug"

	"github.com/angel-one/go-utils/constants"
	"github.com/angel-one/go-utils/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	LogTypeKey = "log_type"

	// Pre-defined log types as per https://angelbrokingpl.atlassian.net/wiki/spaces/Technology/pages/2809299426/Log+Management+Program#Log-Format
	// Will be defaulted to Application Log Type if no log type is specified.
	LogTypeAccess                 = "access"
	LogTypeAudit                  = "audit"
	LogTypeTransaction            = "transaction"
	LogTypeApplication            = "application"
	LogTypeApplicationPerformance = "application-performance"
	LogTypeProcessAudit           = "process-audit"
)

type Level string

func (l Level) zeroLogLevel() zerolog.Level {
	switch l {
	case constants.TraceLevel:
		return zerolog.TraceLevel
	case constants.DebugLevel:
		return zerolog.DebugLevel
	case constants.InfoLevel:
		return zerolog.InfoLevel
	case constants.WarnLevel:
		return zerolog.WarnLevel
	case constants.ErrorLevel:
		return zerolog.ErrorLevel
	case constants.FatalLevel:
		return zerolog.FatalLevel
	case constants.PanicLevel:
		return zerolog.PanicLevel
	default:
		return zerolog.DebugLevel
	}
}

// InitLogger is used to initialize logger
func InitLogger(level Level) {
	zerolog.ErrorStackMarshaler = getErrorStackMarshaller()
	zerolog.SetGlobalLevel(level.zeroLogLevel())
	log.Logger = log.With().Caller().Logger()
}

// InitLoggerWithWriter is used to initialize logger with a writer
func InitLoggerWithWriter(level Level, w io.Writer) {
	zerolog.ErrorStackMarshaler = getErrorStackMarshaller()
	zerolog.SetGlobalLevel(level.zeroLogLevel())
	log.Logger = zerolog.New(w).With().Caller().Timestamp().Logger()
}

// Trace is the for trace log
func Trace(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Trace(), LogTypeApplication)
}

// Debug is the for debug log
func Debug(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Debug(), LogTypeApplication)
}

// Info is the for info log
func Info(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Info(), LogTypeApplication)
}

// Warn is the for warn log
func Warn(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Warn(), LogTypeApplication)
}

// Error is the for error log
func Error(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Error().Stack(), LogTypeApplication)
}

// Panic is the for panic log
func Panic(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Panic().Stack(), LogTypeApplication)
}

// Fatal is the for fatal log
func Fatal(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Fatal().Stack(), LogTypeApplication)
}

// ErrorWarn checks for the error object.
// In case it is corresponding to a 4XX status code, it logs it as warning.
// Otherwise, it logs it as an error.
func ErrorWarn(ctx context.Context, err error) *zerolog.Event {
	if e, ok := err.(*errors.Error); ok && e.StatusCode >= 400 && e.StatusCode < 500 {
		return Warn(ctx).Err(err)
	}
	return Error(ctx).Err(err)
}

func getErrorStackMarshaller() func(err error) interface{} {
	return func(err error) interface{} {
		if err != nil {
			if e, ok := err.(*errors.Error); ok {
				return map[string]interface{}{
					constants.CodeLogParam:    e.Code,
					constants.MessageLogParam: e.Message,
					constants.DetailsLogParam: e.Details,
					constants.TraceLogParam:   e.GetTrace(),
				}
			}
		}
		return string(debug.Stack())
	}
}

func withIDAndPath(ctx context.Context, event *zerolog.Event, logType string) *zerolog.Event {
	if ctx == nil {
		return event
	}
	id := ctx.Value(constants.IDLogParam)
	if id != nil {
		event.Interface(constants.IDLogParam, id)
	}
	path := ctx.Value(constants.PathLogParam)
	if path != nil {
		event.Interface(constants.PathLogParam, path)
	}
	correlationId := ctx.Value(constants.CorrelationLogParam)
	if correlationId != nil {
		event.Interface(constants.CorrelationLogParam, correlationId)
	}
	clientID := ctx.Value(constants.ClientIDLogParam)
	if clientID != nil {
		event.Interface(constants.ClientIDLogParam, clientID)
	}
	action := ctx.Value(constants.ActionLogParam)
	if action != nil {
		event.Interface(constants.ActionLogParam, action)
	}
	event.Str(LogTypeKey, logType)
	return event
}

type logTypeEvent struct {
	logType string
}

func Access() *logTypeEvent {
	return &logTypeEvent{logType: LogTypeAccess}
}

func Audit() *logTypeEvent {
	return &logTypeEvent{logType: LogTypeAudit}
}

func Transaction() *logTypeEvent {
	return &logTypeEvent{logType: LogTypeTransaction}
}

func Application() *logTypeEvent {
	return &logTypeEvent{logType: LogTypeApplication}
}

func ApplicationPerformance() *logTypeEvent {
	return &logTypeEvent{logType: LogTypeApplicationPerformance}
}

func ProcessAudit() *logTypeEvent {
	return &logTypeEvent{logType: LogTypeProcessAudit}
}

// Trace is the for trace log
func (ev *logTypeEvent) Trace(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Trace(), ev.logType)
}

// Debug is the for debug log
func (ev *logTypeEvent) Debug(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Debug(), ev.logType)
}

// Info is the for info log
func (ev *logTypeEvent) Info(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Info(), ev.logType)
}

// Warn is the for warn log
func (ev *logTypeEvent) Warn(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Warn(), ev.logType)
}

// Error is the for error log
func (ev *logTypeEvent) Error(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Error(), ev.logType)
}

// Panic is the for panic log
func (ev *logTypeEvent) Panic(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Panic(), ev.logType)
}

// Fatal is the for fatal log
func (ev *logTypeEvent) Fatal(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Fatal(), ev.logType)
}

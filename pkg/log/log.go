package log

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/samber/lo"
)

type Context zerolog.Context

var Logger *zerolog.Logger

func New() *Context {
	return lo.ToPtr((Context)(Logger.With()))
}

func (e *Context) WithField(key string, value any) *Context {
	return lo.ToPtr((Context)((*zerolog.Context)(e).Str(key, logFields(value))))
}

func (e *Context) WithFields(fields map[string]any) *Context {
	for key, value := range fields {
		fields[key] = logFields(value)
	}

	return lo.ToPtr((Context)((*zerolog.Context)(e).Fields(fields)))
}

func (e *Context) WithError(err error) *Context {
	return lo.ToPtr((Context)((*zerolog.Context)(e).Err(err).Stack()))
}

func (e *Context) WithErrorf(format string, args ...interface{}) *Context {
	return lo.ToPtr((Context)((*zerolog.Context)(e).Err(fmt.Errorf(format, args...)).Stack()))
}

func (e *Context) WithRequest(request interface{}) *Context {
	return e.WithField("request", request)
}

func (e *Context) WithResponse(response interface{}) *Context {
	return e.WithField("response", response)
}

func (e *Context) WithUserID(userID int) *Context {
	return e.WithField("user_id", userID)
}

func (e *Context) WithRequestID(requestID string) *Context {
	return e.WithField("request_id", requestID)
}

func (e *Context) Print(args ...interface{}) {
	logger := (*zerolog.Context)(e).Logger()
	logger.Print(fmt.Sprint(args...))
}

func (e *Context) Debug(args ...interface{}) {
	logger := (*zerolog.Context)(e).Logger()
	logger.Debug().Msg(fmt.Sprint(args...))
}

func (e *Context) Info(args ...interface{}) {
	logger := (*zerolog.Context)(e).Logger()
	logger.Info().Msg(fmt.Sprint(args...))
}

func (e *Context) Infof(format string, args ...interface{}) {
	logger := (*zerolog.Context)(e).Logger()
	logger.Info().Msg(fmt.Sprintf(format, args...))
}

func (e *Context) Warn(args ...interface{}) {
	logger := (*zerolog.Context)(e).Logger()
	logger.Warn().Msg(fmt.Sprint(args...))
}

func (e *Context) Warnf(format string, args ...interface{}) {
	logger := (*zerolog.Context)(e).Logger()
	logger.Warn().Msg(fmt.Sprintf(format, args...))
}

func (e *Context) Error(args ...interface{}) {
	logger := (*zerolog.Context)(e).Logger()
	logger.Error().Msg(fmt.Sprint(args...))
}

func (e *Context) Errorf(format string, args ...interface{}) {
	logger := (*zerolog.Context)(e).Logger()
	logger.Error().Msg(fmt.Sprintf(format, args...))
}

func (e *Context) Fatal(args ...interface{}) {
	logger := (*zerolog.Context)(e).Logger()
	logger.Fatal().Msg(fmt.Sprint(args...))
}

func (e *Context) Fatalf(format string, args ...interface{}) {
	logger := (*zerolog.Context)(e).Logger()
	logger.Fatal().Msg(fmt.Sprintf(format, args...))
}

func (e *Context) Panic(args ...interface{}) {
	logger := (*zerolog.Context)(e).Logger()
	logger.Panic().Msg(fmt.Sprint(args...))
}

func (e *Context) Panicf(format string, args ...interface{}) {
	logger := (*zerolog.Context)(e).Logger()
	logger.Panic().Msg(fmt.Sprintf(format, args...))
}

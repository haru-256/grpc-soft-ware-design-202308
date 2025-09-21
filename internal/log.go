package internal

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger(logger *slog.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

func (l *Logger) logStart(ctx context.Context, req connect.AnyRequest) {
	l.logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		"start processing request",
slog.String("procedure", req.Spec().Procedure),
		slog.String("user-agent", req.Header().Get("User-Agent")),
		slog.String("http_method", req.HTTPMethod()),
		slog.String("peer", req.Peer().Addr),
	)
}

// logEnd はリクエスト/ストリームの終了をログに出力します。
func (l *Logger) logEnd(ctx context.Context, req connect.AnyRequest, startTime time.Time, err error) {
	duration := time.Since(startTime)
	logFields := []slog.Attr{
		slog.String("procedure", req.Spec().Procedure),
		slog.String("code", connect.CodeOf(err).String()),
		slog.Duration("duration", duration),
	}

	if err != nil {
		logFields = append(logFields, slog.String("error", err.Error()))
		l.logger.LogAttrs(ctx, slog.LevelError, "finished with error", logFields...)
	} else {
		l.logger.LogAttrs(ctx, slog.LevelInfo, "finished", logFields...)
	}
}

// file: interceptor/logger.go (続き)

// NewUnaryInterceptor は共通ロガーを使ったUnaryインターセプターを返します。
func (l *Logger) NewUnaryInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			startTime := time.Now()

			l.logStart(ctx, req)
			res, err := next(ctx, req)
			l.logEnd(ctx, req, startTime, err)

			return res, err
		}
	}
}

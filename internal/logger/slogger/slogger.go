package slogger

import (
	"context"
	"log/slog"
	"os"
)

func InitLogging() {
	handler := slog.Handler(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	handler = NewHandlerMiddleware(handler)
	slog.SetDefault(slog.New(handler))
}

type HandlerMiddleware struct {
	next slog.Handler
}

func NewHandlerMiddleware(next slog.Handler) *HandlerMiddleware {
	return &HandlerMiddleware{next: next}
}

func (h *HandlerMiddleware) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

func (h *HandlerMiddleware) Handle(ctx context.Context, rec slog.Record) error {
	return h.next.Handle(ctx, rec)
}

func (h *HandlerMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddleware{next: h.next.WithAttrs(attrs)}
}

// WithGroup implements Handler.WithGroup.
func (h *HandlerMiddleware) WithGroup(name string) slog.Handler {
	return &HandlerMiddleware{next: h.next.WithGroup(name)}
}

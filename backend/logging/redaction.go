package logging

import (
	"context"
	"log/slog"

	"github.com/samber/lo"
	slogmulti "github.com/samber/slog-multi"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/ficsitcli"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
)

type redactGamePathCredentialsMiddleware struct {
	next slog.Handler
}

func newRedactGamePathCredentialsMiddleware() slogmulti.Middleware {
	return func(next slog.Handler) slog.Handler {
		return &redactGamePathCredentialsMiddleware{
			next: next,
		}
	}
}

func (r redactGamePathCredentialsMiddleware) Enabled(ctx context.Context, level slog.Level) bool {
	return r.next.Enabled(ctx, level)
}

func (r redactGamePathCredentialsMiddleware) Handle(ctx context.Context, record slog.Record) error {
	attrs := make([]slog.Attr, 0, record.NumAttrs())

	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, redactPaths(attr))
		return true
	})

	// new record with redacted paths
	record = slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
	record.AddAttrs(attrs...)

	return r.next.Handle(ctx, record) //nolint:wrapcheck
}

func (r redactGamePathCredentialsMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	for i := range attrs {
		attrs[i] = redactPaths(attrs[i])
	}
	return &redactGamePathCredentialsMiddleware{
		next: r.next.WithAttrs(attrs),
	}
}

func (r redactGamePathCredentialsMiddleware) WithGroup(name string) slog.Handler {
	return &redactGamePathCredentialsMiddleware{
		next: r.next.WithGroup(name),
	}
}

func redactPaths(attr slog.Attr) slog.Attr {
	k := attr.Key
	v := attr.Value
	kind := attr.Value.Kind()

	switch kind {
	case slog.KindGroup:
		attrs := v.Group()
		for i := range attrs {
			attrs[i] = redactPaths(attrs[i])
		}
		return slog.Group(k, lo.ToAnySlice(attrs)...)
	case slog.KindString:
		if isGamePath(v.String()) {
			return slog.String(k, utils.RedactPath(v.String()))
		}
	default:
		break
	}
	return attr
}

func isGamePath(str string) bool {
	if ficsitcli.FicsitCLI != nil {
		return ficsitcli.FicsitCLI.GetInstallation(str) != nil
	}
	// if ficsitcli is not initialized, we can't know if it's a game path
	// so any code running before that should not log game paths
	return false
}

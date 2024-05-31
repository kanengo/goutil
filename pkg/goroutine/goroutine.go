package goroutine

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"maps"
	"sync"
)

type inheritDataMapCtxType struct {
}

func GoSafe(fn func(), panicHandler func(r any)) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				if panicHandler != nil {
					panicHandler(r)
				}
			}
		}()
		fn()
	}()
}

func inheritDataMap(dst context.Context, src context.Context) context.Context {
	inheritDataMap := src.Value(inheritDataMapCtxType{})
	if inheritDataMap != nil {
		m := inheritDataMap.(map[string]any)
		if len(m) == 0 {
			return dst
		}
		newMap := make(map[string]any, len(m))
		maps.Copy(newMap, m)
		dst = context.WithValue(dst, inheritDataMapCtxType{}, newMap)
	}

	return dst
}

var initTracer sync.Once
var tracer trace.Tracer

const traceName = "github.com/kanengo/tracer/goroutine"

func getTracer() trace.Tracer {
	initTracer.Do(func() {
		tracer = otel.GetTracerProvider().Tracer(traceName)
	})
	return tracer
}

func GoSafeWithTraceContext(ctx context.Context, spanName string, fn func(ctx context.Context) error, panicHandler func(ctx context.Context, r any)) {
	goCtx := context.Background()
	sc := trace.SpanContextFromContext(ctx)
	needTrace := false
	if sc.IsValid() {
		needTrace = true
		goCtx = trace.ContextWithSpanContext(ctx, sc)
	}

	goCtx = inheritDataMap(goCtx, ctx)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				if panicHandler != nil {
					panicHandler(goCtx, err)
				}
			}
		}()
		var err error
		if needTrace {
			var span trace.Span
			goCtx, span = getTracer().Start(goCtx, spanName, trace.WithSpanKind(trace.SpanKindInternal))
			defer func() {
				if err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
				}
				span.End()
			}()
		}
		err = fn(goCtx)
	}()
}

func GoSafeWithContext(ctx context.Context, fn func(ctx context.Context), panicHandler func(ctx context.Context, r any)) {
	goCtx := context.Background()
	sc := trace.SpanContextFromContext(ctx)
	if sc.IsValid() {
		goCtx = trace.ContextWithSpanContext(ctx, sc)
	}
	goCtx = inheritDataMap(goCtx, ctx)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if panicHandler != nil {
					panicHandler(ctx, err)
				}
			}
		}()
		fn(goCtx)
	}()
}

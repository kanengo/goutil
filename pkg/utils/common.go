package utils

import (
	"context"

	"github.com/kanengo/goutil/pkg/log"
	"go.uber.org/zap"
)

var panicRecoverHandlerFn func(ctx context.Context)

func SetCommonPanicRecoverHandler(fn func(ctx context.Context)) {
	panicRecoverHandlerFn = fn
}

func CheckGoPanic(ctx context.Context, callback func(context.Context)) {
	if r := recover(); r != nil {
		log.DPanic("panic recoverd", zap.Any("msg", r))
		if callback != nil {
			callback(ctx)
		} else if panicRecoverHandlerFn != nil {
			panicRecoverHandlerFn(ctx)
		}
	}
}

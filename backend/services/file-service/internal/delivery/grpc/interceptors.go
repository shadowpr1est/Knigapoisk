package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			st, _ := status.FromError(err)
			logger.Error("grpc request",
				zap.String("method", info.FullMethod),
				zap.Duration("duration", duration),
				zap.String("code", st.Code().String()),
				zap.String("error", st.Message()),
			)
		} else {
			logger.Info("grpc request",
				zap.String("method", info.FullMethod),
				zap.Duration("duration", duration),
			)
		}

		return resp, err
	}
}

func RecoveryInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered",
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
				)
				err = status.Error(13, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}


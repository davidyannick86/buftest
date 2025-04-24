package interceptor

import (
	"context"
	"fmt"
	"log"

	valid "github.com/bufbuild/protovalidate-go"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// initialize validator once at startup
var validator = func() valid.Validator {
	v, err := valid.New()
	if err != nil {
		log.Fatalf("validator init failed: %v", err)
	}
	return v
}()

func UnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if pm, ok := req.(proto.Message); ok {
		if err := validator.Validate(pm); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
	}
	return handler(ctx, req)
}

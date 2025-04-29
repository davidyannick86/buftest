package interceptor

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"

	valid "github.com/bufbuild/protovalidate-go"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// initialize validator once at startup.
var validator = func() valid.Validator {
	v, err := valid.New()
	if err != nil {
		log.Error().Err(err).Msg("failed to create validator")
	}
	return v
}()

// UnaryServerInterceptor est un intercepteur gRPC unary qui retourne une réponse générique.
//
//nolint:ireturn // required by gRPC handler signature
func UnaryServerInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if pm, ok := req.(proto.Message); ok {
		if err := validator.Validate(pm); err != nil {
			var ve *valid.ValidationError
			if errors.As(err, &ve) {
				log.Info().Msgf("Validation violations: %v", ve.Error())
			}
			return nil, fmt.Errorf("validation failed: %w", err)
		}
	}

	log.Info().Msgf("Received request: %v", req)
	return handler(ctx, req)
}

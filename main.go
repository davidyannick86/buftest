package main

import (
	"context"
	"fmt"

	"net"
	"net/http"
	"os"

	valid "github.com/bufbuild/protovalidate-go"
	proto_hello "github.com/davidyannick86/bufbuild/testbuf/protogen/hello/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	proto_hello.UnimplementedHelloServiceServer
}

func (s *server) SayHello(ctx context.Context, req *proto_hello.SayHelloRequest) (*proto_hello.SayHelloResponse, error) {
	validator, err := valid.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create validator: %v", err)
	}

	log.Info().Msgf("Received request: %v", req)

	if err := validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	response := &proto_hello.SayHelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
	}
	return response, nil
}

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	go func() {
		grpcMux := runtime.NewServeMux()
		err := proto_hello.RegisterHelloServiceHandlerServer(context.Background(), grpcMux, &server{})
		if err != nil {
			log.Fatal().Err(err).Msg("failed to register handler")
		}
		mux := http.NewServeMux()
		mux.Handle("/", grpcMux)

		listener, err := net.Listen("tcp", ":8080")
		if err != nil {
			log.Fatal().Err(err).Msg("failed to listen")
		}

		log.Info().Msgf("Gateway server listening at %v", listener.Addr())

		if err := http.Serve(listener, mux); err != nil {
			log.Fatal().Err(err).Msg("failed to serve")
		}

	}()

	grpcServer := grpc.NewServer()
	proto_hello.RegisterHelloServiceServer(grpcServer, &server{})

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	log.Info().Msgf("GRPC server listening at %v", listener.Addr())

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}

}

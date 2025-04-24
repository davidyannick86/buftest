package main

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"syscall"

	"net"
	"net/http"
	"os"

	valid "github.com/bufbuild/protovalidate-go"
	proto_hello "github.com/davidyannick86/bufbuild/testbuf/protogen/hello/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

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

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		interruptSignals...,
	)

	defer cancel()

	log.Info().Msg("Starting server...")

	wg, ctx := errgroup.WithContext(ctx)

	runGrpcServer(ctx, wg)
	runHTTOGateway(ctx, wg)

	err := wg.Wait()
	if err != nil {
		log.Error().Err(err).Msg("Error in server")
	} else {
		log.Info().Msg("Server stopped gracefully")
	}

}

func runHTTOGateway(
	ctx context.Context,
	wg *errgroup.Group,
) {
	grpcMux := runtime.NewServeMux()
	err := proto_hello.RegisterHelloServiceHandlerServer(context.Background(), grpcMux, &server{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register handler")
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	httpServer := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}

	wg.Go(func() error {
		log.Info().Msgf("HTTP Gateway server listening at %v", httpServer.Addr)
		err := httpServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Error().Err(err).Msg("failed to serve")
			return err
		}
		return nil
	})

	wg.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("Stopping HTTP Gateway server...")
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("failed to shutdown HTTP server")
			return err
		}
		log.Info().Msg("HTTP Gateway server stopped")
		return nil
	})
}

func runGrpcServer(
	ctx context.Context,
	wg *errgroup.Group,
) {
	grpcServer := grpc.NewServer()
	proto_hello.RegisterHelloServiceServer(grpcServer, &server{})

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	reflection.Register(grpcServer)

	wg.Go(func() error {
		log.Info().Msgf("GRPC server listening at %v", listener.Addr())

		if err := grpcServer.Serve(listener); err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Error().Err(err).Msg("failed to serve")
			return err
		}
		return nil
	})

	wg.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("Stopping GRPC server...")
		grpcServer.GracefulStop()
		log.Info().Msg("GRPC server stopped")
		return nil
	})

}

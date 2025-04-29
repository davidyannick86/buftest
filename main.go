package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	valid "github.com/bufbuild/protovalidate-go"
	"github.com/davidyannick86/bufbuild/testbuf/interceptor"
	protohello "github.com/davidyannick86/bufbuild/testbuf/protogen/hello/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	protohello.UnimplementedHelloServiceServer
}

func (s *server) SayHello(_ context.Context, req *protohello.SayHelloRequest) (*protohello.SayHelloResponse, error) {
	response := &protohello.SayHelloResponse{
		Message: fmt.Sprintf("Hello, %s aged : %d!", req.Name, req.GetAge()),
	}
	return response, nil
}

func (s *server) GreetManyTimes(req *protohello.GreetManyTimesRequest, stream grpc.ServerStreamingServer[protohello.GreetManyTimesResponse]) error {
	validator, err := valid.New()
	if err != nil {
		return fmt.Errorf("failed to create validator: %w", err)
	}
	log.Info().Msgf("Received request: %v", req)
	if err := validator.Validate(req); err != nil {
		var ve *valid.ValidationError
		if errors.As(err, &ve) {
			log.Info().Msgf("Validation violations: %v", ve.Error())
		}
		return fmt.Errorf("validation failed ➡️  %w", err)
	}

	name := req.GetName()
	result := fmt.Sprintf("Hello, %s!", name)
	for i := range 10 {
		response := &protohello.GreetManyTimesResponse{
			Message: result,
		}
		if err := stream.Send(response); err != nil {
			log.Error().Err(err).Msg("failed to send response")
			return fmt.Errorf("failed to send response: %w", err)
		}
		time.Sleep(1 * time.Second)
		log.Info().Msgf("Sent response %d : %v", i+1, response)
	}
	return nil
}

func (s *server) LongGreet(stream grpc.ClientStreamingServer[protohello.LongGreetRequest, protohello.LongGreetResponse]) error {
	var result string
	log.Info().Msgf("Received request: %v", stream)
	for {
		req, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if sendErr := stream.SendAndClose(&protohello.LongGreetResponse{
					Message: result,
				}); sendErr != nil {
					return fmt.Errorf("failed to send and close: %w", err)
				}
				return nil
			}
			return fmt.Errorf("failed to receive: %w", err)
		}
		log.Info().Msgf("Received request: %v", req)

		result = fmt.Sprintf("Hello, %s!", req.GetName())
		result += " ! "
	}
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
	runHTTPGateway(ctx, wg)

	err := wg.Wait()
	if err != nil {
		log.Error().Err(err).Msg("Error in server")
	} else {
		log.Info().Msg("Server stopped gracefully")
	}
}

func runHTTPGateway(
	ctx context.Context,
	wg *errgroup.Group,
) {
	grpcMux := runtime.NewServeMux()
	err := protohello.RegisterHelloServiceHandlerServer(ctx, grpcMux, &server{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register handler")
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	httpServer := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	wg.Go(func() error {
		log.Info().Msgf("HTTP Gateway server listening at %v", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Error().Err(err).Msg("failed to serve")
			return fmt.Errorf("HTTP server error: %w", err)
		}
		return nil
	})

	wg.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("Stopping HTTP Gateway server...")
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("failed to shutdown HTTP server")
			return fmt.Errorf("failed to shutdown HTTP server: %w", err)
		}
		log.Info().Msg("HTTP Gateway server stopped")
		return nil
	})
}

func runGrpcServer(
	ctx context.Context,
	wg *errgroup.Group,
) {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor),
	)
	protohello.RegisterHelloServiceServer(grpcServer, &server{})

	listener, err := net.Listen("tcp", "localhost:50051")
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
			return fmt.Errorf("GRPC server error: %w", err)
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

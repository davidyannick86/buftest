package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	valid "github.com/bufbuild/protovalidate-go"
	v1 "github.com/davidyannick86/bufbuild/testbuf/protogen/hello"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	v1.UnimplementedHelloServiceServer
}

func (s *server) SayHello(ctx context.Context, req *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	validator, err := valid.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create validator: %v", err)
	}

	if err := validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	log.Printf("Received request: %v", req)
	response := &v1.SayHelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
	}
	return response, nil
}

func main() {

	go func() {
		grpcMux := runtime.NewServeMux()
		err := v1.RegisterHelloServiceHandlerServer(context.Background(), grpcMux, &server{})
		if err != nil {
			log.Fatalf("failed to register handler: %v", err)
		}
		mux := http.NewServeMux()
		mux.Handle("/", grpcMux)

		listener, err := net.Listen("tcp", ":8080")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		log.Printf("HTTP server listening at %v", listener.Addr())
		if err := http.Serve(listener, mux); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

	}()

	grpcServer := grpc.NewServer()
	v1.RegisterHelloServiceServer(grpcServer, &server{})

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", listener.Addr())

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

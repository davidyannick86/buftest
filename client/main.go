package main

import (
	"context"
	"flag"

	proto_hello "github.com/davidyannick86/bufbuild/testbuf/protogen/hello/v1"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const serverAddress = ":50051"

func main() {

	//log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	var name string
	flag.StringVar(&name, "name", "George of the Jungle", "Name of the person to greet")
	flag.Parse()

	creds := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.NewClient(serverAddress, creds)
	if err != nil {

		log.Fatal().Msg("Failed to connect to server")
	}
	defer conn.Close()

	sayhello := proto_hello.NewHelloServiceClient(conn)
	sayhelloReq := &proto_hello.SayHelloRequest{
		Name: name,
		Age:  50,
	}

	sayHelloRespons, err := sayhello.SayHello(context.Background(), sayhelloReq)
	if err != nil {

		st, ok := status.FromError(err)
		if ok {
			log.Info().Str("error::", st.Message()).Msg("Error occurred")
		}
		return
	}
	log.Info().Msgf("Response from SayHello: %s", sayHelloRespons.Message)

}

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: hello/v1/hello.proto

package protohello

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	HelloService_SayHello_FullMethodName       = "/hello.v1.HelloService/SayHello"
	HelloService_GreetManyTimes_FullMethodName = "/hello.v1.HelloService/GreetManyTimes"
	HelloService_LongGreet_FullMethodName      = "/hello.v1.HelloService/LongGreet"
)

// HelloServiceClient is the client API for HelloService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HelloServiceClient interface {
	SayHello(ctx context.Context, in *SayHelloRequest, opts ...grpc.CallOption) (*SayHelloResponse, error)
	// Server streaming
	GreetManyTimes(ctx context.Context, in *GreetManyTimesRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[GreetManyTimesResponse], error)
	// Client streaming
	LongGreet(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[LongGreetRequest, LongGreetResponse], error)
}

type helloServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewHelloServiceClient(cc grpc.ClientConnInterface) HelloServiceClient {
	return &helloServiceClient{cc}
}

func (c *helloServiceClient) SayHello(ctx context.Context, in *SayHelloRequest, opts ...grpc.CallOption) (*SayHelloResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SayHelloResponse)
	err := c.cc.Invoke(ctx, HelloService_SayHello_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *helloServiceClient) GreetManyTimes(ctx context.Context, in *GreetManyTimesRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[GreetManyTimesResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &HelloService_ServiceDesc.Streams[0], HelloService_GreetManyTimes_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[GreetManyTimesRequest, GreetManyTimesResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type HelloService_GreetManyTimesClient = grpc.ServerStreamingClient[GreetManyTimesResponse]

func (c *helloServiceClient) LongGreet(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[LongGreetRequest, LongGreetResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &HelloService_ServiceDesc.Streams[1], HelloService_LongGreet_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[LongGreetRequest, LongGreetResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type HelloService_LongGreetClient = grpc.ClientStreamingClient[LongGreetRequest, LongGreetResponse]

// HelloServiceServer is the server API for HelloService service.
// All implementations must embed UnimplementedHelloServiceServer
// for forward compatibility.
type HelloServiceServer interface {
	SayHello(context.Context, *SayHelloRequest) (*SayHelloResponse, error)
	// Server streaming
	GreetManyTimes(*GreetManyTimesRequest, grpc.ServerStreamingServer[GreetManyTimesResponse]) error
	// Client streaming
	LongGreet(grpc.ClientStreamingServer[LongGreetRequest, LongGreetResponse]) error
	mustEmbedUnimplementedHelloServiceServer()
}

// UnimplementedHelloServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedHelloServiceServer struct{}

func (UnimplementedHelloServiceServer) SayHello(context.Context, *SayHelloRequest) (*SayHelloResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedHelloServiceServer) GreetManyTimes(*GreetManyTimesRequest, grpc.ServerStreamingServer[GreetManyTimesResponse]) error {
	return status.Errorf(codes.Unimplemented, "method GreetManyTimes not implemented")
}
func (UnimplementedHelloServiceServer) LongGreet(grpc.ClientStreamingServer[LongGreetRequest, LongGreetResponse]) error {
	return status.Errorf(codes.Unimplemented, "method LongGreet not implemented")
}
func (UnimplementedHelloServiceServer) mustEmbedUnimplementedHelloServiceServer() {}
func (UnimplementedHelloServiceServer) testEmbeddedByValue()                      {}

// UnsafeHelloServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HelloServiceServer will
// result in compilation errors.
type UnsafeHelloServiceServer interface {
	mustEmbedUnimplementedHelloServiceServer()
}

func RegisterHelloServiceServer(s grpc.ServiceRegistrar, srv HelloServiceServer) {
	// If the following call pancis, it indicates UnimplementedHelloServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&HelloService_ServiceDesc, srv)
}

func _HelloService_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SayHelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HelloServiceServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HelloService_SayHello_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HelloServiceServer).SayHello(ctx, req.(*SayHelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HelloService_GreetManyTimes_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GreetManyTimesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(HelloServiceServer).GreetManyTimes(m, &grpc.GenericServerStream[GreetManyTimesRequest, GreetManyTimesResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type HelloService_GreetManyTimesServer = grpc.ServerStreamingServer[GreetManyTimesResponse]

func _HelloService_LongGreet_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(HelloServiceServer).LongGreet(&grpc.GenericServerStream[LongGreetRequest, LongGreetResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type HelloService_LongGreetServer = grpc.ClientStreamingServer[LongGreetRequest, LongGreetResponse]

// HelloService_ServiceDesc is the grpc.ServiceDesc for HelloService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var HelloService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hello.v1.HelloService",
	HandlerType: (*HelloServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _HelloService_SayHello_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GreetManyTimes",
			Handler:       _HelloService_GreetManyTimes_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "LongGreet",
			Handler:       _HelloService_LongGreet_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "hello/v1/hello.proto",
}

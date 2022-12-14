// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: emitter.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// EmitterServiceClient is the client API for EmitterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EmitterServiceClient interface {
	SendEmitterPayload(ctx context.Context, in *Emitter, opts ...grpc.CallOption) (*EmitterResponse, error)
}

type emitterServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEmitterServiceClient(cc grpc.ClientConnInterface) EmitterServiceClient {
	return &emitterServiceClient{cc}
}

func (c *emitterServiceClient) SendEmitterPayload(ctx context.Context, in *Emitter, opts ...grpc.CallOption) (*EmitterResponse, error) {
	out := new(EmitterResponse)
	err := c.cc.Invoke(ctx, "/pb.EmitterService/SendEmitterPayload", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EmitterServiceServer is the server API for EmitterService service.
// All implementations must embed UnimplementedEmitterServiceServer
// for forward compatibility
type EmitterServiceServer interface {
	SendEmitterPayload(context.Context, *Emitter) (*EmitterResponse, error)
	mustEmbedUnimplementedEmitterServiceServer()
}

// UnimplementedEmitterServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEmitterServiceServer struct {
}

func (UnimplementedEmitterServiceServer) SendEmitterPayload(context.Context, *Emitter) (*EmitterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendEmitterPayload not implemented")
}
func (UnimplementedEmitterServiceServer) mustEmbedUnimplementedEmitterServiceServer() {}

// UnsafeEmitterServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EmitterServiceServer will
// result in compilation errors.
type UnsafeEmitterServiceServer interface {
	mustEmbedUnimplementedEmitterServiceServer()
}

func RegisterEmitterServiceServer(s grpc.ServiceRegistrar, srv EmitterServiceServer) {
	s.RegisterService(&EmitterService_ServiceDesc, srv)
}

func _EmitterService_SendEmitterPayload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Emitter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmitterServiceServer).SendEmitterPayload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.EmitterService/SendEmitterPayload",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmitterServiceServer).SendEmitterPayload(ctx, req.(*Emitter))
	}
	return interceptor(ctx, in, info, handler)
}

// EmitterService_ServiceDesc is the grpc.ServiceDesc for EmitterService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EmitterService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.EmitterService",
	HandlerType: (*EmitterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendEmitterPayload",
			Handler:    _EmitterService_SendEmitterPayload_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "emitter.proto",
}

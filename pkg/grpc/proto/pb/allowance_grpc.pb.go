// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: allowance.proto

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

// UnifrontClient is the client API for Unifront service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UnifrontClient interface {
	Allowance(ctx context.Context, in *AllowanceRequest, opts ...grpc.CallOption) (*AllowanceResponse, error)
}

type unifrontClient struct {
	cc grpc.ClientConnInterface
}

func NewUnifrontClient(cc grpc.ClientConnInterface) UnifrontClient {
	return &unifrontClient{cc}
}

func (c *unifrontClient) Allowance(ctx context.Context, in *AllowanceRequest, opts ...grpc.CallOption) (*AllowanceResponse, error) {
	out := new(AllowanceResponse)
	err := c.cc.Invoke(ctx, "/proto.Unifront/allowance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UnifrontServer is the server API for Unifront service.
// All implementations must embed UnimplementedUnifrontServer
// for forward compatibility
type UnifrontServer interface {
	Allowance(context.Context, *AllowanceRequest) (*AllowanceResponse, error)
	mustEmbedUnimplementedUnifrontServer()
}

// UnimplementedUnifrontServer must be embedded to have forward compatible implementations.
type UnimplementedUnifrontServer struct {
}

func (UnimplementedUnifrontServer) Allowance(context.Context, *AllowanceRequest) (*AllowanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Allowance not implemented")
}
func (UnimplementedUnifrontServer) mustEmbedUnimplementedUnifrontServer() {}

// UnsafeUnifrontServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UnifrontServer will
// result in compilation errors.
type UnsafeUnifrontServer interface {
	mustEmbedUnimplementedUnifrontServer()
}

func RegisterUnifrontServer(s grpc.ServiceRegistrar, srv UnifrontServer) {
	s.RegisterService(&Unifront_ServiceDesc, srv)
}

func _Unifront_Allowance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AllowanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UnifrontServer).Allowance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Unifront/allowance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UnifrontServer).Allowance(ctx, req.(*AllowanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Unifront_ServiceDesc is the grpc.ServiceDesc for Unifront service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Unifront_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Unifront",
	HandlerType: (*UnifrontServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "allowance",
			Handler:    _Unifront_Allowance_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "allowance.proto",
}

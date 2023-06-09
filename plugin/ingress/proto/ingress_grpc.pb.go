// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

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

// IngressClient is the client API for Ingress service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IngressClient interface {
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error)
}

type ingressClient struct {
	cc grpc.ClientConnInterface
}

func NewIngressClient(cc grpc.ClientConnInterface) IngressClient {
	return &ingressClient{cc}
}

func (c *ingressClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error) {
	out := new(GetReply)
	err := c.cc.Invoke(ctx, "/proto.Ingress/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IngressServer is the server API for Ingress service.
// All implementations must embed UnimplementedIngressServer
// for forward compatibility
type IngressServer interface {
	Get(context.Context, *GetRequest) (*GetReply, error)
	mustEmbedUnimplementedIngressServer()
}

// UnimplementedIngressServer must be embedded to have forward compatible implementations.
type UnimplementedIngressServer struct {
}

func (UnimplementedIngressServer) Get(context.Context, *GetRequest) (*GetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedIngressServer) mustEmbedUnimplementedIngressServer() {}

// UnsafeIngressServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IngressServer will
// result in compilation errors.
type UnsafeIngressServer interface {
	mustEmbedUnimplementedIngressServer()
}

func RegisterIngressServer(s grpc.ServiceRegistrar, srv IngressServer) {
	s.RegisterService(&Ingress_ServiceDesc, srv)
}

func _Ingress_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngressServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Ingress/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngressServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Ingress_ServiceDesc is the grpc.ServiceDesc for Ingress service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Ingress_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Ingress",
	HandlerType: (*IngressServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Ingress_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ingress.proto",
}

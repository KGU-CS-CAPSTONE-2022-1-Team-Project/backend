// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: proto/partner/partner.proto

package partner

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

// PartnerServiceClient is the client API for PartnerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PartnerServiceClient interface {
	SaveNFTInfo(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*SaveResponse, error)
	LoadNFTInfo(ctx context.Context, in *LoadRequest, opts ...grpc.CallOption) (*LoadResponse, error)
}

type partnerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPartnerServiceClient(cc grpc.ClientConnInterface) PartnerServiceClient {
	return &partnerServiceClient{cc}
}

func (c *partnerServiceClient) SaveNFTInfo(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*SaveResponse, error) {
	out := new(SaveResponse)
	err := c.cc.Invoke(ctx, "/PartnerService/SaveNFTInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *partnerServiceClient) LoadNFTInfo(ctx context.Context, in *LoadRequest, opts ...grpc.CallOption) (*LoadResponse, error) {
	out := new(LoadResponse)
	err := c.cc.Invoke(ctx, "/PartnerService/LoadNFTInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PartnerServiceServer is the server API for PartnerService service.
// All implementations must embed UnimplementedPartnerServiceServer
// for forward compatibility
type PartnerServiceServer interface {
	SaveNFTInfo(context.Context, *SaveRequest) (*SaveResponse, error)
	LoadNFTInfo(context.Context, *LoadRequest) (*LoadResponse, error)
	mustEmbedUnimplementedPartnerServiceServer()
}

// UnimplementedPartnerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPartnerServiceServer struct {
}

func (UnimplementedPartnerServiceServer) SaveNFTInfo(context.Context, *SaveRequest) (*SaveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveNFTInfo not implemented")
}
func (UnimplementedPartnerServiceServer) LoadNFTInfo(context.Context, *LoadRequest) (*LoadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadNFTInfo not implemented")
}
func (UnimplementedPartnerServiceServer) mustEmbedUnimplementedPartnerServiceServer() {}

// UnsafePartnerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PartnerServiceServer will
// result in compilation errors.
type UnsafePartnerServiceServer interface {
	mustEmbedUnimplementedPartnerServiceServer()
}

func RegisterPartnerServiceServer(s grpc.ServiceRegistrar, srv PartnerServiceServer) {
	s.RegisterService(&PartnerService_ServiceDesc, srv)
}

func _PartnerService_SaveNFTInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PartnerServiceServer).SaveNFTInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/PartnerService/SaveNFTInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PartnerServiceServer).SaveNFTInfo(ctx, req.(*SaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PartnerService_LoadNFTInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PartnerServiceServer).LoadNFTInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/PartnerService/LoadNFTInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PartnerServiceServer).LoadNFTInfo(ctx, req.(*LoadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PartnerService_ServiceDesc is the grpc.ServiceDesc for PartnerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PartnerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "PartnerService",
	HandlerType: (*PartnerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SaveNFTInfo",
			Handler:    _PartnerService_SaveNFTInfo_Handler,
		},
		{
			MethodName: "LoadNFTInfo",
			Handler:    _PartnerService_LoadNFTInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/partner/partner.proto",
}

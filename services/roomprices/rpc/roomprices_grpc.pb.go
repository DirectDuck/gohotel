// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.1
// source: roomprices.proto

package rpc

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

// RoomPricesServiceClient is the client API for RoomPricesService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RoomPricesServiceClient interface {
	GetRoomPrice(ctx context.Context, in *RoomPriceRequest, opts ...grpc.CallOption) (*RoomPriceResponse, error)
}

type roomPricesServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRoomPricesServiceClient(cc grpc.ClientConnInterface) RoomPricesServiceClient {
	return &roomPricesServiceClient{cc}
}

func (c *roomPricesServiceClient) GetRoomPrice(ctx context.Context, in *RoomPriceRequest, opts ...grpc.CallOption) (*RoomPriceResponse, error) {
	out := new(RoomPriceResponse)
	err := c.cc.Invoke(ctx, "/roomprices_rpc.RoomPricesService/GetRoomPrice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RoomPricesServiceServer is the server API for RoomPricesService service.
// All implementations must embed UnimplementedRoomPricesServiceServer
// for forward compatibility
type RoomPricesServiceServer interface {
	GetRoomPrice(context.Context, *RoomPriceRequest) (*RoomPriceResponse, error)
	mustEmbedUnimplementedRoomPricesServiceServer()
}

// UnimplementedRoomPricesServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRoomPricesServiceServer struct {
}

func (UnimplementedRoomPricesServiceServer) GetRoomPrice(context.Context, *RoomPriceRequest) (*RoomPriceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoomPrice not implemented")
}
func (UnimplementedRoomPricesServiceServer) mustEmbedUnimplementedRoomPricesServiceServer() {}

// UnsafeRoomPricesServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RoomPricesServiceServer will
// result in compilation errors.
type UnsafeRoomPricesServiceServer interface {
	mustEmbedUnimplementedRoomPricesServiceServer()
}

func RegisterRoomPricesServiceServer(s grpc.ServiceRegistrar, srv RoomPricesServiceServer) {
	s.RegisterService(&RoomPricesService_ServiceDesc, srv)
}

func _RoomPricesService_GetRoomPrice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoomPriceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomPricesServiceServer).GetRoomPrice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/roomprices_rpc.RoomPricesService/GetRoomPrice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomPricesServiceServer).GetRoomPrice(ctx, req.(*RoomPriceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RoomPricesService_ServiceDesc is the grpc.ServiceDesc for RoomPricesService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RoomPricesService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "roomprices_rpc.RoomPricesService",
	HandlerType: (*RoomPricesServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRoomPrice",
			Handler:    _RoomPricesService_GetRoomPrice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "roomprices.proto",
}

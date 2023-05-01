// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: SFS/SFS.proto

package __

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

// NodeClient is the client API for Node service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NodeClient interface {
	GetLoad(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Load, error)
	Download(ctx context.Context, in *DownloadMessage, opts ...grpc.CallOption) (*Load, error)
	GetList(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ACK, error)
}

type nodeClient struct {
	cc grpc.ClientConnInterface
}

func NewNodeClient(cc grpc.ClientConnInterface) NodeClient {
	return &nodeClient{cc}
}

func (c *nodeClient) GetLoad(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Load, error) {
	out := new(Load)
	err := c.cc.Invoke(ctx, "/SFS.Node/GetLoad", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) Download(ctx context.Context, in *DownloadMessage, opts ...grpc.CallOption) (*Load, error) {
	out := new(Load)
	err := c.cc.Invoke(ctx, "/SFS.Node/Download", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) GetList(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ACK, error) {
	out := new(ACK)
	err := c.cc.Invoke(ctx, "/SFS.Node/GetList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NodeServer is the server API for Node service.
// All implementations must embed UnimplementedNodeServer
// for forward compatibility
type NodeServer interface {
	GetLoad(context.Context, *Empty) (*Load, error)
	Download(context.Context, *DownloadMessage) (*Load, error)
	GetList(context.Context, *Empty) (*ACK, error)
	mustEmbedUnimplementedNodeServer()
}

// UnimplementedNodeServer must be embedded to have forward compatible implementations.
type UnimplementedNodeServer struct {
}

func (UnimplementedNodeServer) GetLoad(context.Context, *Empty) (*Load, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLoad not implemented")
}
func (UnimplementedNodeServer) Download(context.Context, *DownloadMessage) (*Load, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedNodeServer) GetList(context.Context, *Empty) (*ACK, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetList not implemented")
}
func (UnimplementedNodeServer) mustEmbedUnimplementedNodeServer() {}

// UnsafeNodeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NodeServer will
// result in compilation errors.
type UnsafeNodeServer interface {
	mustEmbedUnimplementedNodeServer()
}

func RegisterNodeServer(s grpc.ServiceRegistrar, srv NodeServer) {
	s.RegisterService(&Node_ServiceDesc, srv)
}

func _Node_GetLoad_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).GetLoad(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SFS.Node/GetLoad",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).GetLoad(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_Download_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).Download(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SFS.Node/Download",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).Download(ctx, req.(*DownloadMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_GetList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).GetList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SFS.Node/GetList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).GetList(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Node_ServiceDesc is the grpc.ServiceDesc for Node service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Node_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "SFS.Node",
	HandlerType: (*NodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLoad",
			Handler:    _Node_GetLoad_Handler,
		},
		{
			MethodName: "Download",
			Handler:    _Node_Download_Handler,
		},
		{
			MethodName: "GetList",
			Handler:    _Node_GetList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "SFS/SFS.proto",
}

// TrackingClient is the client API for Tracking service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TrackingClient interface {
	Find(ctx context.Context, in *String, opts ...grpc.CallOption) (*IDs, error)
	UpdateList(ctx context.Context, in *UpdateMessage, opts ...grpc.CallOption) (*ACK, error)
}

type trackingClient struct {
	cc grpc.ClientConnInterface
}

func NewTrackingClient(cc grpc.ClientConnInterface) TrackingClient {
	return &trackingClient{cc}
}

func (c *trackingClient) Find(ctx context.Context, in *String, opts ...grpc.CallOption) (*IDs, error) {
	out := new(IDs)
	err := c.cc.Invoke(ctx, "/SFS.Tracking/Find", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *trackingClient) UpdateList(ctx context.Context, in *UpdateMessage, opts ...grpc.CallOption) (*ACK, error) {
	out := new(ACK)
	err := c.cc.Invoke(ctx, "/SFS.Tracking/UpdateList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TrackingServer is the server API for Tracking service.
// All implementations must embed UnimplementedTrackingServer
// for forward compatibility
type TrackingServer interface {
	Find(context.Context, *String) (*IDs, error)
	UpdateList(context.Context, *UpdateMessage) (*ACK, error)
	mustEmbedUnimplementedTrackingServer()
}

// UnimplementedTrackingServer must be embedded to have forward compatible implementations.
type UnimplementedTrackingServer struct {
}

func (UnimplementedTrackingServer) Find(context.Context, *String) (*IDs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Find not implemented")
}
func (UnimplementedTrackingServer) UpdateList(context.Context, *UpdateMessage) (*ACK, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateList not implemented")
}
func (UnimplementedTrackingServer) mustEmbedUnimplementedTrackingServer() {}

// UnsafeTrackingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TrackingServer will
// result in compilation errors.
type UnsafeTrackingServer interface {
	mustEmbedUnimplementedTrackingServer()
}

func RegisterTrackingServer(s grpc.ServiceRegistrar, srv TrackingServer) {
	s.RegisterService(&Tracking_ServiceDesc, srv)
}

func _Tracking_Find_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackingServer).Find(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SFS.Tracking/Find",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackingServer).Find(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tracking_UpdateList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrackingServer).UpdateList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SFS.Tracking/UpdateList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrackingServer).UpdateList(ctx, req.(*UpdateMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// Tracking_ServiceDesc is the grpc.ServiceDesc for Tracking service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Tracking_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "SFS.Tracking",
	HandlerType: (*TrackingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Find",
			Handler:    _Tracking_Find_Handler,
		},
		{
			MethodName: "UpdateList",
			Handler:    _Tracking_UpdateList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "SFS/SFS.proto",
}

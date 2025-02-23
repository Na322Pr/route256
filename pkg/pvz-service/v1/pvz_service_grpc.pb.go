// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.1
// source: pvz-service/v1/pvz_service.proto

package pvz_service

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
	PVZService_ReceiveCourier_FullMethodName = "/pvz.PVZService/ReceiveCourier"
	PVZService_ReturnCourier_FullMethodName  = "/pvz.PVZService/ReturnCourier"
	PVZService_GiveOutClient_FullMethodName  = "/pvz.PVZService/GiveOutClient"
	PVZService_RefundClient_FullMethodName   = "/pvz.PVZService/RefundClient"
	PVZService_OrderList_FullMethodName      = "/pvz.PVZService/OrderList"
	PVZService_RefundList_FullMethodName     = "/pvz.PVZService/RefundList"
)

// PVZServiceClient is the client API for PVZService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PVZServiceClient interface {
	ReceiveCourier(ctx context.Context, in *ReceiveCourierRequest, opts ...grpc.CallOption) (*ReceiveCourierResponse, error)
	ReturnCourier(ctx context.Context, in *ReturnCourierRequest, opts ...grpc.CallOption) (*ReturnCourierResponse, error)
	GiveOutClient(ctx context.Context, in *GiveOutClientRequest, opts ...grpc.CallOption) (*GiveOutClientResponse, error)
	RefundClient(ctx context.Context, in *RefundClientRequest, opts ...grpc.CallOption) (*RefundClientResponse, error)
	OrderList(ctx context.Context, in *OrderListRequest, opts ...grpc.CallOption) (*OrderListResponse, error)
	RefundList(ctx context.Context, in *RefundListRequest, opts ...grpc.CallOption) (*RefundListResponse, error)
}

type pVZServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPVZServiceClient(cc grpc.ClientConnInterface) PVZServiceClient {
	return &pVZServiceClient{cc}
}

func (c *pVZServiceClient) ReceiveCourier(ctx context.Context, in *ReceiveCourierRequest, opts ...grpc.CallOption) (*ReceiveCourierResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReceiveCourierResponse)
	err := c.cc.Invoke(ctx, PVZService_ReceiveCourier_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pVZServiceClient) ReturnCourier(ctx context.Context, in *ReturnCourierRequest, opts ...grpc.CallOption) (*ReturnCourierResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReturnCourierResponse)
	err := c.cc.Invoke(ctx, PVZService_ReturnCourier_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pVZServiceClient) GiveOutClient(ctx context.Context, in *GiveOutClientRequest, opts ...grpc.CallOption) (*GiveOutClientResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GiveOutClientResponse)
	err := c.cc.Invoke(ctx, PVZService_GiveOutClient_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pVZServiceClient) RefundClient(ctx context.Context, in *RefundClientRequest, opts ...grpc.CallOption) (*RefundClientResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RefundClientResponse)
	err := c.cc.Invoke(ctx, PVZService_RefundClient_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pVZServiceClient) OrderList(ctx context.Context, in *OrderListRequest, opts ...grpc.CallOption) (*OrderListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrderListResponse)
	err := c.cc.Invoke(ctx, PVZService_OrderList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pVZServiceClient) RefundList(ctx context.Context, in *RefundListRequest, opts ...grpc.CallOption) (*RefundListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RefundListResponse)
	err := c.cc.Invoke(ctx, PVZService_RefundList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PVZServiceServer is the server API for PVZService service.
// All implementations must embed UnimplementedPVZServiceServer
// for forward compatibility.
type PVZServiceServer interface {
	ReceiveCourier(context.Context, *ReceiveCourierRequest) (*ReceiveCourierResponse, error)
	ReturnCourier(context.Context, *ReturnCourierRequest) (*ReturnCourierResponse, error)
	GiveOutClient(context.Context, *GiveOutClientRequest) (*GiveOutClientResponse, error)
	RefundClient(context.Context, *RefundClientRequest) (*RefundClientResponse, error)
	OrderList(context.Context, *OrderListRequest) (*OrderListResponse, error)
	RefundList(context.Context, *RefundListRequest) (*RefundListResponse, error)
	mustEmbedUnimplementedPVZServiceServer()
}

// UnimplementedPVZServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPVZServiceServer struct{}

func (UnimplementedPVZServiceServer) ReceiveCourier(context.Context, *ReceiveCourierRequest) (*ReceiveCourierResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReceiveCourier not implemented")
}
func (UnimplementedPVZServiceServer) ReturnCourier(context.Context, *ReturnCourierRequest) (*ReturnCourierResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReturnCourier not implemented")
}
func (UnimplementedPVZServiceServer) GiveOutClient(context.Context, *GiveOutClientRequest) (*GiveOutClientResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GiveOutClient not implemented")
}
func (UnimplementedPVZServiceServer) RefundClient(context.Context, *RefundClientRequest) (*RefundClientResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefundClient not implemented")
}
func (UnimplementedPVZServiceServer) OrderList(context.Context, *OrderListRequest) (*OrderListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderList not implemented")
}
func (UnimplementedPVZServiceServer) RefundList(context.Context, *RefundListRequest) (*RefundListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefundList not implemented")
}
func (UnimplementedPVZServiceServer) mustEmbedUnimplementedPVZServiceServer() {}
func (UnimplementedPVZServiceServer) testEmbeddedByValue()                    {}

// UnsafePVZServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PVZServiceServer will
// result in compilation errors.
type UnsafePVZServiceServer interface {
	mustEmbedUnimplementedPVZServiceServer()
}

func RegisterPVZServiceServer(s grpc.ServiceRegistrar, srv PVZServiceServer) {
	// If the following call pancis, it indicates UnimplementedPVZServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PVZService_ServiceDesc, srv)
}

func _PVZService_ReceiveCourier_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReceiveCourierRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PVZServiceServer).ReceiveCourier(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PVZService_ReceiveCourier_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PVZServiceServer).ReceiveCourier(ctx, req.(*ReceiveCourierRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PVZService_ReturnCourier_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReturnCourierRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PVZServiceServer).ReturnCourier(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PVZService_ReturnCourier_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PVZServiceServer).ReturnCourier(ctx, req.(*ReturnCourierRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PVZService_GiveOutClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GiveOutClientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PVZServiceServer).GiveOutClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PVZService_GiveOutClient_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PVZServiceServer).GiveOutClient(ctx, req.(*GiveOutClientRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PVZService_RefundClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefundClientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PVZServiceServer).RefundClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PVZService_RefundClient_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PVZServiceServer).RefundClient(ctx, req.(*RefundClientRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PVZService_OrderList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PVZServiceServer).OrderList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PVZService_OrderList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PVZServiceServer).OrderList(ctx, req.(*OrderListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PVZService_RefundList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefundListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PVZServiceServer).RefundList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PVZService_RefundList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PVZServiceServer).RefundList(ctx, req.(*RefundListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PVZService_ServiceDesc is the grpc.ServiceDesc for PVZService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PVZService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pvz.PVZService",
	HandlerType: (*PVZServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReceiveCourier",
			Handler:    _PVZService_ReceiveCourier_Handler,
		},
		{
			MethodName: "ReturnCourier",
			Handler:    _PVZService_ReturnCourier_Handler,
		},
		{
			MethodName: "GiveOutClient",
			Handler:    _PVZService_GiveOutClient_Handler,
		},
		{
			MethodName: "RefundClient",
			Handler:    _PVZService_RefundClient_Handler,
		},
		{
			MethodName: "OrderList",
			Handler:    _PVZService_OrderList_Handler,
		},
		{
			MethodName: "RefundList",
			Handler:    _PVZService_RefundList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pvz-service/v1/pvz_service.proto",
}

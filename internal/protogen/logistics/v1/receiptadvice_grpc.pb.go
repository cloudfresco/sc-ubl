// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: logistics/v1/receiptadvice.proto

package v1

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
	ReceiptAdviceHeaderService_CreateReceiptAdviceHeader_FullMethodName  = "/logistics.v1.ReceiptAdviceHeaderService/CreateReceiptAdviceHeader"
	ReceiptAdviceHeaderService_GetReceiptAdviceHeaders_FullMethodName    = "/logistics.v1.ReceiptAdviceHeaderService/GetReceiptAdviceHeaders"
	ReceiptAdviceHeaderService_GetReceiptAdviceHeader_FullMethodName     = "/logistics.v1.ReceiptAdviceHeaderService/GetReceiptAdviceHeader"
	ReceiptAdviceHeaderService_GetReceiptAdviceHeaderByPk_FullMethodName = "/logistics.v1.ReceiptAdviceHeaderService/GetReceiptAdviceHeaderByPk"
	ReceiptAdviceHeaderService_CreateReceiptAdviceLine_FullMethodName    = "/logistics.v1.ReceiptAdviceHeaderService/CreateReceiptAdviceLine"
	ReceiptAdviceHeaderService_GetReceiptAdviceLines_FullMethodName      = "/logistics.v1.ReceiptAdviceHeaderService/GetReceiptAdviceLines"
	ReceiptAdviceHeaderService_UpdateReceiptAdviceHeader_FullMethodName  = "/logistics.v1.ReceiptAdviceHeaderService/UpdateReceiptAdviceHeader"
)

// ReceiptAdviceHeaderServiceClient is the client API for ReceiptAdviceHeaderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// The ReceiptAdviceHeaderService service definition.
type ReceiptAdviceHeaderServiceClient interface {
	CreateReceiptAdviceHeader(ctx context.Context, in *CreateReceiptAdviceHeaderRequest, opts ...grpc.CallOption) (*CreateReceiptAdviceHeaderResponse, error)
	GetReceiptAdviceHeaders(ctx context.Context, in *GetReceiptAdviceHeadersRequest, opts ...grpc.CallOption) (*GetReceiptAdviceHeadersResponse, error)
	GetReceiptAdviceHeader(ctx context.Context, in *GetReceiptAdviceHeaderRequest, opts ...grpc.CallOption) (*GetReceiptAdviceHeaderResponse, error)
	GetReceiptAdviceHeaderByPk(ctx context.Context, in *GetReceiptAdviceHeaderByPkRequest, opts ...grpc.CallOption) (*GetReceiptAdviceHeaderByPkResponse, error)
	CreateReceiptAdviceLine(ctx context.Context, in *CreateReceiptAdviceLineRequest, opts ...grpc.CallOption) (*CreateReceiptAdviceLineResponse, error)
	GetReceiptAdviceLines(ctx context.Context, in *GetReceiptAdviceLinesRequest, opts ...grpc.CallOption) (*GetReceiptAdviceLinesResponse, error)
	UpdateReceiptAdviceHeader(ctx context.Context, in *UpdateReceiptAdviceHeaderRequest, opts ...grpc.CallOption) (*UpdateReceiptAdviceHeaderResponse, error)
}

type receiptAdviceHeaderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewReceiptAdviceHeaderServiceClient(cc grpc.ClientConnInterface) ReceiptAdviceHeaderServiceClient {
	return &receiptAdviceHeaderServiceClient{cc}
}

func (c *receiptAdviceHeaderServiceClient) CreateReceiptAdviceHeader(ctx context.Context, in *CreateReceiptAdviceHeaderRequest, opts ...grpc.CallOption) (*CreateReceiptAdviceHeaderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateReceiptAdviceHeaderResponse)
	err := c.cc.Invoke(ctx, ReceiptAdviceHeaderService_CreateReceiptAdviceHeader_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *receiptAdviceHeaderServiceClient) GetReceiptAdviceHeaders(ctx context.Context, in *GetReceiptAdviceHeadersRequest, opts ...grpc.CallOption) (*GetReceiptAdviceHeadersResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetReceiptAdviceHeadersResponse)
	err := c.cc.Invoke(ctx, ReceiptAdviceHeaderService_GetReceiptAdviceHeaders_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *receiptAdviceHeaderServiceClient) GetReceiptAdviceHeader(ctx context.Context, in *GetReceiptAdviceHeaderRequest, opts ...grpc.CallOption) (*GetReceiptAdviceHeaderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetReceiptAdviceHeaderResponse)
	err := c.cc.Invoke(ctx, ReceiptAdviceHeaderService_GetReceiptAdviceHeader_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *receiptAdviceHeaderServiceClient) GetReceiptAdviceHeaderByPk(ctx context.Context, in *GetReceiptAdviceHeaderByPkRequest, opts ...grpc.CallOption) (*GetReceiptAdviceHeaderByPkResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetReceiptAdviceHeaderByPkResponse)
	err := c.cc.Invoke(ctx, ReceiptAdviceHeaderService_GetReceiptAdviceHeaderByPk_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *receiptAdviceHeaderServiceClient) CreateReceiptAdviceLine(ctx context.Context, in *CreateReceiptAdviceLineRequest, opts ...grpc.CallOption) (*CreateReceiptAdviceLineResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateReceiptAdviceLineResponse)
	err := c.cc.Invoke(ctx, ReceiptAdviceHeaderService_CreateReceiptAdviceLine_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *receiptAdviceHeaderServiceClient) GetReceiptAdviceLines(ctx context.Context, in *GetReceiptAdviceLinesRequest, opts ...grpc.CallOption) (*GetReceiptAdviceLinesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetReceiptAdviceLinesResponse)
	err := c.cc.Invoke(ctx, ReceiptAdviceHeaderService_GetReceiptAdviceLines_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *receiptAdviceHeaderServiceClient) UpdateReceiptAdviceHeader(ctx context.Context, in *UpdateReceiptAdviceHeaderRequest, opts ...grpc.CallOption) (*UpdateReceiptAdviceHeaderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateReceiptAdviceHeaderResponse)
	err := c.cc.Invoke(ctx, ReceiptAdviceHeaderService_UpdateReceiptAdviceHeader_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReceiptAdviceHeaderServiceServer is the server API for ReceiptAdviceHeaderService service.
// All implementations must embed UnimplementedReceiptAdviceHeaderServiceServer
// for forward compatibility.
//
// The ReceiptAdviceHeaderService service definition.
type ReceiptAdviceHeaderServiceServer interface {
	CreateReceiptAdviceHeader(context.Context, *CreateReceiptAdviceHeaderRequest) (*CreateReceiptAdviceHeaderResponse, error)
	GetReceiptAdviceHeaders(context.Context, *GetReceiptAdviceHeadersRequest) (*GetReceiptAdviceHeadersResponse, error)
	GetReceiptAdviceHeader(context.Context, *GetReceiptAdviceHeaderRequest) (*GetReceiptAdviceHeaderResponse, error)
	GetReceiptAdviceHeaderByPk(context.Context, *GetReceiptAdviceHeaderByPkRequest) (*GetReceiptAdviceHeaderByPkResponse, error)
	CreateReceiptAdviceLine(context.Context, *CreateReceiptAdviceLineRequest) (*CreateReceiptAdviceLineResponse, error)
	GetReceiptAdviceLines(context.Context, *GetReceiptAdviceLinesRequest) (*GetReceiptAdviceLinesResponse, error)
	UpdateReceiptAdviceHeader(context.Context, *UpdateReceiptAdviceHeaderRequest) (*UpdateReceiptAdviceHeaderResponse, error)
	mustEmbedUnimplementedReceiptAdviceHeaderServiceServer()
}

// UnimplementedReceiptAdviceHeaderServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedReceiptAdviceHeaderServiceServer struct{}

func (UnimplementedReceiptAdviceHeaderServiceServer) CreateReceiptAdviceHeader(context.Context, *CreateReceiptAdviceHeaderRequest) (*CreateReceiptAdviceHeaderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateReceiptAdviceHeader not implemented")
}
func (UnimplementedReceiptAdviceHeaderServiceServer) GetReceiptAdviceHeaders(context.Context, *GetReceiptAdviceHeadersRequest) (*GetReceiptAdviceHeadersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReceiptAdviceHeaders not implemented")
}
func (UnimplementedReceiptAdviceHeaderServiceServer) GetReceiptAdviceHeader(context.Context, *GetReceiptAdviceHeaderRequest) (*GetReceiptAdviceHeaderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReceiptAdviceHeader not implemented")
}
func (UnimplementedReceiptAdviceHeaderServiceServer) GetReceiptAdviceHeaderByPk(context.Context, *GetReceiptAdviceHeaderByPkRequest) (*GetReceiptAdviceHeaderByPkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReceiptAdviceHeaderByPk not implemented")
}
func (UnimplementedReceiptAdviceHeaderServiceServer) CreateReceiptAdviceLine(context.Context, *CreateReceiptAdviceLineRequest) (*CreateReceiptAdviceLineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateReceiptAdviceLine not implemented")
}
func (UnimplementedReceiptAdviceHeaderServiceServer) GetReceiptAdviceLines(context.Context, *GetReceiptAdviceLinesRequest) (*GetReceiptAdviceLinesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReceiptAdviceLines not implemented")
}
func (UnimplementedReceiptAdviceHeaderServiceServer) UpdateReceiptAdviceHeader(context.Context, *UpdateReceiptAdviceHeaderRequest) (*UpdateReceiptAdviceHeaderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateReceiptAdviceHeader not implemented")
}
func (UnimplementedReceiptAdviceHeaderServiceServer) mustEmbedUnimplementedReceiptAdviceHeaderServiceServer() {
}
func (UnimplementedReceiptAdviceHeaderServiceServer) testEmbeddedByValue() {}

// UnsafeReceiptAdviceHeaderServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReceiptAdviceHeaderServiceServer will
// result in compilation errors.
type UnsafeReceiptAdviceHeaderServiceServer interface {
	mustEmbedUnimplementedReceiptAdviceHeaderServiceServer()
}

func RegisterReceiptAdviceHeaderServiceServer(s grpc.ServiceRegistrar, srv ReceiptAdviceHeaderServiceServer) {
	// If the following call pancis, it indicates UnimplementedReceiptAdviceHeaderServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ReceiptAdviceHeaderService_ServiceDesc, srv)
}

func _ReceiptAdviceHeaderService_CreateReceiptAdviceHeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateReceiptAdviceHeaderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReceiptAdviceHeaderServiceServer).CreateReceiptAdviceHeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReceiptAdviceHeaderService_CreateReceiptAdviceHeader_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReceiptAdviceHeaderServiceServer).CreateReceiptAdviceHeader(ctx, req.(*CreateReceiptAdviceHeaderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReceiptAdviceHeaderService_GetReceiptAdviceHeaders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReceiptAdviceHeadersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReceiptAdviceHeaderServiceServer).GetReceiptAdviceHeaders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReceiptAdviceHeaderService_GetReceiptAdviceHeaders_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReceiptAdviceHeaderServiceServer).GetReceiptAdviceHeaders(ctx, req.(*GetReceiptAdviceHeadersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReceiptAdviceHeaderService_GetReceiptAdviceHeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReceiptAdviceHeaderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReceiptAdviceHeaderServiceServer).GetReceiptAdviceHeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReceiptAdviceHeaderService_GetReceiptAdviceHeader_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReceiptAdviceHeaderServiceServer).GetReceiptAdviceHeader(ctx, req.(*GetReceiptAdviceHeaderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReceiptAdviceHeaderService_GetReceiptAdviceHeaderByPk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReceiptAdviceHeaderByPkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReceiptAdviceHeaderServiceServer).GetReceiptAdviceHeaderByPk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReceiptAdviceHeaderService_GetReceiptAdviceHeaderByPk_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReceiptAdviceHeaderServiceServer).GetReceiptAdviceHeaderByPk(ctx, req.(*GetReceiptAdviceHeaderByPkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReceiptAdviceHeaderService_CreateReceiptAdviceLine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateReceiptAdviceLineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReceiptAdviceHeaderServiceServer).CreateReceiptAdviceLine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReceiptAdviceHeaderService_CreateReceiptAdviceLine_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReceiptAdviceHeaderServiceServer).CreateReceiptAdviceLine(ctx, req.(*CreateReceiptAdviceLineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReceiptAdviceHeaderService_GetReceiptAdviceLines_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReceiptAdviceLinesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReceiptAdviceHeaderServiceServer).GetReceiptAdviceLines(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReceiptAdviceHeaderService_GetReceiptAdviceLines_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReceiptAdviceHeaderServiceServer).GetReceiptAdviceLines(ctx, req.(*GetReceiptAdviceLinesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReceiptAdviceHeaderService_UpdateReceiptAdviceHeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateReceiptAdviceHeaderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReceiptAdviceHeaderServiceServer).UpdateReceiptAdviceHeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReceiptAdviceHeaderService_UpdateReceiptAdviceHeader_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReceiptAdviceHeaderServiceServer).UpdateReceiptAdviceHeader(ctx, req.(*UpdateReceiptAdviceHeaderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ReceiptAdviceHeaderService_ServiceDesc is the grpc.ServiceDesc for ReceiptAdviceHeaderService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReceiptAdviceHeaderService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "logistics.v1.ReceiptAdviceHeaderService",
	HandlerType: (*ReceiptAdviceHeaderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateReceiptAdviceHeader",
			Handler:    _ReceiptAdviceHeaderService_CreateReceiptAdviceHeader_Handler,
		},
		{
			MethodName: "GetReceiptAdviceHeaders",
			Handler:    _ReceiptAdviceHeaderService_GetReceiptAdviceHeaders_Handler,
		},
		{
			MethodName: "GetReceiptAdviceHeader",
			Handler:    _ReceiptAdviceHeaderService_GetReceiptAdviceHeader_Handler,
		},
		{
			MethodName: "GetReceiptAdviceHeaderByPk",
			Handler:    _ReceiptAdviceHeaderService_GetReceiptAdviceHeaderByPk_Handler,
		},
		{
			MethodName: "CreateReceiptAdviceLine",
			Handler:    _ReceiptAdviceHeaderService_CreateReceiptAdviceLine_Handler,
		},
		{
			MethodName: "GetReceiptAdviceLines",
			Handler:    _ReceiptAdviceHeaderService_GetReceiptAdviceLines_Handler,
		},
		{
			MethodName: "UpdateReceiptAdviceHeader",
			Handler:    _ReceiptAdviceHeaderService_UpdateReceiptAdviceHeader_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "logistics/v1/receiptadvice.proto",
}

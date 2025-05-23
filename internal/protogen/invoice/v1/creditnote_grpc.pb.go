// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: invoice/v1/creditnote.proto

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
	CreditNoteHeaderService_CreateCreditNoteHeader_FullMethodName  = "/invoice.v1.CreditNoteHeaderService/CreateCreditNoteHeader"
	CreditNoteHeaderService_GetCreditNoteHeaders_FullMethodName    = "/invoice.v1.CreditNoteHeaderService/GetCreditNoteHeaders"
	CreditNoteHeaderService_GetCreditNoteHeader_FullMethodName     = "/invoice.v1.CreditNoteHeaderService/GetCreditNoteHeader"
	CreditNoteHeaderService_GetCreditNoteHeaderByPk_FullMethodName = "/invoice.v1.CreditNoteHeaderService/GetCreditNoteHeaderByPk"
	CreditNoteHeaderService_CreateCreditNoteLine_FullMethodName    = "/invoice.v1.CreditNoteHeaderService/CreateCreditNoteLine"
	CreditNoteHeaderService_GetCreditNoteLines_FullMethodName      = "/invoice.v1.CreditNoteHeaderService/GetCreditNoteLines"
	CreditNoteHeaderService_UpdateCreditNoteHeader_FullMethodName  = "/invoice.v1.CreditNoteHeaderService/UpdateCreditNoteHeader"
)

// CreditNoteHeaderServiceClient is the client API for CreditNoteHeaderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// The CreditNoteHeaderService service definition.
type CreditNoteHeaderServiceClient interface {
	CreateCreditNoteHeader(ctx context.Context, in *CreateCreditNoteHeaderRequest, opts ...grpc.CallOption) (*CreateCreditNoteHeaderResponse, error)
	GetCreditNoteHeaders(ctx context.Context, in *GetCreditNoteHeadersRequest, opts ...grpc.CallOption) (*GetCreditNoteHeadersResponse, error)
	GetCreditNoteHeader(ctx context.Context, in *GetCreditNoteHeaderRequest, opts ...grpc.CallOption) (*GetCreditNoteHeaderResponse, error)
	GetCreditNoteHeaderByPk(ctx context.Context, in *GetCreditNoteHeaderByPkRequest, opts ...grpc.CallOption) (*GetCreditNoteHeaderByPkResponse, error)
	CreateCreditNoteLine(ctx context.Context, in *CreateCreditNoteLineRequest, opts ...grpc.CallOption) (*CreateCreditNoteLineResponse, error)
	GetCreditNoteLines(ctx context.Context, in *GetCreditNoteLinesRequest, opts ...grpc.CallOption) (*GetCreditNoteLinesResponse, error)
	UpdateCreditNoteHeader(ctx context.Context, in *UpdateCreditNoteHeaderRequest, opts ...grpc.CallOption) (*UpdateCreditNoteHeaderResponse, error)
}

type creditNoteHeaderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCreditNoteHeaderServiceClient(cc grpc.ClientConnInterface) CreditNoteHeaderServiceClient {
	return &creditNoteHeaderServiceClient{cc}
}

func (c *creditNoteHeaderServiceClient) CreateCreditNoteHeader(ctx context.Context, in *CreateCreditNoteHeaderRequest, opts ...grpc.CallOption) (*CreateCreditNoteHeaderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateCreditNoteHeaderResponse)
	err := c.cc.Invoke(ctx, CreditNoteHeaderService_CreateCreditNoteHeader_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *creditNoteHeaderServiceClient) GetCreditNoteHeaders(ctx context.Context, in *GetCreditNoteHeadersRequest, opts ...grpc.CallOption) (*GetCreditNoteHeadersResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetCreditNoteHeadersResponse)
	err := c.cc.Invoke(ctx, CreditNoteHeaderService_GetCreditNoteHeaders_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *creditNoteHeaderServiceClient) GetCreditNoteHeader(ctx context.Context, in *GetCreditNoteHeaderRequest, opts ...grpc.CallOption) (*GetCreditNoteHeaderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetCreditNoteHeaderResponse)
	err := c.cc.Invoke(ctx, CreditNoteHeaderService_GetCreditNoteHeader_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *creditNoteHeaderServiceClient) GetCreditNoteHeaderByPk(ctx context.Context, in *GetCreditNoteHeaderByPkRequest, opts ...grpc.CallOption) (*GetCreditNoteHeaderByPkResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetCreditNoteHeaderByPkResponse)
	err := c.cc.Invoke(ctx, CreditNoteHeaderService_GetCreditNoteHeaderByPk_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *creditNoteHeaderServiceClient) CreateCreditNoteLine(ctx context.Context, in *CreateCreditNoteLineRequest, opts ...grpc.CallOption) (*CreateCreditNoteLineResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateCreditNoteLineResponse)
	err := c.cc.Invoke(ctx, CreditNoteHeaderService_CreateCreditNoteLine_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *creditNoteHeaderServiceClient) GetCreditNoteLines(ctx context.Context, in *GetCreditNoteLinesRequest, opts ...grpc.CallOption) (*GetCreditNoteLinesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetCreditNoteLinesResponse)
	err := c.cc.Invoke(ctx, CreditNoteHeaderService_GetCreditNoteLines_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *creditNoteHeaderServiceClient) UpdateCreditNoteHeader(ctx context.Context, in *UpdateCreditNoteHeaderRequest, opts ...grpc.CallOption) (*UpdateCreditNoteHeaderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateCreditNoteHeaderResponse)
	err := c.cc.Invoke(ctx, CreditNoteHeaderService_UpdateCreditNoteHeader_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CreditNoteHeaderServiceServer is the server API for CreditNoteHeaderService service.
// All implementations must embed UnimplementedCreditNoteHeaderServiceServer
// for forward compatibility.
//
// The CreditNoteHeaderService service definition.
type CreditNoteHeaderServiceServer interface {
	CreateCreditNoteHeader(context.Context, *CreateCreditNoteHeaderRequest) (*CreateCreditNoteHeaderResponse, error)
	GetCreditNoteHeaders(context.Context, *GetCreditNoteHeadersRequest) (*GetCreditNoteHeadersResponse, error)
	GetCreditNoteHeader(context.Context, *GetCreditNoteHeaderRequest) (*GetCreditNoteHeaderResponse, error)
	GetCreditNoteHeaderByPk(context.Context, *GetCreditNoteHeaderByPkRequest) (*GetCreditNoteHeaderByPkResponse, error)
	CreateCreditNoteLine(context.Context, *CreateCreditNoteLineRequest) (*CreateCreditNoteLineResponse, error)
	GetCreditNoteLines(context.Context, *GetCreditNoteLinesRequest) (*GetCreditNoteLinesResponse, error)
	UpdateCreditNoteHeader(context.Context, *UpdateCreditNoteHeaderRequest) (*UpdateCreditNoteHeaderResponse, error)
	mustEmbedUnimplementedCreditNoteHeaderServiceServer()
}

// UnimplementedCreditNoteHeaderServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCreditNoteHeaderServiceServer struct{}

func (UnimplementedCreditNoteHeaderServiceServer) CreateCreditNoteHeader(context.Context, *CreateCreditNoteHeaderRequest) (*CreateCreditNoteHeaderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCreditNoteHeader not implemented")
}
func (UnimplementedCreditNoteHeaderServiceServer) GetCreditNoteHeaders(context.Context, *GetCreditNoteHeadersRequest) (*GetCreditNoteHeadersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCreditNoteHeaders not implemented")
}
func (UnimplementedCreditNoteHeaderServiceServer) GetCreditNoteHeader(context.Context, *GetCreditNoteHeaderRequest) (*GetCreditNoteHeaderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCreditNoteHeader not implemented")
}
func (UnimplementedCreditNoteHeaderServiceServer) GetCreditNoteHeaderByPk(context.Context, *GetCreditNoteHeaderByPkRequest) (*GetCreditNoteHeaderByPkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCreditNoteHeaderByPk not implemented")
}
func (UnimplementedCreditNoteHeaderServiceServer) CreateCreditNoteLine(context.Context, *CreateCreditNoteLineRequest) (*CreateCreditNoteLineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCreditNoteLine not implemented")
}
func (UnimplementedCreditNoteHeaderServiceServer) GetCreditNoteLines(context.Context, *GetCreditNoteLinesRequest) (*GetCreditNoteLinesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCreditNoteLines not implemented")
}
func (UnimplementedCreditNoteHeaderServiceServer) UpdateCreditNoteHeader(context.Context, *UpdateCreditNoteHeaderRequest) (*UpdateCreditNoteHeaderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCreditNoteHeader not implemented")
}
func (UnimplementedCreditNoteHeaderServiceServer) mustEmbedUnimplementedCreditNoteHeaderServiceServer() {
}
func (UnimplementedCreditNoteHeaderServiceServer) testEmbeddedByValue() {}

// UnsafeCreditNoteHeaderServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CreditNoteHeaderServiceServer will
// result in compilation errors.
type UnsafeCreditNoteHeaderServiceServer interface {
	mustEmbedUnimplementedCreditNoteHeaderServiceServer()
}

func RegisterCreditNoteHeaderServiceServer(s grpc.ServiceRegistrar, srv CreditNoteHeaderServiceServer) {
	// If the following call pancis, it indicates UnimplementedCreditNoteHeaderServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CreditNoteHeaderService_ServiceDesc, srv)
}

func _CreditNoteHeaderService_CreateCreditNoteHeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCreditNoteHeaderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CreditNoteHeaderServiceServer).CreateCreditNoteHeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CreditNoteHeaderService_CreateCreditNoteHeader_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CreditNoteHeaderServiceServer).CreateCreditNoteHeader(ctx, req.(*CreateCreditNoteHeaderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CreditNoteHeaderService_GetCreditNoteHeaders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCreditNoteHeadersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CreditNoteHeaderServiceServer).GetCreditNoteHeaders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CreditNoteHeaderService_GetCreditNoteHeaders_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CreditNoteHeaderServiceServer).GetCreditNoteHeaders(ctx, req.(*GetCreditNoteHeadersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CreditNoteHeaderService_GetCreditNoteHeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCreditNoteHeaderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CreditNoteHeaderServiceServer).GetCreditNoteHeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CreditNoteHeaderService_GetCreditNoteHeader_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CreditNoteHeaderServiceServer).GetCreditNoteHeader(ctx, req.(*GetCreditNoteHeaderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CreditNoteHeaderService_GetCreditNoteHeaderByPk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCreditNoteHeaderByPkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CreditNoteHeaderServiceServer).GetCreditNoteHeaderByPk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CreditNoteHeaderService_GetCreditNoteHeaderByPk_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CreditNoteHeaderServiceServer).GetCreditNoteHeaderByPk(ctx, req.(*GetCreditNoteHeaderByPkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CreditNoteHeaderService_CreateCreditNoteLine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCreditNoteLineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CreditNoteHeaderServiceServer).CreateCreditNoteLine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CreditNoteHeaderService_CreateCreditNoteLine_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CreditNoteHeaderServiceServer).CreateCreditNoteLine(ctx, req.(*CreateCreditNoteLineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CreditNoteHeaderService_GetCreditNoteLines_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCreditNoteLinesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CreditNoteHeaderServiceServer).GetCreditNoteLines(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CreditNoteHeaderService_GetCreditNoteLines_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CreditNoteHeaderServiceServer).GetCreditNoteLines(ctx, req.(*GetCreditNoteLinesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CreditNoteHeaderService_UpdateCreditNoteHeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCreditNoteHeaderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CreditNoteHeaderServiceServer).UpdateCreditNoteHeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CreditNoteHeaderService_UpdateCreditNoteHeader_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CreditNoteHeaderServiceServer).UpdateCreditNoteHeader(ctx, req.(*UpdateCreditNoteHeaderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CreditNoteHeaderService_ServiceDesc is the grpc.ServiceDesc for CreditNoteHeaderService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CreditNoteHeaderService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "invoice.v1.CreditNoteHeaderService",
	HandlerType: (*CreditNoteHeaderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCreditNoteHeader",
			Handler:    _CreditNoteHeaderService_CreateCreditNoteHeader_Handler,
		},
		{
			MethodName: "GetCreditNoteHeaders",
			Handler:    _CreditNoteHeaderService_GetCreditNoteHeaders_Handler,
		},
		{
			MethodName: "GetCreditNoteHeader",
			Handler:    _CreditNoteHeaderService_GetCreditNoteHeader_Handler,
		},
		{
			MethodName: "GetCreditNoteHeaderByPk",
			Handler:    _CreditNoteHeaderService_GetCreditNoteHeaderByPk_Handler,
		},
		{
			MethodName: "CreateCreditNoteLine",
			Handler:    _CreditNoteHeaderService_CreateCreditNoteLine_Handler,
		},
		{
			MethodName: "GetCreditNoteLines",
			Handler:    _CreditNoteHeaderService_GetCreditNoteLines_Handler,
		},
		{
			MethodName: "UpdateCreditNoteHeader",
			Handler:    _CreditNoteHeaderService_UpdateCreditNoteHeader_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "invoice/v1/creditnote.proto",
}

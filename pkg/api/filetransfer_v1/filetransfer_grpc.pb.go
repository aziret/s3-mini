// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: api/filetransfer_v1/filetransfer.proto

package filetransfer_v1

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
	FileTransferServiceV1_UploadFile_FullMethodName     = "/file_v1.FileTransferServiceV1/UploadFile"
	FileTransferServiceV1_DownloadFile_FullMethodName   = "/file_v1.FileTransferServiceV1/DownloadFile"
	FileTransferServiceV1_RegisterClient_FullMethodName = "/file_v1.FileTransferServiceV1/RegisterClient"
)

// FileTransferServiceV1Client is the client API for FileTransferServiceV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FileTransferServiceV1Client interface {
	UploadFile(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[FileChunk, UploadStatus], error)
	DownloadFile(ctx context.Context, in *FileChunkRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[FileChunk], error)
	RegisterClient(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PongResponse, error)
}

type fileTransferServiceV1Client struct {
	cc grpc.ClientConnInterface
}

func NewFileTransferServiceV1Client(cc grpc.ClientConnInterface) FileTransferServiceV1Client {
	return &fileTransferServiceV1Client{cc}
}

func (c *fileTransferServiceV1Client) UploadFile(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[FileChunk, UploadStatus], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &FileTransferServiceV1_ServiceDesc.Streams[0], FileTransferServiceV1_UploadFile_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[FileChunk, UploadStatus]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type FileTransferServiceV1_UploadFileClient = grpc.ClientStreamingClient[FileChunk, UploadStatus]

func (c *fileTransferServiceV1Client) DownloadFile(ctx context.Context, in *FileChunkRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[FileChunk], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &FileTransferServiceV1_ServiceDesc.Streams[1], FileTransferServiceV1_DownloadFile_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[FileChunkRequest, FileChunk]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type FileTransferServiceV1_DownloadFileClient = grpc.ServerStreamingClient[FileChunk]

func (c *fileTransferServiceV1Client) RegisterClient(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PongResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PongResponse)
	err := c.cc.Invoke(ctx, FileTransferServiceV1_RegisterClient_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FileTransferServiceV1Server is the server API for FileTransferServiceV1 service.
// All implementations must embed UnimplementedFileTransferServiceV1Server
// for forward compatibility.
type FileTransferServiceV1Server interface {
	UploadFile(grpc.ClientStreamingServer[FileChunk, UploadStatus]) error
	DownloadFile(*FileChunkRequest, grpc.ServerStreamingServer[FileChunk]) error
	RegisterClient(context.Context, *PingRequest) (*PongResponse, error)
	mustEmbedUnimplementedFileTransferServiceV1Server()
}

// UnimplementedFileTransferServiceV1Server must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedFileTransferServiceV1Server struct{}

func (UnimplementedFileTransferServiceV1Server) UploadFile(grpc.ClientStreamingServer[FileChunk, UploadStatus]) error {
	return status.Errorf(codes.Unimplemented, "method UploadFile not implemented")
}
func (UnimplementedFileTransferServiceV1Server) DownloadFile(*FileChunkRequest, grpc.ServerStreamingServer[FileChunk]) error {
	return status.Errorf(codes.Unimplemented, "method DownloadFile not implemented")
}
func (UnimplementedFileTransferServiceV1Server) RegisterClient(context.Context, *PingRequest) (*PongResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterClient not implemented")
}
func (UnimplementedFileTransferServiceV1Server) mustEmbedUnimplementedFileTransferServiceV1Server() {}
func (UnimplementedFileTransferServiceV1Server) testEmbeddedByValue()                               {}

// UnsafeFileTransferServiceV1Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileTransferServiceV1Server will
// result in compilation errors.
type UnsafeFileTransferServiceV1Server interface {
	mustEmbedUnimplementedFileTransferServiceV1Server()
}

func RegisterFileTransferServiceV1Server(s grpc.ServiceRegistrar, srv FileTransferServiceV1Server) {
	// If the following call pancis, it indicates UnimplementedFileTransferServiceV1Server was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&FileTransferServiceV1_ServiceDesc, srv)
}

func _FileTransferServiceV1_UploadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FileTransferServiceV1Server).UploadFile(&grpc.GenericServerStream[FileChunk, UploadStatus]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type FileTransferServiceV1_UploadFileServer = grpc.ClientStreamingServer[FileChunk, UploadStatus]

func _FileTransferServiceV1_DownloadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FileChunkRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FileTransferServiceV1Server).DownloadFile(m, &grpc.GenericServerStream[FileChunkRequest, FileChunk]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type FileTransferServiceV1_DownloadFileServer = grpc.ServerStreamingServer[FileChunk]

func _FileTransferServiceV1_RegisterClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileTransferServiceV1Server).RegisterClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileTransferServiceV1_RegisterClient_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileTransferServiceV1Server).RegisterClient(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FileTransferServiceV1_ServiceDesc is the grpc.ServiceDesc for FileTransferServiceV1 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FileTransferServiceV1_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "file_v1.FileTransferServiceV1",
	HandlerType: (*FileTransferServiceV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterClient",
			Handler:    _FileTransferServiceV1_RegisterClient_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadFile",
			Handler:       _FileTransferServiceV1_UploadFile_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "DownloadFile",
			Handler:       _FileTransferServiceV1_DownloadFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/filetransfer_v1/filetransfer.proto",
}
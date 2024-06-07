// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: pkg/surfstore/SurfStore.proto

package surfstore

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// BlockStoreClient is the client API for BlockStore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BlockStoreClient interface {
	GetBlock(ctx context.Context, in *BlockHash, opts ...grpc.CallOption) (*Block, error)
	PutBlock(ctx context.Context, in *Block, opts ...grpc.CallOption) (*Success, error)
	MissingBlocks(ctx context.Context, in *BlockHashes, opts ...grpc.CallOption) (*BlockHashes, error)
	GetBlockHashes(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*BlockHashes, error)
}

type blockStoreClient struct {
	cc grpc.ClientConnInterface
}

func NewBlockStoreClient(cc grpc.ClientConnInterface) BlockStoreClient {
	return &blockStoreClient{cc}
}

func (c *blockStoreClient) GetBlock(ctx context.Context, in *BlockHash, opts ...grpc.CallOption) (*Block, error) {
	out := new(Block)
	err := c.cc.Invoke(ctx, "/surfstore.BlockStore/GetBlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockStoreClient) PutBlock(ctx context.Context, in *Block, opts ...grpc.CallOption) (*Success, error) {
	out := new(Success)
	err := c.cc.Invoke(ctx, "/surfstore.BlockStore/PutBlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockStoreClient) MissingBlocks(ctx context.Context, in *BlockHashes, opts ...grpc.CallOption) (*BlockHashes, error) {
	out := new(BlockHashes)
	err := c.cc.Invoke(ctx, "/surfstore.BlockStore/MissingBlocks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockStoreClient) GetBlockHashes(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*BlockHashes, error) {
	out := new(BlockHashes)
	err := c.cc.Invoke(ctx, "/surfstore.BlockStore/GetBlockHashes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BlockStoreServer is the server API for BlockStore service.
// All implementations must embed UnimplementedBlockStoreServer
// for forward compatibility
type BlockStoreServer interface {
	GetBlock(context.Context, *BlockHash) (*Block, error)
	PutBlock(context.Context, *Block) (*Success, error)
	MissingBlocks(context.Context, *BlockHashes) (*BlockHashes, error)
	GetBlockHashes(context.Context, *empty.Empty) (*BlockHashes, error)
	mustEmbedUnimplementedBlockStoreServer()
}

// UnimplementedBlockStoreServer must be embedded to have forward compatible implementations.
type UnimplementedBlockStoreServer struct {
}

func (UnimplementedBlockStoreServer) GetBlock(context.Context, *BlockHash) (*Block, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlock not implemented")
}
func (UnimplementedBlockStoreServer) PutBlock(context.Context, *Block) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutBlock not implemented")
}
func (UnimplementedBlockStoreServer) MissingBlocks(context.Context, *BlockHashes) (*BlockHashes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MissingBlocks not implemented")
}
func (UnimplementedBlockStoreServer) GetBlockHashes(context.Context, *empty.Empty) (*BlockHashes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockHashes not implemented")
}
func (UnimplementedBlockStoreServer) mustEmbedUnimplementedBlockStoreServer() {}

// UnsafeBlockStoreServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BlockStoreServer will
// result in compilation errors.
type UnsafeBlockStoreServer interface {
	mustEmbedUnimplementedBlockStoreServer()
}

func RegisterBlockStoreServer(s grpc.ServiceRegistrar, srv BlockStoreServer) {
	s.RegisterService(&BlockStore_ServiceDesc, srv)
}

func _BlockStore_GetBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BlockHash)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockStoreServer).GetBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/surfstore.BlockStore/GetBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockStoreServer).GetBlock(ctx, req.(*BlockHash))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockStore_PutBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Block)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockStoreServer).PutBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/surfstore.BlockStore/PutBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockStoreServer).PutBlock(ctx, req.(*Block))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockStore_MissingBlocks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BlockHashes)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockStoreServer).MissingBlocks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/surfstore.BlockStore/MissingBlocks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockStoreServer).MissingBlocks(ctx, req.(*BlockHashes))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockStore_GetBlockHashes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockStoreServer).GetBlockHashes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/surfstore.BlockStore/GetBlockHashes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockStoreServer).GetBlockHashes(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// BlockStore_ServiceDesc is the grpc.ServiceDesc for BlockStore service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BlockStore_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "surfstore.BlockStore",
	HandlerType: (*BlockStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetBlock",
			Handler:    _BlockStore_GetBlock_Handler,
		},
		{
			MethodName: "PutBlock",
			Handler:    _BlockStore_PutBlock_Handler,
		},
		{
			MethodName: "MissingBlocks",
			Handler:    _BlockStore_MissingBlocks_Handler,
		},
		{
			MethodName: "GetBlockHashes",
			Handler:    _BlockStore_GetBlockHashes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/surfstore/SurfStore.proto",
}

// MetaStoreClient is the client API for MetaStore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetaStoreClient interface {
	GetFileInfoMap(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*FileInfoMap, error)
	UpdateFile(ctx context.Context, in *FileMetaData, opts ...grpc.CallOption) (*Version, error)
	GetBlockStoreMap(ctx context.Context, in *BlockHashes, opts ...grpc.CallOption) (*BlockStoreMap, error)
	GetBlockStoreAddrs(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*BlockStoreAddrs, error)
}

type metaStoreClient struct {
	cc grpc.ClientConnInterface
}

func NewMetaStoreClient(cc grpc.ClientConnInterface) MetaStoreClient {
	return &metaStoreClient{cc}
}

func (c *metaStoreClient) GetFileInfoMap(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*FileInfoMap, error) {
	out := new(FileInfoMap)
	err := c.cc.Invoke(ctx, "/surfstore.MetaStore/GetFileInfoMap", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaStoreClient) UpdateFile(ctx context.Context, in *FileMetaData, opts ...grpc.CallOption) (*Version, error) {
	out := new(Version)
	err := c.cc.Invoke(ctx, "/surfstore.MetaStore/UpdateFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaStoreClient) GetBlockStoreMap(ctx context.Context, in *BlockHashes, opts ...grpc.CallOption) (*BlockStoreMap, error) {
	out := new(BlockStoreMap)
	err := c.cc.Invoke(ctx, "/surfstore.MetaStore/GetBlockStoreMap", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaStoreClient) GetBlockStoreAddrs(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*BlockStoreAddrs, error) {
	out := new(BlockStoreAddrs)
	err := c.cc.Invoke(ctx, "/surfstore.MetaStore/GetBlockStoreAddrs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetaStoreServer is the server API for MetaStore service.
// All implementations must embed UnimplementedMetaStoreServer
// for forward compatibility
type MetaStoreServer interface {
	GetFileInfoMap(context.Context, *empty.Empty) (*FileInfoMap, error)
	UpdateFile(context.Context, *FileMetaData) (*Version, error)
	GetBlockStoreMap(context.Context, *BlockHashes) (*BlockStoreMap, error)
	GetBlockStoreAddrs(context.Context, *empty.Empty) (*BlockStoreAddrs, error)
	mustEmbedUnimplementedMetaStoreServer()
}

// UnimplementedMetaStoreServer must be embedded to have forward compatible implementations.
type UnimplementedMetaStoreServer struct {
}

func (UnimplementedMetaStoreServer) GetFileInfoMap(context.Context, *empty.Empty) (*FileInfoMap, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFileInfoMap not implemented")
}
func (UnimplementedMetaStoreServer) UpdateFile(context.Context, *FileMetaData) (*Version, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateFile not implemented")
}
func (UnimplementedMetaStoreServer) GetBlockStoreMap(context.Context, *BlockHashes) (*BlockStoreMap, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockStoreMap not implemented")
}
func (UnimplementedMetaStoreServer) GetBlockStoreAddrs(context.Context, *empty.Empty) (*BlockStoreAddrs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockStoreAddrs not implemented")
}
func (UnimplementedMetaStoreServer) mustEmbedUnimplementedMetaStoreServer() {}

// UnsafeMetaStoreServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetaStoreServer will
// result in compilation errors.
type UnsafeMetaStoreServer interface {
	mustEmbedUnimplementedMetaStoreServer()
}

func RegisterMetaStoreServer(s grpc.ServiceRegistrar, srv MetaStoreServer) {
	s.RegisterService(&MetaStore_ServiceDesc, srv)
}

func _MetaStore_GetFileInfoMap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaStoreServer).GetFileInfoMap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/surfstore.MetaStore/GetFileInfoMap",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaStoreServer).GetFileInfoMap(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetaStore_UpdateFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileMetaData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaStoreServer).UpdateFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/surfstore.MetaStore/UpdateFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaStoreServer).UpdateFile(ctx, req.(*FileMetaData))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetaStore_GetBlockStoreMap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BlockHashes)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaStoreServer).GetBlockStoreMap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/surfstore.MetaStore/GetBlockStoreMap",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaStoreServer).GetBlockStoreMap(ctx, req.(*BlockHashes))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetaStore_GetBlockStoreAddrs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaStoreServer).GetBlockStoreAddrs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/surfstore.MetaStore/GetBlockStoreAddrs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaStoreServer).GetBlockStoreAddrs(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// MetaStore_ServiceDesc is the grpc.ServiceDesc for MetaStore service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MetaStore_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "surfstore.MetaStore",
	HandlerType: (*MetaStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFileInfoMap",
			Handler:    _MetaStore_GetFileInfoMap_Handler,
		},
		{
			MethodName: "UpdateFile",
			Handler:    _MetaStore_UpdateFile_Handler,
		},
		{
			MethodName: "GetBlockStoreMap",
			Handler:    _MetaStore_GetBlockStoreMap_Handler,
		},
		{
			MethodName: "GetBlockStoreAddrs",
			Handler:    _MetaStore_GetBlockStoreAddrs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/surfstore/SurfStore.proto",
}

// RaftSurfstoreClient is the client API for RaftSurfstore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RaftSurfstoreClient interface {
	// raft
	AppendEntries(ctx context.Context, in *AppendEntryInput, opts ...grpc.CallOption) (*AppendEntryOutput, error)
	SetLeader(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Success, error)
	SendHeartbeat(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Success, error)
	// metastore
	GetFileInfoMap(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*FileInfoMap, error)
	UpdateFile(ctx context.Context, in *FileMetaData, opts ...grpc.CallOption) (*Version, error)
	GetBlockStoreMap(ctx context.Context, in *BlockHashes, opts ...grpc.CallOption) (*BlockStoreMap, error)
	GetBlockStoreAddrs(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*BlockStoreAddrs, error)
	// testing interface
	GetInternalState(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*RaftInternalState, error)
	Restore(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Success, error)
	Crash(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Success, error)
	MakeServerUnreachableFrom(ctx context.Context, in *UnreachableFromServers, opts ...grpc.CallOption) (*Success, error)
}

type raftSurfstoreClient struct {
	cc grpc.ClientConnInterface
}

func NewRaftSurfstoreClient(cc grpc.ClientConnInterface) RaftSurfstoreClient {
	return &raftSurfstoreClient{cc}
}

func (c *raftSurfstoreClient) AppendEntries(ctx context.Context, in *AppendEntryInput, opts ...grpc.CallOption) (*AppendEntryOutput, error) {
	out := new(AppendEntryOutput)
	err := c.cc.Invoke(ctx, "/surfstore.RaftSurfstore/AppendEntries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftSurfstoreClient) SetLeader(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Success, error) {
	out := new(Success)
	err := c.cc.Invoke(ctx, "/surfstore.RaftSurfstore/SetLeader", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftSurfstoreClient) SendHeartbeat(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Success, error) {
	out := new(Success)
	err := c.cc.Invoke(ctx, "/surfstore.RaftSurfstore/SendHeartbeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftSurfstoreClient) GetFileInfoMap(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*FileInfoMap, error) {
	out := new(FileInfoMap)
	err := c.cc.Invoke(ctx, "/surfstore.RaftSurfstore/GetFileInfoMap", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftSurfstoreClient) UpdateFile(ctx context.Context, in *FileMetaData, opts ...grpc.CallOption) (*Version, error) {
	out := new(Version)
	err := c.cc.Invoke(ctx, "/surfstore.RaftSurfstore/UpdateFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftSurfstoreClient) GetBlockStoreMap(ctx context.Context, in *BlockHashes, opts ...grpc.CallOption) (*BlockStoreMap, error) {
	out := new(BlockStoreMap)
	err := c.cc.Invoke(ctx, "/surfstore.RaftSurfstore/GetBlockStoreMap", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftSurfstoreClient) GetBlockStoreAddrs(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*BlockStoreAddrs, error) {
	out := new(BlockStoreAddrs)
	err := c.cc.Invoke(ctx, "/surfstore.RaftSurfstore/GetBlockStoreAddrs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftSurfstoreClient) GetInternalState(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*RaftInternalState, error) {
	out := new(RaftInternalState)
	err := c.cc.Invoke(ctx, "/surfstore.RaftSurfstore/GetInternalState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftSurfstoreClient) Restore(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Success, error) {
	out := new(Success)
	err := c.cc.Invoke(ctx, "/surfstore.RaftSurfstore/Restore", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftSurfstoreClient) Crash(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Success, error) {
	out := new(Success)
	err := c.cc.Invoke(ctx, "/surfstore.RaftSurfstore/Crash", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftSurfstoreClient) MakeServerUnreachableFrom(ctx context.Context, in *UnreachableFromServers, opts ...grpc.CallOption) (*Success, error) {
	out := new(Success)
	err := c.cc.Invoke(ctx, "/surfstore.RaftSurfstore/MakeServerUnreachableFrom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RaftSurfstoreServer is the server API for RaftSurfstore service.
// All implementations must embed UnimplementedRaftSurfstoreServer
// for forward compatibility
type RaftSurfstoreServer interface {
	// raft
	AppendEntries(context.Context, *AppendEntryInput) (*AppendEntryOutput, error)
	SetLeader(context.Context, *empty.Empty) (*Success, error)
	SendHeartbeat(context.Context, *empty.Empty) (*Success, error)
	// metastore
	GetFileInfoMap(context.Context, *empty.Empty) (*FileInfoMap, error)
	UpdateFile(context.Context, *FileMetaData) (*Version, error)
	GetBlockStoreMap(context.Context, *BlockHashes) (*BlockStoreMap, error)
	GetBlockStoreAddrs(context.Context, *empty.Empty) (*BlockStoreAddrs, error)
	// testing interface
	GetInternalState(context.Context, *empty.Empty) (*RaftInternalState, error)
	Restore(context.Context, *empty.Empty) (*Success, error)
	Crash(context.Context, *empty.Empty) (*Success, error)
	MakeServerUnreachableFrom(context.Context, *UnreachableFromServers) (*Success, error)
	mustEmbedUnimplementedRaftSurfstoreServer()
}

// UnimplementedRaftSurfstoreServer must be embedded to have forward compatible implementations.
type UnimplementedRaftSurfstoreServer struct {
}

func (UnimplementedRaftSurfstoreServer) AppendEntries(context.Context, *AppendEntryInput) (*AppendEntryOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppendEntries not implemented")
}
func (UnimplementedRaftSurfstoreServer) SetLeader(context.Context, *empty.Empty) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetLeader not implemented")
}
func (UnimplementedRaftSurfstoreServer) SendHeartbeat(context.Context, *empty.Empty) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendHeartbeat not implemented")
}
func (UnimplementedRaftSurfstoreServer) GetFileInfoMap(context.Context, *empty.Empty) (*FileInfoMap, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFileInfoMap not implemented")
}
func (UnimplementedRaftSurfstoreServer) UpdateFile(context.Context, *FileMetaData) (*Version, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateFile not implemented")
}
func (UnimplementedRaftSurfstoreServer) GetBlockStoreMap(context.Context, *BlockHashes) (*BlockStoreMap, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockStoreMap not implemented")
}
func (UnimplementedRaftSurfstoreServer) GetBlockStoreAddrs(context.Context, *empty.Empty) (*BlockStoreAddrs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockStoreAddrs not implemented")
}
func (UnimplementedRaftSurfstoreServer) GetInternalState(context.Context, *empty.Empty) (*RaftInternalState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInternalState not implemented")
}
func (UnimplementedRaftSurfstoreServer) Restore(context.Context, *empty.Empty) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Restore not implemented")
}
func (UnimplementedRaftSurfstoreServer) Crash(context.Context, *empty.Empty) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Crash not implemented")
}
func (UnimplementedRaftSurfstoreServer) MakeServerUnreachableFrom(context.Context, *UnreachableFromServers) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MakeServerUnreachableFrom not implemented")
}
func (UnimplementedRaftSurfstoreServer) mustEmbedUnimplementedRaftSurfstoreServer() {}

// UnsafeRaftSurfstoreServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RaftSurfstoreServer will
// result in compilation errors.
type UnsafeRaftSurfstoreServer interface {
	mustEmbedUnimplementedRaftSurfstoreServer()
}

func RegisterRaftSurfstoreServer(s grpc.ServiceRegistrar, srv RaftSurfstoreServer) {
	s.RegisterService(&RaftSurfstore_ServiceDesc, srv)
}

func _RaftSurfstore_AppendEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendEntryInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftSurfstoreServer).AppendEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/surfstore.RaftSurfstore/AppendEntries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftSurfstoreServer).AppendEntries(ctx, req.(*AppendEntryInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftSurfstore_SetLeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftSurfstoreServer).SetLeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/surfstore.RaftSurfstore/SetLeader",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftSurfstoreServer).SetLeader(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftSurfstore_SendHeartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftSurfstoreServer).SendHeartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/surfstore.RaftSurfstore/SendHeartbeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftSurfstoreServer).SendHeartbeat(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftSurfstore_GetFileInfoMap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftSurfstoreServer).GetFileInfoMap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/surfstore.RaftSurfstore/GetFileInfoMap",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftSurfstoreServer).GetFileInfoMap(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftSurfstore_UpdateFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileMetaData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftSurfstoreServer).UpdateFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/surfstore.RaftSurfstore/UpdateFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftSurfstoreServer).UpdateFile(ctx, req.(*FileMetaData))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftSurfstore_GetBlockStoreMap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BlockHashes)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftSurfstoreServer).GetBlockStoreMap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/surfstore.RaftSurfstore/GetBlockStoreMap",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftSurfstoreServer).GetBlockStoreMap(ctx, req.(*BlockHashes))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftSurfstore_GetBlockStoreAddrs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftSurfstoreServer).GetBlockStoreAddrs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/surfstore.RaftSurfstore/GetBlockStoreAddrs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftSurfstoreServer).GetBlockStoreAddrs(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftSurfstore_GetInternalState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftSurfstoreServer).GetInternalState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/surfstore.RaftSurfstore/GetInternalState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftSurfstoreServer).GetInternalState(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftSurfstore_Restore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftSurfstoreServer).Restore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/surfstore.RaftSurfstore/Restore",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftSurfstoreServer).Restore(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftSurfstore_Crash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftSurfstoreServer).Crash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/surfstore.RaftSurfstore/Crash",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftSurfstoreServer).Crash(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftSurfstore_MakeServerUnreachableFrom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnreachableFromServers)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftSurfstoreServer).MakeServerUnreachableFrom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/surfstore.RaftSurfstore/MakeServerUnreachableFrom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftSurfstoreServer).MakeServerUnreachableFrom(ctx, req.(*UnreachableFromServers))
	}
	return interceptor(ctx, in, info, handler)
}

// RaftSurfstore_ServiceDesc is the grpc.ServiceDesc for RaftSurfstore service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RaftSurfstore_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "surfstore.RaftSurfstore",
	HandlerType: (*RaftSurfstoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AppendEntries",
			Handler:    _RaftSurfstore_AppendEntries_Handler,
		},
		{
			MethodName: "SetLeader",
			Handler:    _RaftSurfstore_SetLeader_Handler,
		},
		{
			MethodName: "SendHeartbeat",
			Handler:    _RaftSurfstore_SendHeartbeat_Handler,
		},
		{
			MethodName: "GetFileInfoMap",
			Handler:    _RaftSurfstore_GetFileInfoMap_Handler,
		},
		{
			MethodName: "UpdateFile",
			Handler:    _RaftSurfstore_UpdateFile_Handler,
		},
		{
			MethodName: "GetBlockStoreMap",
			Handler:    _RaftSurfstore_GetBlockStoreMap_Handler,
		},
		{
			MethodName: "GetBlockStoreAddrs",
			Handler:    _RaftSurfstore_GetBlockStoreAddrs_Handler,
		},
		{
			MethodName: "GetInternalState",
			Handler:    _RaftSurfstore_GetInternalState_Handler,
		},
		{
			MethodName: "Restore",
			Handler:    _RaftSurfstore_Restore_Handler,
		},
		{
			MethodName: "Crash",
			Handler:    _RaftSurfstore_Crash_Handler,
		},
		{
			MethodName: "MakeServerUnreachableFrom",
			Handler:    _RaftSurfstore_MakeServerUnreachableFrom_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/surfstore/SurfStore.proto",
}

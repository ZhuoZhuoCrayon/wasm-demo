package main

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/ZhuoZhuoCrayon/wasm-demo/src/userinfoservice/protos/userinfo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
)

// server must be embedded to have forward compatible implementations.
type server struct {
	pb.UnimplementedUserInfoServiceServer
}

type getUserInfoError struct {
	userID int32
}

var (
	openIdNotFoundError = errors.New("openid not found")
)

func (e getUserInfoError) Error() string {
	return fmt.Sprintf("faied to get user(%d) info", e.userID)
}

func (s *server) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.UserInfo, error) {

	reqMetadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, getUserInfoError{userID: req.UserId}
	}

	openid, ok := reqMetadata["openid"]
	if !ok {
		return nil, openIdNotFoundError
	}

	md := metadata.Pairs("openid", openid[0])
	// 响应前发送
	grpc.SetHeader(ctx, md)
	// 在响应的 DATA 帧后发送
	grpc.SetTrailer(ctx, md)

	return &pb.UserInfo{
		Openid:   openid[0],
		UserId:   req.UserId,
		Username: "John Doe",
		Email:    "john.doe@example.com",
		Phones: []*pb.PhoneNumber{
			{
				Number: "555-1234",
				Type:   pb.PhoneType_PHONE_TYPE_MOBILE,
			},
			{
				Number: "555-5678",
				Type:   pb.PhoneType_PHONE_TYPE_HOME,
			},
			{
				Number: "555-9012",
				Type:   pb.PhoneType_PHONE_TYPE_HOME,
			},
		},
	}, nil
}

func main() {
	// 创建一个 gRPC 服务器
	s := grpc.NewServer()

	// 注册 UserService 服务的实现
	pb.RegisterUserInfoServiceServer(s, &server{})

	// 监听一个端口，例如 9000
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 启动 gRPC 服务器
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

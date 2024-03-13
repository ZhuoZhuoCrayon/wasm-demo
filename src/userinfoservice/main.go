package main

import (
	"context"
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

func (s *server) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.UserInfo, error) {

	md := metadata.Pairs("open-id", "a-1234567Resp", "trace-id", "1234567890abcdef")
	grpc.SetHeader(ctx, md)

	println(req.UserId)
	return &pb.UserInfo{
		Openid:   "1234567890",
		UserId:   42,
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

// var _ pb.UserInfoServiceServer = server{}

func main() {
	// 创建一个 gRPC 服务器
	s := grpc.NewServer()

	// 注册 UserService 服务的实现
	pb.RegisterUserInfoServiceServer(s, &server{})

	// 监听一个端口，例如 9000
	lis, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 启动 gRPC 服务器
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

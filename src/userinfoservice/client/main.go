package main

import (
	"context"
	"fmt"
	pb "github.com/ZhuoZhuoCrayon/wasm-demo/src/userinfoservice/protos/userinfo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
)

func main() {
	// 连接到 gRPC 服务器
	conn, err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	// 创建一个 UserInfoService 客户端
	client := pb.NewUserInfoServiceClient(conn)

	// 构造一个 GetUserRequest
	request := &pb.GetUserInfoRequest{
		UserId: 42,
	}

	md := metadata.New(map[string]string{"request-open-id": "a-1234567Req"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// 调用 GetUser 方法
	response, err := client.GetUserInfo(ctx, request)
	if err != nil {
		log.Fatalf("failed to get user: %v", err)
	}

	// 打印 GetUser 方法的返回结果
	fmt.Printf("Response: %+v\n", response)
}
